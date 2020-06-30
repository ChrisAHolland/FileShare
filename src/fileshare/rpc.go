package fileshare

import (
	"os"
	"strconv"
)

/*
*	Request RPC for Peer's to connect
 */
type ConnectRequest struct {
	PeerID  int
	Address string
	Port    int
}

/*
*	Reply RPC for Peer's to connect
 */
type ConnectReply struct {
	PeerID   int
	Address  string
	Port     int
	Accepted bool
}

type PeerSendFile struct {
	PeerID       int
	FileContents string
}

type ServerReceiveFile struct {
	FileContents string
	Received     bool
}

// Add your RPC definitions here.

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func masterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
