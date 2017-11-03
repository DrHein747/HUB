package main

import (
	"fmt"
	"net"
	"os"
)

var clients = make(map[net.Conn]chan<- string)

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

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {

	_, err := conn.WriteToUDP([]byte("From server: Hello I got your mesage "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
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
	ch := make(chan string)

	go handleClient(addchan, ServerConn)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		go sendResponse(ServerConn, addr)
		addchan <- Client{ServerConn, ch}

	}
}

func handleClient(addchan <-chan Client, c net.Conn) {
	fmt.Println(c)
	for {
		select {
		case client := <-addchan:
			fmt.Printf("New client: %v\n", client.conn)
			clients[client.conn] = client.ch
		}
	}
}
