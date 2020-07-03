package fileshare

import (
	//"../labrpc"
	"fmt"
	"log"
	//"net"
	//"net/http"
	"net/rpc"
	//"os"
	"io/ioutil"
	"sync"
)

type Peer struct {
	PeerID int
	Files  []string
	mu     sync.Mutex
}

func SendFile(file string, id int) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
	}
	data := string(f)
	fileRpcArgs := PeerSendFile{}
	fileRpcReply := ServerReceiveFile{}
	fileRpcArgs.FileContents = data
	fileRpcArgs.FileName = file
	fileRpcArgs.PeerID = id
	call("SwarmMaster.RegisterFile", &fileRpcArgs, &fileRpcReply)
}

func (p *Peer) Connect(file string) {
	request := ConnectRequest{}
	reply := ConnectReply{}
	request.PeerID = p.PeerID
	call("SwarmMaster.ConnectPeer", &request, &reply)
	if reply.Accepted == true {
		SendFile(file, p.PeerID)
	}
}

func (p *Peer) RequestFile(file string) {
	requestFileArgs := RequestFileArgs{}
	requestFileReply := RequestFileReply{}

	requestFileArgs.PeerID = p.PeerID
	requestFileArgs.File = file
	call("SwarmMaster.ServeFile", &requestFileArgs, &requestFileReply)
	fmt.Printf("Peer %v received %v from SwarmMaster\n", p.PeerID, requestFileReply.File)
}

func call(rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
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

func MakePeer(id int) *Peer {
	p := Peer{}
	p.PeerID = id
	return &p
}
