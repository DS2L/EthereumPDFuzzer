package p2p

import (
	"fmt"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	fuzz "github.com/google/gofuzz"
	//"net"
	"time"
)

func v5_fuzzMessages(ntab *discover.UDPv5) {

	// Give time for the rest of the node to spin up
	time.Sleep(time.Duration(1) * time.Second)

	node, err := enode.Parse(enode.ValidSchemes, "enode://7cab9110de2dd3421d1451bd268880e53697bcfa0fdd86485b9269a0d6e42f6bc657a6cec60ac6be22e857eb686aed4add65fc2eb189bc0cc5767fb085cc9720@127.0.0.1:11004")
	if err != nil {
		fmt.Println(err)
	}

	ntab.Ping(node) // We send a ping message to the node we want to fuzz, so that we are not a new node to it

	fmt.Println("Starting V5Fuzzer!")
	 PINGV5Fuzz(ntab, node)
	// FindNodeV5Fuzz(ntab, node)
	//TalkRequestFuzz(ntab, node)
}

func PINGV5Fuzz(ntab *discover.UDPv5, node *enode.Node) {
	fmt.Println("Starting PingV5 Fuzzer!")
        count := 0
	f := fuzz.New()
	for {
		//fmt.Println("generate a PING message!")
                count++
                if (count%10000 == 0) {
                    log.Error("V5Fuzzer", "Ping message count", count)
                }
		reqid:= make([]byte, 8)
		f.Fuzz(&reqid)
		ntab.FuzzPING(node, reqid)
	}
}

func FindNodeV5Fuzz(ntab *discover.UDPv5, node *enode.Node) {
        fmt.Println("Starting FindNodeV5 Fuzzer!")
        count := 0
	f := fuzz.New()
	for {
		//fmt.Println("generate a FindNode message!")
		count++
                if (count%10000 == 0) {
                    log.Error("V5Fuzzer", "FindNode message count", count)
                }
		reqid:= make([]byte, 8)
                f.Fuzz(&reqid)
		var distances []uint
		f.Fuzz(&distances)
		_, err := ntab.FuzzFindNode(node, distances, reqid)
		if err != nil { // Whenever we get an error it's most likely because we sent too many invalid distances, so we resend a findnode message with a valid distance.
			fmt.Println(fmt.Sprintf("Error returned for reqid: %v, distances: %v, Error: %v", reqid, distances, err))
			ntab.FuzzFindNode(node, []uint{253}, reqid) // We send a valid distance of 253
			//return
		}
	}
}

func TalkRequestFuzz(ntab *discover.UDPv5, node *enode.Node) {
	fmt.Println("Starting TalkRes Fuzzer!")
        count := 0
	f := fuzz.New()
	for {
		//fmt.Println("generate a Talk message!")
		count++
                if (count%10000 == 0) {
                    log.Error("V5Fuzzer", "TalkRes message count", count)
                }
                reqid:= make([]byte, 8)
                f.Fuzz(&reqid)
		var protocol string
		f.Fuzz(&protocol)
		var message []byte
		f.Fuzz(&message)
		_, err := ntab.FuzzTalkRequest(node, protocol, message, reqid)
		if err != nil { 
		    fmt.Println(fmt.Sprintf("Error returned for reqid: %v, protocol: %v, message: %v,  Error: %v", reqid, protocol, message, err))
                }
	}
}

/*
func NodesFuzz(f *fuzz.Fuzzer, ntab *discover.UDPv5, eid enode.ID, toaddr *net.UDPAddr) {
	p := v5wire.Findnode{}
	f.Fuzz(&p.ReqID)
	fmt.Println("Generated Request ID: ", &p.ReqID)
	ntab.GenerateNodesAndSendMessage(f, &p, eid, toaddr)
}
*/

/*
func v5_fuzzMessages(ntab *discover.UDPv5) {

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:11003")
	if err != nil {
		return
	}

	var messageCount int
	startTime := time.Now()
	duration := time.Hour // Set the duration to 1 hour
	f := fuzz.New()

	//vanillaGeth enode
	node, err := enode.Parse(enode.ValidSchemes, "enode://e9f2d107095f8a5a276588e972a8105bcc6373a425a9d1a50ae0e6b7f383b92952bb5d791d3685c6258bd86ebef6eb398d85755fa4868cc6801a0bae04f61000@146.244.227.19:11003")
	if err != nil {
		fmt.Println(err)
	}
	for {
		// Check if the time elapsed is less than 1 hour
		if time.Since(startTime) < duration {
			fmt.Println("Sending Talk Request!")
			TalkRequestFuzz(f, ntab, node)
			messageCount++ // Increment the counter for each TalkRequestFuzz messages sent
		} else {
			break // Exit the loop once 1 hour has passed
		}
	}

	fmt.Printf("Total Talk Request Fuzz messages sent in one hour: %d\n", messageCount)
}
*/
