// Copyright 2018 The FUSION Foundation Authors
// This file is part of the fusion-dcrm library.
//
// The fusion-dcrm library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package implements the DCRM P2P.

package dcrm

import (
    //"bufio"
    "fmt"
    "os"
    "sync"

    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/p2p"
    "github.com/ethereum/go-ethereum/p2p/nat"
    "github.com/ethereum/go-ethereum/p2p/discover"
    "gopkg.in/urfave/cli.v1"
    //"github.com/ethereum/go-ethereum/log"
    //"github.com/ethereum/go-ethereum/params"
)

var (
    //args
    port        int
    bootnode    string
    transaction string

    //globle
    emitter     *Emitter
    recordfd    *os.File
    err         error
)

const (
    msgTalk   = 0
    msgLength = iota
    //log file path
    recordfile   = "/var/log/dcrmrecord.txt"
    keyfile   = "/etc/dcrm-self.key"
    dcrmenode = "/etc/dcrm-enodes.json"
)

//init log file pointer
func recordInit(){
    recordfd, err = os.OpenFile(recordfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        os.Exit(1)
    }
    defer func () {
        if err = recordfd.Close(); err != nil {
            os.Exit(1)
        }
    }()
}

func P2pInit() {
    fmt.Println("main")
    go func(){
        recordInit()

        app := cli.NewApp()
        app.Usage = "p2p dcrm Init"
        app.Action = startP2pNode
        app.Flags = []cli.Flag{
            //命令行解析得到的port
            cli.IntFlag{Name: "port", Value:0, Usage: "listen port", Destination: &port},
            //命令行解析得到bootnode
            cli.StringFlag{Name: "bootnode", Value: "", Usage: "boot node", Destination: &bootnode},
            //--transaction args to start
            cli.StringFlag{Name: "transaction", Value: "", Usage: "transaction", Destination: &transaction},
        }

        if err = app.Run(os.Args); err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
    }()
}

func startP2pNode(c *cli.Context) error {
    go func () error{
        fmt.Println("startP2pNode")
        fmt.Println("transaction: ", transaction)
        //logger := log.New()
        //logger.SetHandler(log.StderrHandler)
        emitter = NewEmitter()
        nodeKey, _ := crypto.GenerateKey()
        //nodeKey, errkey := crypto.LoadECDSA(keyfile)
        //if errkey != nil {
        //    nodeKey, _ = crypto.GenerateKey()
        //    err = crypto.SaveECDSA(keyfile, nodeKey)
        //    var kfd *os.File
        //    kfd, err = os.OpenFile(keyfile, os.O_WRONLY|os.O_APPEND, 0600)
        //    _, err = kfd.WriteString(fmt.Sprintf("\nenode://%v\n", discover.PubkeyID(&nodeKey.PublicKey)))
        //    kfd.Close()
        //}
        fmt.Printf("nodekey: %+v\n", nodeKey)


        nodeserv := p2p.Server{
            Config: p2p.Config{
                MaxPeers:   100,
                MaxPendingPeers: 100,
                NoDiscovery:   false,
                PrivateKey: nodeKey,
                Name:       "p2p DCRM",
                ListenAddr: fmt.Sprintf(":%d", port),
                Protocols:  []p2p.Protocol{emitter.MyProtocol()},
                NAT:        nat.Any(),
                //Logger:     logger,
            },
        }


        //TODO: Config.StaticNodes
        //nodeserv.Config.StaticNodes = nodeserv.Config.StaticNodes()
        //bootNodetest, errtest := discover.ParseNode("enode://2d5ca9ed6c40878196a57dc04b8e0761e229311ef36766a105e3e4d9a476b26532b2af84630689d82a88b473095d92eaab7ae799d399c0fa0cf61ae3438d6ce0@10.192.32.72:1234")
        //if errtest != nil {
        //    return errtest
        //}
        //nodeserv.Config.StaticNodes = []*discover.Node{bootNodetest}



        //TODO: 1, from args
        //从bootnode字符串中解析得到bootNode节点
        bootNode, err := discover.ParseNode(bootnode)
        //TODO: 2, read node from params
        /*var TestDcrmBootnodes = []string{
            "enode://8468f9041f9b87b5cecab15cf6f5a773fee0b2b7d3c806b43250ed9751a07e1282b1015e72104866815b0b68aa84dbbc5328f08518b6ad6ef90780f4da829517@101.132.45.75:30301",
        }*/
        //bootnode := params.TestDcrmBootnodes[0]
        //bootNode, err := discover.ParseNode(bootnode)
        if err != nil {
            return err
        }
        //p2p服务器从BootstrapNodes中得到相邻节点
        nodeserv.Config.BootstrapNodes = []*discover.Node{bootNode}
        //TODO: 3, read nodes from params
        /*bootnode := params.TestDcrmBootnodes
        cfg.BootstrapNodes = make([]*discover.Node, 0, len(urls))
        for _, url := range urls {
                node, err := discover.ParseNode(url)
                if err != nil {
                        log.Error("Bootstrap URL invalid", "enode", url, "err", err)
                        continue
                }
                cfg.BootstrapNodes = append(cfg.BootstrapNodes, node)
        }*/

        //nodeserv.Start()开启p2p服务
        if err := nodeserv.Start(); err != nil {
            return err
        }

        //emitter.self = nodeserv.NodeInfo().ID[:8]
        emitter.self = nodeserv.NodeInfo().ID[:]
        fmt.Printf("self id: %s\n", emitter.self)
        fmt.Printf("\nNodeInfo: %+v\n", nodeserv.NodeInfo())
        //go emitter.sendMsg("send message")
        select {}
    }()
    return nil
}

func (e *Emitter) MyProtocol() p2p.Protocol {
    fmt.Println("MyProtocol")
    return p2p.Protocol{
        Name:    "MyProtocol",
        Version: 1,
        Length:  msgLength,
        Run:     e.msgHandler,
    }
}

type peer struct {
    peer *p2p.Peer
    ws   p2p.MsgReadWriter
    RecvMessage []string
}

type Emitter struct {
    self  string
    peers map[string]*peer
    sync.Mutex
}

func NewEmitter() *Emitter {
    fmt.Println("NewEmitter")
    return &Emitter{peers: make(map[string]*peer)}
}

func (e *Emitter) printpeer(p *peer){
    fmt.Printf("\n\nenode: %v\npeer: %+v\nws.Protocol: %+v\n", e.self, p.peer, p.ws)
    recordfd.Seek(0, os.SEEK_END)
    n3, err := recordfd.WriteString(fmt.Sprintf("%v\n", p.peer))
    if err == nil {
        fmt.Printf("\n\nwrite %d bytes\n", n3)
    }
    recordfd.Sync()
}

func (e *Emitter) addPeer(p *p2p.Peer, ws p2p.MsgReadWriter) {
    fmt.Println("addPeer")
    //e.Lock()
    //defer e.Unlock()
    id := fmt.Sprintf("%x", p.ID().String()[:8])
    e.peers[id] = &peer{ws: ws, peer: p}

    e.printpeer(e.peers[id])
    fmt.Printf("add peer: %v\n", p)
}

func SendMsg(msg string){
    emitter.sendMsg(msg)
}

func RecvMsg() string{
    s := "{["
    for _, p := range emitter.peers {
        if p.RecvMessage == nil {
            continue
        }
        s += p.RecvMessage[0]
        s += "],["
    }
    s += "]}"
    return s
}

func (e *Emitter) sendMsg(msg string) {
    fmt.Println("sendMsg")
    //for {
        func() {
            //e.Lock()
            //defer e.Unlock()
            //inputReader := bufio.NewReader(os.Stdin)
            //fmt.Println("Please enter some input: ")
            //input, err := inputReader.ReadString('\n')
            if err == nil {
                //fmt.Printf("The input was: %s\n", input)
                for _, p := range e.peers {
                    if err := p2p.SendItems(p.ws, msgTalk, msg); err != nil {
                        //log.Println("Emitter.loopSendMsg p2p.SendItems err", err, "peer id", p.peer.ID())
                        continue
                    }
                }
            }
        }()
    //}
}

func (e *Emitter) msgHandler(peer *p2p.Peer, ws p2p.MsgReadWriter) error {
    fmt.Println("--------  msgHandler  --------")
    e.addPeer(peer, ws)
    id := fmt.Sprintf("%x", peer.ID().String()[:8])
    //e.peers[id] = &peer{ws: ws, peer: p}
    for {
        msg, err := ws.ReadMsg()
        fmt.Printf("ReadMsg: %+v\n", msg)
        if err != nil {
            return err
        }

        switch msg.Code {
        case msgTalk:
            fmt.Printf("from peer: %v, msg: %+v\n", peer, e.peers[id].RecvMessage)
            if err := msg.Decode(&e.peers[id].RecvMessage); err != nil {
                fmt.Println("decode msg err", err)
            } else {
                fmt.Println("read msg[0]:", e.peers[id].RecvMessage[0])
            }

        default:
            fmt.Println("default")
            fmt.Println("unkown msg code")
        }
    }
    return nil
}
