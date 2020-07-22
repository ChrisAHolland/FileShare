/*
	This file contains the functions and RPC handlers that Peers
	use to handle files.

	RequestFile():
		- Requests a file from a Peer using a RequestFileArgs RPC.

	ServeFile():
		- Handles RequestFileArgs RPCs from Peers and returns the requested
		  file (if possible) to the requesting Peer using a RequestFileReply RPC.

	RegisterFile():
		- Peers use this function to register a file in the system. This means
		  to make the file publicly shareable with other peers.

	saveFile():
		- Private function that Peers use to save a file to 'disk' once obtained from
		  another Peer.
*/

package fileshare

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
	Requests a given file from a given Peer
*/
func (p *Peer) RequestFile(peer *Peer, file string) {
	requestFileArgs := RequestFileArgs{}
	requestFileReply := RequestFileReply{}
	requestFileArgs.PeerID = p.PeerID
	requestFileArgs.File = file
	call("Peer.ServeFile", &requestFileArgs, &requestFileReply, peer.Port)

	if requestFileReply.FileExists == false {
		fmt.Printf("Peer %v: Did not receive %v from Peer %v, the file does not exist\n", p.PeerID, file, peer.PeerID)
		return
	}

	fmt.Printf("Peer %v: Received %v from Peer %v\n", p.PeerID, requestFileReply.File, peer.PeerID)
	saveFile(requestFileReply.File, requestFileReply.FileContents, p.PeerID, p.directory)
}

/*
	Handles file request RPCs (RequestFileArgs{}) from other Peers
*/
func (p *Peer) ServeFile(request *RequestFileArgs, reply *RequestFileReply) error {
	for i := 0; i <= p.numFiles; i++ {
		if p.files[i] == request.File {
			reply.FileExists = true
			reply.File = p.files[i]
			reply.FileContents = p.fileContents[i]
			reply.PeerID = request.PeerID
			fmt.Printf("Peer %v: Served file %v to Peer %v\n", p.PeerID, request.File, request.PeerID)
			return nil
		}
	}
	reply.FileExists = false
	reply.ErrorMessage = "File not found on the SwarmMaster\n"
	reply.File = request.File
	reply.PeerID = request.PeerID
	fmt.Printf("Peer %v: Peer %v requested %v, but the file does not exist\n", p.PeerID, request.PeerID, request.File)
	return nil
}

/*
	Registers a file that a Peer has on disk into the FileShare system
*/
func (p *Peer) RegisterFile(fileName string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	f, err := ioutil.ReadFile(p.directory + fileName)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	data := string(f)

	p.files[p.numFiles] = fileName
	p.fileContents[p.numFiles] = data
	p.numFiles = p.numFiles + 1

	request := PeerSendFile{}
	reply := ServerReceiveFile{}
	request.FileName = fileName
	request.PeerID = p.PeerID
	serverCall("SwarmMaster.Register", &request, &reply)
	fmt.Printf("Peer %v: Registered file %v\n", p.PeerID, fileName)
	return nil
}

/*
	Saves a newly received file to the Peer's directory
*/
func saveFile(fileName string, fileContents string, id int, directory string) {
	filePath, _ := filepath.Abs(directory + fileName)
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
