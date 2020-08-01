/*
	This file contains the SwarmMaster structs, functions and RPC handlers.
*/

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
	SwarmMaster data type for the server.
*/
type SwarmMaster struct {
	peers    []PeerInfo
	numPeers int
	mu       sync.Mutex
}

/*
	RPC handler for when a Peer wishes to connect
	to the SwarmMaster.
*/
func (m *SwarmMaster) ConnectPeer(request *ConnectRequest, reply *ConnectReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Accepted = true
	reply.PeerID = request.PeerID
	m.peers[m.numPeers].PeerID = request.PeerID
	m.peers[m.numPeers].Port = request.Port
	m.peers[m.numPeers].isConnected = true
	m.numPeers = m.numPeers + 1
	fmt.Printf("SwarmMaster: Connected to Peer: %v\n", request.PeerID)
	return nil
}

/*
	Simple function to let us know when the
	SwarmMaster has successfully been built.
*/
func (m *SwarmMaster) MasterTest() {
	fmt.Printf("SwarmMaster is ready...\n")
}

/*
	RPC handler for when a Peer registers a file in
	the FileShare system to be shareable. This function will
	update the SwarmMaster's peers data to include the new
	file.
*/
func (m *SwarmMaster) Register(request *PeerSendFile, reply *ServerReceiveFile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Accepted = false
	reply.FileName = request.FileName
	reply.Received = true
	for i := 0; i <= m.numPeers; i++ {
		if m.peers[i].PeerID == request.PeerID {
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
	RPC handler for when a Peer is in search of a file.
	This function will search the registered files in each
	Peer's file list to find which Peer contains the requested
	file. Then a FindPeerReply RPC will be sent to the requesting
	Peer telling it how to contact the Peer with the desired file.
*/
func (m *SwarmMaster) SearchFile(request *RequestFileArgs, reply *FindPeerReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Found = false
	reply.File = request.File
	fmt.Printf("SwarmMaster: Peer %v requested a search for file %v\n", request.PeerID, request.File)
	for i := 0; i <= m.numPeers; i++ {
		for j := 0; j <= m.peers[i].numFiles; j++ {
			if request.File == m.peers[i].Files[j] {
				reply.Found = true
				reply.PeerID = m.peers[i].PeerID
				reply.Port = m.peers[i].Port
				fmt.Printf("SwarmMaster: Found file %v for Peer %v on Peer %v\n", request.File, request.PeerID, m.peers[i].PeerID)
				return nil
			}
		}
	}
	fmt.Printf("SwarmMaster: Cannot find a Peer containing file %v for Peer %v", request.File, request.PeerID)
	return nil
}

/*
	Starts a SwarmMaster's server
*/
func (m *SwarmMaster) server() {
	rpc.Register(m)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":3123")
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
