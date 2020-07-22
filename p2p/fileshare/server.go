package fileshare

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type SwarmMaster struct {
	peers    []PeerInfo
	numPeers int
	mu       sync.Mutex
}

type PeerInfo struct {
	PeerId   int
	Port     string
	Files    [10]string
	numFiles int
}

func (m *SwarmMaster) ConnectPeer(request *ConnectRequest, reply *ConnectReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Accepted = true
	reply.PeerID = request.PeerID
	m.peers[m.numPeers].PeerId = request.PeerID
	m.peers[m.numPeers].Port = request.Port
	m.numPeers = m.numPeers + 1
	fmt.Printf("SwarmMaster: Connected to Peer: %v\n", request.PeerID)
	return nil
}

func (m *SwarmMaster) MasterTest() {
	fmt.Printf("SwarmMaster is ready...\n")
}

func (m *SwarmMaster) Register(request *PeerSendFile, reply *ServerReceiveFile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Accepted = false
	reply.FileName = request.FileName
	reply.Received = true
	for i := 0; i <= m.numPeers; i++ {
		if m.peers[i].PeerId == request.PeerID {
			m.peers[i].numFiles++
			m.peers[i].Files[m.peers[i].numFiles] = request.FileName
			reply.Accepted = true
			fmt.Printf("SwarmMaster: Registered %v from Peer %v\n", request.FileName, request.PeerID)
			break
		}
	}
	return nil
}

/*
	Starts a SwarmMaster's server
*/
func (m *SwarmMaster) server() {
	rpc.Register(m)
	rpc.HandleHTTP()

	sockname := masterSock()
	os.Remove(sockname)

	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

/*
	Creates a SwarmMaster
	Called in main.go
*/
func MakeSwarmMaster() *SwarmMaster {
	m := SwarmMaster{}
	// 10 Peers is arbitrary
	m.peers = make([]PeerInfo, 10)
	m.numPeers = 0
	m.server()
	return &m
}
