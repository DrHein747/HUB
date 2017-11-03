package main

import (
	"fmt"
	"net"
	"os"
)

type Client struct {
	conn *net.UDPConn
	addr *net.UDPAddr
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
	fmt.Println(ServerConn)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)
	addchan := make(chan Client)
	msgchan := make(chan string)
	rmchan := make(chan net.Conn)

	go handleMessage(addchan, ServerConn, msgchan, rmchan)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		msgchan <- string(buf[0:n])
		addchan <- Client{ServerConn, addr}
	}

}

func handleMessage(addchan <-chan Client, conn *net.UDPConn, msgchan <-chan string, rmchan <-chan net.Conn) {

	var clients = make(map[*net.UDPAddr]*net.UDPConn)

	for {
		select {
		case msg := <-msgchan:
			fmt.Printf("new message: %s\n", msg)
			for addr, conn := range clients {
				conn.WriteToUDP([]byte(msg), addr)
			}
		case client := <-addchan:
			clients[client.addr] = client.conn
			fmt.Println("New client", int(len(clients)), "is connected!")
		case conn := <-rmchan:
			fmt.Printf("Client disconnects: %v\n", conn)
		}
	}
}
