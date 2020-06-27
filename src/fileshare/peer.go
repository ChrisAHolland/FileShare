package fileshare

import (
	"../labrpc"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type Peer struct {
	PeerID  int
	Address string
	Files   []string
	mu      sync.Mutex
	peers   []*labrpc.ClientEnd
}

func (p *Peer) RequestConnect(request *ConnectRequest, reply *ConnectReply) {

}

func (p *Peer) sendRequestConnect(peer int, request *ConnectRequest, reply *ConnectReply) {
	ok := p.peers[peer].Call("Peer.RequestConnect", request, reply)
	p.mu.Lock()
	defer p.mu.Unlock()

	if ok {
	}
}

func (p *Peer) Hello() {
	fmt.Printf("hello\n")
}

func (p *Peer) server() {
	rpc.Register(p)
	rpc.HandleHTTP()

	sockname := masterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func Make() *Peer {
	p := Peer{}
	//p.PeerID = me

	return &p
}
