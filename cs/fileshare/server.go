package fileshare

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
)

type SwarmMaster struct {
	files        []string
	fileContents []string
	numFiles     int
	peers        []int
	numPeers     int
	directory    string
	mu           sync.Mutex
}

func (m *SwarmMaster) ConnectPeer(request *ConnectRequest, reply *ConnectReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Accepted = true
	reply.PeerID = request.PeerID
	m.peers[m.numPeers] = request.PeerID
	m.numPeers = m.numPeers + 1
	fmt.Printf("SwarmMaster: Connected to Peer: %v\n", request.PeerID)
	return nil
}

func (m *SwarmMaster) MasterTest() {
	fmt.Printf("SwarmMaster is ready...\n")
}

/*
	Receives a file from a Peer and saves it
*/
func (m *SwarmMaster) RegisterFile(request *PeerSendFile, reply *ServerReceiveFile) error {
	m.mu.Lock()

	reply.Received = true
	if m.CheckExisting(request.FileName) {
		reply.Accepted = false
		reply.FileName = request.FileName
		m.mu.Unlock()
		return nil
	}
	m.files[m.numFiles] = request.FileName
	m.fileContents[m.numFiles] = request.FileContents
	m.numFiles = m.numFiles + 1
	fmt.Printf("SwarmMaster: Received %v from Peer: %v\n", request.FileName, request.PeerID)

	m.mu.Unlock()
	m.serverSaveFile(request.FileName, request.FileContents)
	return nil
}

func (m *SwarmMaster) ServeFile(request *RequestFileArgs, reply *RequestFileReply) error {
	for i := 0; i <= m.numFiles; i++ {
		if m.files[i] == request.File {
			reply.FileExists = true
			reply.File = m.files[i]
			reply.FileContents = m.fileContents[i]
			reply.PeerID = request.PeerID
			fmt.Printf("SwarmMaster: Served file %v to Peer %v\n", request.File, request.PeerID)
			return nil
		}
	}
	reply.FileExists = false
	reply.ErrorMessage = "File not found on the SwarmMaster\n"
	reply.File = request.File
	reply.PeerID = request.PeerID
	fmt.Printf("SwarmMaster: Peer %v requested %v, but the file does not exist\n", request.PeerID, request.File)
	return nil
}

func (m *SwarmMaster) serverSaveFile(fileName string, fileContents string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	filePath, _ := filepath.Abs("files/" + fileName)
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating the file: %v\n", err)
		return
	}

	l, err := f.WriteString(fileContents)
	if err != nil {
		fmt.Printf("Error writing the file: %v %v\n", err, l)
		return
	}

	fmt.Printf("SwarmMaster: Saved file successfully %v\n", fileName)
}

/*
	Checks if the server already has the given file
*/
func (m *SwarmMaster) CheckExisting(file string) bool {
	for i := 0; i <= m.numFiles; i++ {
		if m.files[i] == file {
			return true
		}
	}
	return false
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
	m.files = make([]string, 10)
	m.fileContents = make([]string, 10)
	m.peers = make([]int, 10)
	m.numPeers = 0
	m.numFiles = 0
	m.server()
	return &m
}
