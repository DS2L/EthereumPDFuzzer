# EthereumPDFuzzer
A fuzzing framework for detecting DoS vulnerabilities in Ethereum's peer discovery protocols.

## Modified files:
`server.go`
`v4_udp.go`
`v5_udp.go`
## Added new files:
`v4_fuzzPeers.go`
`v5_fuzzPeers.go`

## Instructions
This fuzzing project is developed based on Geth 1.11.6.

### Compile geth binaries for both nodes (Please use go 1.19 or higher versions)
  - `cd` into the `../vanillaNode/go-ethereum-1.11.6` directory and run `make geth`
  - `cd` into the `../fuzzerNode/go-ethereum-1.11.6` directory and run `make geth`

### Fuzzing DiscV4
 - Run `./startVanillaNode.sh` inside the vanillaNode directory
 - Run `./go-ethereum-1.11.6/build/bin/geth attach v4data/geth.ipc` inside of the vanillaNode directory
 - Run ` admin.nodeInfo` and look for the enode record inside of the console
 - Modify the line with `eid, err := enode.ParseID([ENODE HERE])` to have the vanilla enode in `/EthereumPeerFuzzing/fuzzerNode/go-ethereum-1.11.6/p2p/v4_fuzzPeers.go`
 - Modify the line with `addr, err := net.ResolveUDPAddr("udp", "[IP ADDRESS HERE]:[PORT HERE])` to have the vanilla node IP and Port
 - Run `./startFuzzerNode.sh` inside the `fuzzerNode` directory

### Fuzzing DiscV5
 - Run `./startVanillaNodeV5.sh` inside the vanillaNode directory
 - Run `./go-ethereum-1.11.6/build/bin/geth attach v5data/geth.ipc` inside of the vanillaNode directory
 - Run ` admin.nodeInfo` and look for the enode record inside of the console
 - Modify the line with `eid, err := enode.ParseID([ENODE HERE])` to have the vanilla enode `/EthereumPeerFuzzing/fuzzerNode/go-ethereum-1.11.6/p2p/v5_fuzzPeers.go`
 - Modify the line with `addr, err := net.ResolveUDPAddr("udp", "[IP ADDRESS HERE]:[PORT HERE])` to have the vanilla node IP and Port
 - Run `./startFuzzerNodeV5.sh` inside the `fuzzerNode` directory

## Acknowledgment
This project is supported by the Ethereum Academic Grant from the Ethereum Foundation. The project is accomplished jointly by Shixuan Guan and Darren Lee, under the guidance of Dr. Kai Li.
