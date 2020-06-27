package main

import (
	"../fileshare"
	"fmt"
	//"time"
)

func main() {
	p := fileshare.Make()
	fmt.Printf("success\n")
	p.Hello()
}
