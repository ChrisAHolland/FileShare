package main

import (
	"../fileshare"
	//"fmt"
	//"time"
)

func main() {
	//m := fileshare.MakeSwarmMaster()
	//m.MasterTest()

	p := fileshare.MakePeer(1, "peer1/", ":60122")
	p2 := fileshare.MakePeer(2, "peer2/", ":60123")
	p.Connect(p2)
	p2.Connect(p)
	/*
		p.SendFile("test.txt")
		p2.SendFile("test2.txt")
		p2.RequestFile("test.txt")
		p.RequestFile("test3.txt")
		p2.SendFile("test3.txt")
		p2.RequestFile("test3.txt")
	*/
}
