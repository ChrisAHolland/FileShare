package fileshare

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"path/filepath"
	"sync"
)

type Peer struct {
	PeerID    int
	files     []string
	directory string
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
	call("SwarmMaster.RegisterFile", &fileRpcArgs, &fileRpcReply)
}

func (p *Peer) Connect() {
	request := ConnectRequest{}
	reply := ConnectReply{}
	request.PeerID = p.PeerID
	call("SwarmMaster.ConnectPeer", &request, &reply)
	if reply.Accepted == true {
		fmt.Printf("Peer %v: Connected to SwarmMaster\n", p.PeerID)
	}
}

func (p *Peer) RequestFile(file string) {
	requestFileArgs := RequestFileArgs{}
	requestFileReply := RequestFileReply{}

	requestFileArgs.PeerID = p.PeerID
	requestFileArgs.File = file
	call("SwarmMaster.ServeFile", &requestFileArgs, &requestFileReply)
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

func call(rpcname string, args interface{}, reply interface{}) bool {
	sockname := masterSock()
	c, err := rpc.DialHTTP("unix", sockname)
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

func MakePeer(id int, directory string) *Peer {
	p := Peer{}
	p.PeerID = id
	p.directory = directory
	p.files = make([]string, 10)
	return &p
}
