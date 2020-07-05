package main

import (
	"../fileshare"
	//"fmt"
	//"time"
)

func main() {
	m := fileshare.MakeSwarmMaster()
	m.MasterTest()

	p := fileshare.MakePeer(1)
	p2 := fileshare.MakePeer(2)
	p.Connect()
	p2.Connect()
	p.SendFile("test.txt")
	p2.SendFile("test2.txt")
	p2.RequestFile("test.txt")
}
