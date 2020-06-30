package main

import (
	"../fileshare"
	//"fmt"
	//"time"
)

func main() {
	m := fileshare.MakeSwarmMaster()
	p := fileshare.MakePeer()
	p.Connect()
	m.MasterTest()
}
