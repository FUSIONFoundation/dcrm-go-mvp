# Author: huangweijun
# Date: 20180808

### make
* ./build/env.sh go build go-ethereum/p2p/dcrmtest.go

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

### Usage
* bootnode
    Infomation: test server peer
    # bootnode -genkey bootnode.key
    # bootnode -nodekey=bootnode.key
    enode://fa6b16e6a9a1e772fa04f129c1149b7804cd79d4b335b8d19eaf6c5fbb5e9c5e45a7308ba4ebf37dffbcfe1627ce7231f6ce1e10d1080bc0a42147c199ba2fa6@101.132.45.75:8008
    
* run more than two peers to test
    (--port must different)
    1) run peer1
      ./dcrmtest --port 3401 --bootnode enode://fa6b16e6a9a1e772fa04f129c1149b7804cd79d4b335b8d19eaf6c5fbb5e9c5e45a7308ba4ebf37dffbcfe1627ce7231f6ce1e10d1080bc0a42147c199ba2fa6@101.132.45.75:8008
    2) run peer2
      ./dcrmtest --port 3402 --bootnode enode://fa6b16e6a9a1e772fa04f129c1149b7804cd79d4b335b8d19eaf6c5fbb5e9c5e45a7308ba4ebf37dffbcfe1627ce7231f6ce1e10d1080bc0a42147c199ba2fa6@101.132.45.75:8008

   
### Problem: NAT
  the peers in intranet cannot connect with each other
