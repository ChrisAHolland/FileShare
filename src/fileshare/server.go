package fileshare

import (
	"../labrpc"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type SwarmMaster struct {
	// Your definitions here.
	peers []*labrpc.ClientEnd
}

func (m *SwarmMaster) ConnectPeer(request *ConnectRequest, reply *ConnectReply) error {
	reply.Accepted = true
	return nil
}

func (m *SwarmMaster) MasterTest() {
	fmt.Printf("swarm master is made\n")
}

func (m *SwarmMaster) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func MakeSwarmMaster() *SwarmMaster {
	m := SwarmMaster{}
	m.server()
	return &m
}
