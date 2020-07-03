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
	p.Connect("test.txt")
	p2.Connect("test2.txt")
	p2.RequestFile("test.txt")
}
