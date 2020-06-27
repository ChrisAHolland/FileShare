package fileshare

import (
	"../labrpc"
	"fmt"
	"log"
	//"net"
	//"net/http"
	"net/rpc"
	//"os"
	"sync"
)

type Peer struct {
	PeerID int
	Files  []string
	mu     sync.Mutex
	peers  []*labrpc.ClientEnd
}

func Connect() {
	request := ConnectRequest{}
	reply := ConnectReply{}
	call("SwarmMaster.ConnectPeer", &request, &reply)
	if reply.Accepted == true {
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
