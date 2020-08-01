package fileshare

import (
	"os"
	"strconv"
)

/*
	Request RPC for Peer's to connect.
*/
type ConnectRequest struct {
	PeerID int
	Port   string
}

/*
	Reply RPC for Peer's to connect.
*/
type ConnectReply struct {
	PeerID   int
	Accepted bool
}

/*
	RPC for a Peer to send a file to the server.
*/
type PeerSendFile struct {
	PeerID   int
	FileName string
}

/*
	RPC for the server to confirm it received the file.
*/
type ServerReceiveFile struct {
	FileName string
	Received bool
	Accepted bool
}

/*
	Sent by the Peer to the SwarmMaster when searching for
	a file in the network using Peer.SearchForFile().
*/
type RequestFileArgs struct {
	PeerID int
	File   string
}

/*
	Used by a peer to send another Peer a file in Peer.RequestFile()
	and Peer.ServeFile().
*/
type RequestFileReply struct {
	PeerID       int
	FileExists   bool
	ErrorMessage string
	File         string
	FileContents string
}

/*
	Sent by the SwarmMaster to a Peer indicating the details
	regarding a Peer that possesses a particular file. Used
	in Peer.SearchForFile() and SwarmMaster.SearchFile().
*/
type FindPeerReply struct {
	PeerID int
	Port   string
	File   string
	Found  bool
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
