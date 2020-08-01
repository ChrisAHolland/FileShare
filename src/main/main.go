/*
	This file acts as a "test" program to simulate a FileShare network.
*/

package main

import (
	"../fileshare"
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	m := fileshare.MakeSwarmMaster()
	t1 := time.Now()
	elapsed := t1.Sub(start)

	fmt.Printf("SwarmMaster start time: %v\n", elapsed)

	m.MasterTest()

	p := fileshare.MakePeer(1, "testdirs/peer1/", ":60122")
	p2 := fileshare.MakePeer(2, "testdirs/peer2/", ":60123")
	p.ConnectServer()
	p2.ConnectServer()
	p.ConnectPeer(":60123", 2)
	p2.RegisterFile("test.txt")
	//p.RequestFile(":60123", 1, "test.txt")
	p.RegisterFile("test2.txt")
	//p2.RequestFile(":60122", 2, "test2.txt")
	//p.RequestFile(":60123", 1, "titanic.txt")
	p.SearchForFile("test.txt")
	p2.SearchForFile("test2.txt")
}
