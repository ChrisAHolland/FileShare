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

	t2 := time.Now()

	p1 := fileshare.MakePeer(1, "testdirs/peer1/", ":60121")
	p2 := fileshare.MakePeer(2, "testdirs/peer2/", ":60122")
	p3 := fileshare.MakePeer(3, "testdirs/peer3/", ":60123")
	p4 := fileshare.MakePeer(4, "testdirs/peer4/", ":60124")
	p5 := fileshare.MakePeer(5, "testdirs/peer5/", ":60125")

	t3 := time.Now()
	peerCreationTime := t3.Sub(t2)
	fmt.Printf("Peer creation time (5): %v\n", peerCreationTime)

	t4 := time.Now()

	p1.ConnectServer()
	p2.ConnectServer()
	p3.ConnectServer()
	p4.ConnectServer()
	p5.ConnectServer()

	t5 := time.Now()
	peerConnectServerTime := t5.Sub(t4)
	fmt.Printf("Peer connect to server time (5): %v\n", peerConnectServerTime)

	t6 := time.Now()

	p1.RegisterFile("test.txt")
	p2.RegisterFile("test2.txt")
	p3.RegisterFile("test3.txt")
	p4.RegisterFile("test4.txt")
	p5.RegisterFile("test5.txt")

	t7 := time.Now()
	peerRegisterTime := t7.Sub(t6)
	fmt.Printf("Peer register file time (5): %v", peerRegisterTime)

	p1.ConnectPeer(":60122", 2)

	p1.ConnectPeer(":60123", 3)

	p2.ConnectPeer(":60123", 3)
	p2.ConnectPeer(":60124", 4)
	p2.ConnectPeer(":60125", 5)

	p3.ConnectPeer(":60124", 4)
	p4.ConnectPeer(":60125", 5)

	p1.RequestFile(":60122", 2, "test2.txt")
	p2.RequestFile(":60123", 3, "test3.txt")
	p3.RequestFile(":60124", 4, "test4.txt")
	p4.RequestFile(":60125", 5, "test5.txt")
	p5.RequestFile(":60121", 1, "test.txt")

	p1.SearchForFile("test3.txt")
	p2.SearchForFile("test4.txt")
	p3.SearchForFile("test5.txt")
	p4.SearchForFile("test.txt")
	p5.SearchForFile("test2.txt")
}
