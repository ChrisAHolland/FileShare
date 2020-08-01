/*
	This file contains the functions and RPC handlers for Peers to
	connect to the SwarmMaster and to other Peers.
	AcceptConnect():
		- Handles incoming ConnectRequest RPCs from other peers.
	ConnectPeer():
		- Requests a connection with a provided Peer via a ConnectRequest RPC.
	ConnectServer():
		- Requests a connection with the SwarmMaster via a ConnectRequest RPC.
*/

package fileshare

import "fmt"

/*
	Handles incoming connection RPCs (ConnectRequest{}) from other Peers.
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
	Connects the Peer to the provided Peer.
*/
func (p *Peer) ConnectPeer(port string, id int) {
	request := ConnectRequest{}
	reply := ConnectReply{}
	request.PeerID = p.PeerID
	request.Port = p.Port
	call("Peer.AcceptConnect", &request, &reply, port)
	if reply.Accepted == false {
		fmt.Printf("Peer %v: Connection refused from Peer %v\n", p.PeerID, id)
		return
	}
	p.peers[p.numPeers] = id
	p.numPeers = p.numPeers + 1
	fmt.Printf("Peer %v: Connected to Peer %v\n", p.PeerID, id)
}

/*
	Connects the Peer to the SwarmMaster.
*/
func (p *Peer) ConnectServer() {
	request := ConnectRequest{}
	reply := ConnectReply{}
	request.PeerID = p.PeerID
	request.Port = p.Port
	serverCall("SwarmMaster.ConnectPeer", &request, &reply)
	if reply.Accepted == true {
		fmt.Printf("Peer %v: Connected to SwarmMaster\n", p.PeerID)
	}
}
