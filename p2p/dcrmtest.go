package main

import (
    "fmt"
    "time"
    "github.com/ethereum/go-ethereum/p2p/dcrm"
)

func main(){
    dcrm.P2pInit()

    time.Sleep(time.Duration(10)*time.Second)
    dcrm.SendMsg("test.........................")

    time.Sleep(time.Duration(2)*time.Second)
    ret := dcrm.RecvMsg()
    fmt.Printf("TEST Recv: %s\n", ret)

    select {}
}
