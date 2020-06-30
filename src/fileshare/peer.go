package fileshare

import (
	"../labrpc"
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
	peers  []*labrpc.ClientEnd
}

func SendFile() {
	f, err := ioutil.ReadFile("test.txt")
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
	}
	file := string(f)
	fileRpcArgs := PeerSendFile{}
	fileRpcReply := ServerReceiveFile{}
	fileRpcArgs.FileContents = file
	call("SwarmMaster.RegisterFile", &fileRpcArgs, &fileRpcReply)
}

func (p *Peer) Connect() {
	request := ConnectRequest{}
	reply := ConnectReply{}
	call("SwarmMaster.ConnectPeer", &request, &reply)
	if reply.Accepted == true {
		SendFile()
	}
}

func (p *Peer) Hello() {
	fmt.Printf("hello\n")
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

func MakePeer() *Peer {
	p := Peer{}
	p.PeerID = 1
	return &p
}
