package p2p

import (
	crand "crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/discover/v4wire"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	fuzz "github.com/google/gofuzz"
	"net"
	"time"
)

func fuzzMessages(ntab *discover.UDPv4) {
	// Give time for the rest of the node to spin up
	time.Sleep(time.Duration(1) * time.Second)

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:11003")
	if err != nil {
		return
	}
	eid, err := enode.ParseID("enode://47428955a3f73b56b15bb4c6cc8f7780a3fde6f132b167d9c0d762531ae811edd74d3ec35384ce0cd91d65c98d3af352c3d548cbd719a90bc74be37415ac3b4f@4.155.157.91:11003")
	if err != nil {
		fmt.Println(err)
	}
	// PINGFuzz(ntab, eid, addr)
	// FindNodeFuzz(ntab, eid, addr)
	go PINGFuzz(ntab, eid, addr)
	go FindNodeFuzz(ntab, eid, addr)
}

func FindNodeFuzz(ntab *discover.UDPv4, eid enode.ID, toaddr *net.UDPAddr) {
	f := fuzz.New()
	fmt.Println("Starting V4 FindNode Fuzzer!")
	count := 0
	for{
		//fmt.Println("Generate a FindNode Message!")
		count++
		if (count%10000 == 0) {
		    log.Error("V4Fuzzer", "FindNode message count", count)
	        }

		var key v4wire.Pubkey
		crand.Read(key[:]) // Randomly generate a key
		f.Fuzz(&key)
		//fmt.Println("Generated Fuzz Pubkey:", key)

		var fuzzedExpirationTime uint64
		f.Fuzz(&fuzzedExpirationTime)

		//fmt.Println("Generated Fuzzed Expiration Time:", fuzzedExpirationTime)

		err := ntab.FuzzFindnode(eid, toaddr, key, fuzzedExpirationTime)
		if err != nil {
			fmt.Println("FindNode Error:", err, "Generated Fuzz Pubkey:", key, "Generated Fuzzed Expiration Time:", fuzzedExpirationTime)
			return
		}
	}
}

func PINGFuzz(ntab *discover.UDPv4, eid enode.ID, toaddr *net.UDPAddr) {
	f := fuzz.New()
	fmt.Println("Starting V4 PING Fuzzer!")
	count := 0
	for {
		count++
		if (count%10000 == 0) {
		    log.Error("V4Fuzzer", "PING message count", count)
	        }
		var fuzzedExpirationTime uint64
		f.Fuzz(&fuzzedExpirationTime)

		err := ntab.FuzzPING(eid, toaddr, fuzzedExpirationTime)
		if err != nil {
			fmt.Println("PING Error:", err, "Generated Fuzzed Expiration Time:", fuzzedExpirationTime)
			return
		} 
	}
}

/*
func generateNeighborMessagesFuzz(f *fuzz.Fuzzer, ntab *discover.UDPv4, eid enode.ID, toaddr *net.UDPAddr) {
	ntab.FuzzNeighborMessages(f, toaddr, eid)
}
*/

/*
func fuzzMessages(ntab *discover.UDPv4) {

	time.Sleep(time.Duration(1) * time.Second)

	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:11003")
	if err != nil {
		return
	}
	eid, err := enode.ParseID("enode://075b09ef2247c1753a462bfc5bf421a737203f3629245de510b2e23179a503355a5b6f2269e838f84ebf950980bcda374fba8cebf1ddd00f939aade667000eb6@146.244.227.175:11003")
	if err != nil {
		fmt.Println(err)
	}

	// Begin fuzzing code for generating FindNode messages
	f := fuzz.New()
	var messageCount int // Counter for the number of messages sent
	startTime := time.Now()
	duration := time.Hour // Set the duration to 1 hour

	for {
		// Check if the time elapsed is less than 1 hour
		if time.Since(startTime) < duration {
			fmt.Println("Send FindNode Request!")
			generateFindNodeFuzz(f, ntab, eid, addr)
			messageCount++ // Increment the counter for each message
		} else {
			break // Exit the loop once 1 hour has passed
		}
	}
	fmt.Printf("Total Find Node Fuzz messages sent in one hour: %d\n", messageCount)
}
*/
