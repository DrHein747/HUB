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

var clients = make(map[*net.UDPAddr]*net.UDPConn)

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

	go handleMessage(addchan, ServerConn, msgchan)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		addchan <- Client{ServerConn, addr}
		msgchan <- string(buf[0:n])
	}

}

func handleMessage(addchan <-chan Client, conn *net.UDPConn, msgchan <-chan string) {

	for {
		select {
		case msg := <-msgchan:
			fmt.Printf("new message: %s\n", msg)
			for addr, conn := range clients {
				conn.WriteToUDP([]byte(msg), addr)
			}
		case client := <-addchan:
			clientfound := false
			fmt.Println("Zahl", int(len(clients)))
			for adr := range clients {
				if adr.String() == client.addr.String() {
					fmt.Println("Client ist inder liste")
					clientfound = true
				}
			}
			if clientfound == false {
				fmt.Println("Client ist nicht in der liste")
				clients[client.addr] = conn
			}
		}
	}
}
