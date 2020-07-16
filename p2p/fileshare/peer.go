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

type Peer struct {
	PeerID    int
	files     []string
	directory string
	Port      string
	mu        sync.Mutex
}

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

func (p *Peer) AcceptConnect(request *ConnectRequest, reply *ConnectReply) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	reply.Accepted = true
	reply.PeerID = request.PeerID

	fmt.Printf("Peer %v: Connected to Peer: %v\n", p.PeerID, request.PeerID)
	return nil
}

func (p *Peer) Connect(port string) {
	request := ConnectRequest{}
	reply := ConnectReply{}
	request.PeerID = p.PeerID
	call("Peer.AcceptConnect", &request, &reply, port)
	if reply.Accepted == true {
		fmt.Printf("Peer %v: Connected to SwarmMaster\n", p.PeerID)
	}
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

func MakePeer(id int, directory string, port string) *Peer {
	p := Peer{}
	p.PeerID = id
	p.directory = directory
	p.files = make([]string, 10)
	p.Port = port
	p.peerServer(port)
	return &p
}
