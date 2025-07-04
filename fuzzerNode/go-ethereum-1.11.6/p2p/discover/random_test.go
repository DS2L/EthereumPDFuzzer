package discover

import (
	"fmt"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"testing"
)

func TestNodeGeneration(t *testing.T) {
	var r enr.Record
	r.Set(enr.IP{127, 0, 0, 1})
	r.Set(enr.UDP(11001))
	r.Set(enr.TCP(11001))
	node, err := enode.GenerateNodeForFuzz(&r)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(node.IP())
}
