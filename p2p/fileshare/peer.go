package fileshare

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
)

/*
	Struct type for the Peers
*/
type Peer struct {
	PeerID       int
	files        []string
	fileContents []string
	peers        []int
	numFiles     int
	numPeers     int
	directory    string
	Port         string
	mu           sync.Mutex
}

/*
	Sends a given file to a given Peer
*/
func (p *Peer) SendFile(file string) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
	}
	data := string(f)
	fileRpcArgs := PeerSendFile{}
	fileRpcReply := ServerReceiveFile{}
	fileRpcArgs.FileContents = data
	fileRpcArgs.FileName = file
	fileRpcArgs.PeerID = p.PeerID
	call("SwarmMaster.RegisterFile", &fileRpcArgs, &fileRpcReply, p.Port)
}

/*
	Handles incoming connection RPCs (ConnectRequest{}) from other Peers
*/
func (p *Peer) AcceptConnect(request *ConnectRequest, reply *ConnectReply) error {
	fmt.Printf("Peer %v: Received ConnectRequest from Peer %v\n", p.PeerID, request.PeerID)

	p.mu.Lock()
	defer p.mu.Unlock()

	reply.Accepted = true
	reply.PeerID = request.PeerID

	p.peers[p.numPeers] = request.PeerID
	p.numPeers = p.numPeers + 1

	fmt.Printf("Peer %v: Connected to Peer: %v\n", p.PeerID, request.PeerID)
	return nil
}

/*
	Connects the Peer to the provided Peer
*/
func (p *Peer) Connect(peer *Peer) {
	request := ConnectRequest{}
	reply := ConnectReply{}
	request.PeerID = p.PeerID
	call("Peer.AcceptConnect", &request, &reply, peer.Port)
	if reply.Accepted == false {
		fmt.Printf("Peer %v: Connection refused from Peer %v\n", p.PeerID, peer.PeerID)
		return
	}
	p.peers[p.numPeers] = peer.PeerID
	p.numPeers = p.numPeers + 1
	fmt.Printf("Peer %v: Connected to Peer %v\n", p.PeerID, peer.PeerID)
}

func (p *Peer) RequestFile(file string) {
	requestFileArgs := RequestFileArgs{}
	requestFileReply := RequestFileReply{}

	requestFileArgs.PeerID = p.PeerID
	requestFileArgs.File = file
	call("SwarmMaster.ServeFile", &requestFileArgs, &requestFileReply, p.Port)
	fmt.Printf("Peer %v: Received %v from SwarmMaster\n", p.PeerID, requestFileReply.File)
	saveFile(requestFileReply.File, requestFileReply.FileContents, p.PeerID, p.directory)
}

/*
	Saves a newly received file to the Peer's directory
*/
func saveFile(fileName string, fileContents string, id int, directory string) {
	filePath, _ := filepath.Abs("peerdata/" + directory + fileName)
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
	fmt.Printf("Peer %v: Saved file successfully %v\n", id, fileName)
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
	p.fileContents = make([]string, 10)
	p.numFiles = 0
	p.peers = make([]int, 10)
	p.numPeers = 0

	p.peerServer(port)
	return &p
}
