package fileshare

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

/*
	Struct type for the Peers
*/
type Peer struct {
	PeerID    int
	files     []string
	peers     []int
	numFiles  int
	numPeers  int
	directory string
	Port      string
	mu        sync.Mutex
}

/*
	A lightweight data type for the SwarmMaster and Peers to
	hold relevant information about the Peers connected
	to it, including their port and the files they posses.
*/
type PeerInfo struct {
	PeerID      int
	Port        string
	Files       [10]string
	numFiles    int
	isConnected bool
}

/*
	Method used to make Remote Procedure Calls (RPCs)
	Adopted from provided lab code
*/
func call(rpcname string, args interface{}, reply interface{}, port string) bool {
	c, err := rpc.DialHTTP("tcp", port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

func serverCall(rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("tcp", ":3123")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

/*
	Creates a server for the Peer so that other Peers can connect
*/
func (p *Peer) peerServer(port string) {
	rpc.Register(p)
	serv := rpc.NewServer()
	serv.Register(p)
	oldMux := http.DefaultServeMux
	mux := http.NewServeMux()
	http.DefaultServeMux = mux
	serv.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	http.DefaultServeMux = oldMux
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	go http.Serve(l, mux)
}

/*
	Method called to create a new Peer
*/
func MakePeer(id int, directory string, port string) *Peer {
	p := Peer{}

	p.PeerID = id
	p.directory = directory
	p.files = make([]string, 10)
	p.Port = port
	p.numFiles = 0
	p.peers = make([]int, 10)
	p.numPeers = 0

	p.peerServer(port)
	return &p
}
