package main

import (
	"testworkspace/UDP/server7/sockets"
)

func main() {
	go sockets.UdpInit()
	sockets.TCPInit()
}
