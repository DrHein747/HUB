package sockets

import (
	"fmt"
	"log"
	"net"
)

type TCPClient struct {
	iDTCP int
	conn  net.Conn
}

type TCPMessage struct {
	conn net.Conn
	data []byte
}

var iDTCP = 1

// Called in udp.go
func TCPInit() {

	ln, err := net.Listen("tcp", ":554")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TCPS-Server started")
	msgchan := make(chan TCPMessage)
	addchan := make(chan TCPClient)

	go handleMessages(msgchan, addchan)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, msgchan, addchan)
	}
}

func handleConnection(c net.Conn, msgchan chan<- TCPMessage, addchan chan<- TCPClient) {

	var buf [1024]byte
	addchan <- TCPClient{iDTCP, c}

	for {
		n, err := c.Read(buf[0:])
		if err != nil {
			return
		}
		data := buf[0:n]
		msgchan <- TCPMessage{c, data}
	}
}

func handleMessages(msgchan <-chan TCPMessage, addchan <-chan TCPClient) {

	clients := make(map[int]net.Conn)
	fmt.Println(len(clients))
	for {
		select {
		case message := <-msgchan:
			go sendToAll(message, clients)
		case client := <-addchan:
			clientfound := false
			for id, conn := range clients {
				if conn == client.conn {
					clientfound = true
					fmt.Println("Client exists with the ID: ", id)
				}
			}
			if clientfound == false {
				iDTCP = iDTCP + 1
				clients[client.iDTCP] = client.conn
				fmt.Println("Client ", client.iDTCP, "is connected!")
			}
		}
	}
}

func sendToAll(message TCPMessage, allClient map[int]net.Conn) {
	for _, conn := range allClient {
		conn.Write(message.data)
	}
}
