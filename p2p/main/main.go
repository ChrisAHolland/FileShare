package main

import (
	"../fileshare"
	//"fmt"
	//"time"
)

func main() {
	//m := fileshare.MakeSwarmMaster()
	//m.MasterTest()

	p := fileshare.MakePeer(1, "testdirs/peer1/", ":60122")
	p2 := fileshare.MakePeer(2, "testdirs/peer2/", ":60123")
	p.Connect(p2)
	p2.Connect(p)
	p2.RegisterFile("test.txt")
	p.RequestFile(p2, "test.txt")
	/*
		p.SendFile("test.txt")
		p2.SendFile("test2.txt")
		p2.RequestFile("test.txt")
		p.RequestFile("test3.txt")
		p2.SendFile("test3.txt")
		p2.RequestFile("test3.txt")
	*/
}
