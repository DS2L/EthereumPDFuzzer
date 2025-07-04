package discover

import (
	"bytes"
	"github.com/ethereum/go-ethereum/p2p/discover/v4wire"
	"github.com/ethereum/go-ethereum/p2p/discover/v5wire"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"net"
	"testing"
	"time"
)

func findNodeV4Fuzz(f *testing.F) {

	f.Add([]byte{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}, uint64(0))
	f.Add([]byte{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}, ^uint64(0))
	f.Add([]byte{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}, uint64(1)<<63)

	f.Fuzz(func(t *testing.T, pubkeyPayload []byte, expirationPayload uint64) {
		test := newUDPTest(t)
		defer test.close()

		fuzzedKey := v4wire.Pubkey{}
		copy(fuzzedKey[:], pubkeyPayload)

		// put a few nodes into the table. their exact
		// distribution shouldn't matter much, although we need to
		// take care not to overflow any bucket.
		nodes := &nodesByDistance{target: testTarget.ID()}
		live := make(map[enode.ID]bool)
		numCandidates := 2 * bucketSize
		for i := 0; i < numCandidates; i++ {
			key := newkey()
			ip := net.IP{10, 13, 0, byte(i)}
			n := wrapNode(enode.NewV4(&key.PublicKey, ip, 0, 2000))
			// Ensure half of table content isn't verified live yet.
			if i > numCandidates/2 {
				n.livenessChecks = 1
				live[n.ID()] = true
			}
			nodes.push(n, numCandidates)
		}
		fillTable(test.table, nodes.entries)

		// ensure there's a bond with the test node,
		// findnode won't be accepted otherwise.
		remoteID := v4wire.EncodePubkey(&test.remotekey.PublicKey).ID()
		test.table.db.UpdateLastPongReceived(remoteID, test.remoteaddr.IP, time.Now())

		// check that closest neighbors are returned.
		expected := test.table.findnodeByID(fuzzedKey.ID(), bucketSize, true)
		test.packetIn(nil, &v4wire.Findnode{Target: fuzzedKey, Expiration: expirationPayload})
		waitNeighbors := func(want []*node) {
			test.waitPacketOut(func(p *v4wire.Neighbors, to *net.UDPAddr, hash []byte) {
				if len(p.Nodes) != len(want) {
					t.Errorf("wrong number of results: got %d, want %d", len(p.Nodes), bucketSize)
					return
				}
				for i, n := range p.Nodes {
					if n.ID.ID() != want[i].ID() {
						t.Errorf("result mismatch at %d:\n  got:  %v\n  want: %v", i, n, expected.entries[i])
					}
					if !live[n.ID.ID()] {
						t.Errorf("result includes dead node %v", n.ID.ID())
					} // logic errors for these checks or remove for just crashes
				}
			})
		}
		// Receive replies.
		want := expected.entries
		if len(want) > v4wire.MaxNeighbors {
			waitNeighbors(want[:v4wire.MaxNeighbors])
			want = want[v4wire.MaxNeighbors:]
		}
		waitNeighbors(want)
	})
}

func talkRequestV5Fuzz(f *testing.F) {
	// This test checks that TALKREQ calls the registered handler function.
	f.Fuzz(func(t *testing.T, payloadReqID []byte, payloadMessage []byte, protocol string) {
		t.Parallel()
		test := newUDPV5Test(t)
		defer test.close()

		var recvMessage []byte
		test.udp.RegisterTalkHandler("test", func(id enode.ID, addr *net.UDPAddr, message []byte) []byte {
			recvMessage = message
			return []byte("test response")
		})

		// Successful case:
		test.packetIn(&v5wire.TalkRequest{
			ReqID:    payloadReqID,
			Protocol: protocol,
			Message:  payloadMessage,
		})
		test.waitPacketOut(func(p *v5wire.TalkResponse, addr *net.UDPAddr, _ v5wire.Nonce) {
			if !bytes.Equal(p.ReqID, []byte("foo")) {
				t.Error("wrong request ID in response:", p.ReqID)
			}
			if string(p.Message) != "test response" {
				t.Errorf("wrong talk response message: %q", p.Message)
			}
			if string(recvMessage) != "test request" {
				t.Errorf("wrong message received in handler: %q", recvMessage)
			}
		})

		// Check that empty response is returned for unregistered protocols.
		recvMessage = nil
		test.packetIn(&v5wire.TalkRequest{
			ReqID:    []byte("2"),
			Protocol: "wrong",
			Message:  []byte("test request"),
		})
		test.waitPacketOut(func(p *v5wire.TalkResponse, addr *net.UDPAddr, _ v5wire.Nonce) {
			if !bytes.Equal(p.ReqID, []byte("2")) {
				t.Error("wrong request ID in response:", p.ReqID)
			}
			if string(p.Message) != "" {
				t.Errorf("wrong talk response message: %q", p.Message)
			}
			if recvMessage != nil {
				t.Errorf("handler was called for wrong protocol: %q", recvMessage)
			}
		})
	})
}

// This test checks that incoming FINDNODE calls are handled correctly.
func findNodeV5Fuzz(f *testing.F) {
	f.Fuzz(func(t *testing.T, payloadReqID []byte, payloadDistances []uint) {
		t.Parallel()
		test := newUDPV5Test(t)
		defer test.close()

		// Create test nodes and insert them into the table.
		nodes253 := nodesAtDistance(test.table.self().ID(), 253, 16)
		nodes249 := nodesAtDistance(test.table.self().ID(), 249, 4)
		nodes248 := nodesAtDistance(test.table.self().ID(), 248, 10)
		fillTable(test.table, wrapNodes(nodes253))
		fillTable(test.table, wrapNodes(nodes249))
		fillTable(test.table, wrapNodes(nodes248))

		// Requesting with distance zero should return the node's own record.
		test.packetIn(&v5wire.Findnode{ReqID: payloadReqID, Distances: payloadDistances})
	})
}
