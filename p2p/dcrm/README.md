### FILE
go-ethereum/p2p/dcrm/dcrm.go
go-ethereum/p2p/dcrmtest.go

### P2P API
* import "dcrm"
  API: dcrm.P2pInit()
       dcrm.SendMsg(string)
       dcrm.RecvMsg() string

### DEMO
  dcrmtest.go

### make
  ./build/env.sh go build go-ethereum/p2p/dcrmtest.go

### Usage
* bootnode
     1) bootnode -genkey bootnode.key
     2) bootnode -nodekey=bootnode.key
    
* run more than two peers to test
    (--port must different)
    1) run peer1
      ./dcrmtest --port 3401 --bootnode enode://[bootnode.enodes]@[:::]:30301
    2) run peer2
      ./dcrmtest --port 3402 --bootnode enode://[bootnode.enodes]@[:::]:30301

   
### Problem: NAT
  the peers in intranet cannot connect with each other
