# Submission for the EthGlobal 2023 Hackathon. 

## EthnodComm
 Ethnod Communicator is derived from methods created in the [go-ethereum](https://github.com/ethereum/go-ethereum) library
 Mainly those logic from [cmd/devp2p](https://github.com/ethereum/go-ethereum/tree/master/cmd/devp2p)

 It Uses the `lookuppubkey` method with in the [`p2p/discover/v4_udp.go`](https://github.com/ethereum/go-ethereum/blob/master/p2p/discover/v4_udp.go) library To discover the neighbors of TargetNodes.

## How to use it 
Takes in an `enode://` encoded bootnode via the `targetnode` flag. It them derives a set of nearby neighbors
based on that bootnodes public key.  

  Example with bootnode:
  ```
   go run main.go --targetnode enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166:30303
  ```
   
## Flags
`targetnode` - "Target node to connect to in format `enode://` format"

