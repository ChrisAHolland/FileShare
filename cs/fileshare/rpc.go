package fileshare

import (
	"os"
	"strconv"
)

/*
	Request RPC for Peer's to connect
*/
type ConnectRequest struct {
	PeerID int
}

/*
	Reply RPC for Peer's to connect
*/
type ConnectReply struct {
	PeerID   int
	Accepted bool
}

/*
	RPC for a Peer to send a file to the server
*/
type PeerSendFile struct {
	PeerID       int
	FileName     string
	FileContents string
}

/*
	RPC for the server to confirm it received the file
*/
type ServerReceiveFile struct {
	FileName string
	Received bool
	Accepted bool
}

type RequestFileArgs struct {
	PeerID int
	File   string
}

type RequestFileReply struct {
	PeerID       int
	FileExists   bool
	ErrorMessage string
	File         string
	FileContents string
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
