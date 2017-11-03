package main

import (
	"fmt"
	"net"
	"os"
)

type Client struct {
	conn net.Conn
	ch   chan<- string
}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {
	/* Lets prepare a address at any address at port 8080*/
	ServerAddr, err := net.ResolveUDPAddr("udp", ":8080")
	CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	addchan := make(chan Client)
	msgchan := make(chan string)
	ch := make(chan string)

	go handleMessage(addchan, ServerConn, msgchan)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		msgchan <- string(buf[0:n])
		addchan <- Client{ServerConn, ch}
	}
}

func handleMessage(addchan <-chan Client, c net.Conn, msgchan <-chan string) {
	var clients = make(map[net.Conn]chan<- string)
	for {
		select {
		case msg := <-msgchan:
			fmt.Printf("new message: %s\n", msg)
			for _, ch := range clients {
				go func(mch chan<- string) { mch <- "From server: Hello I got your mesage " + msg }(ch)
			}
		case client := <-addchan:
			fmt.Printf("New client: %v\n", client.conn)
			clients[client.conn] = client.ch

		}
	}
}
