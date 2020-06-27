package main

import (
	"../fileshare"
	"fmt"
	//"time"
)

func main() {
	m := fileshare.MakeSwarmMaster()
	fileshare.Connect()
	m.MasterTest()
}
