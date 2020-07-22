package main

import (
	"../fileshare"
	//"fmt"
	//"time"
)

func main() {
	m := fileshare.MakeSwarmMaster()
	m.MasterTest()

	p := fileshare.MakePeer(1, "testdirs/peer1/", ":60122")
	p2 := fileshare.MakePeer(2, "testdirs/peer2/", ":60123")
	p.ConnectServer()
	p2.ConnectServer()
	p.ConnectPeer(p2)
	p2.RegisterFile("test.txt")
	p.RequestFile(p2, "test.txt")
	p.RegisterFile("test2.txt")
	p2.RequestFile(p, "test2.txt")
	p.RequestFile(p2, "titanic.txt")
}
