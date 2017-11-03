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

type Message struct {
	addrM *net.UDPAddr
	data  string
}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

var clients = make(map[*net.UDPAddr]*net.UDPConn)
var messages = make(map[*net.UDPAddr]string)

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
	msgchan := make(chan Message)

	go handleMessage(addchan, ServerConn, msgchan)

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		data := string(buf[0:n])
		addchan <- Client{ServerConn, addr}
		msgchan <- Message{addr, data}
	}
}

func handleMessage(addchan <-chan Client, conn *net.UDPConn, msgchan <-chan Message) {

	for {
		select {
		case message := <-msgchan:
			messages[message.addrM] = message.data
			for addr, conn := range clients {
				go sentToAllOthers(message, conn, addr)
			}
		case client := <-addchan:
			clientfound := false
			for addr, _ := range clients {
				if addr.String() == client.addr.String() {
					clientfound = true
				}
			}
			if clientfound == false {
				clients[client.addr] = conn
				message := ("Avatar:")
				go sendToAll(message)
			}
		}
	}
}

func sentToAllOthers(message Message, conn *net.UDPConn, addr *net.UDPAddr) {
	if message.addrM.String() == addr.String() {
		return
	}
	conn.WriteToUDP([]byte(message.data), addr)
}
func sendToAll(message string) {
	for addr, conn := range clients {
		conn.WriteToUDP([]byte(message), addr)
	}
}
