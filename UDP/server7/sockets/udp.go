package sockets

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

type Client struct {
	ID   int
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

var clients = make(map[int]map[*net.UDPAddr]*net.UDPConn)
var clientinfo = make(map[*net.UDPAddr]*net.UDPConn)
var messages = make(map[*net.UDPAddr]string)
var ID = 1

func UdpInit() {
	/* Lets prepare a address at any address at port 8080*/
	ServerAddr, err := net.ResolveUDPAddr("udp", ":8554")
	CheckError(err)
	fmt.Println("UDP-Server started")

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
		addchan <- Client{ID, ServerConn, addr}
		msgchan <- Message{addr, data}

	}
}

func handleMessage(addchan <-chan Client, conn *net.UDPConn, msgchan <-chan Message) {
	for {
		select {
		case message := <-msgchan:
			messages[message.addrM] = message.data
			for _, value := range clients {
				go sentToAllOthers(message, value)
			}
		case client := <-addchan:
			clientfound := false
			for _, value := range clients {
				for addr := range value {
					if addr.String() == client.addr.String() {
						clientfound = true
					}
				}
			}
			if clientfound == false {
				ID = ID + 1
				clients[client.ID] = clientinfo
				clients[client.ID][client.addr] = conn
				fmt.Println(client.ID)
				go sendToOneCLient(client.ID, client.addr)
			}
		}
	}
}

func sentToAllOthers(message Message, value map[*net.UDPAddr]*net.UDPConn) {
	for addr, conn := range value {
		if message.addrM.String() == addr.String() {
			return
		}
		fmt.Println(message)
		conn.WriteToUDP([]byte(message.data), addr)
	}
}

func sendToOneCLient(ClientID int, Clientaddr *net.UDPAddr) {
	fmt.Println(ClientID)
	for _, value := range clients {
		for addr, conn := range value {
			if addr == Clientaddr {
				fmt.Println(strconv.Itoa(ClientID))
				message := "socketrec.setID(" + strconv.Itoa(ClientID) + ")|"
				conn.WriteToUDP([]byte(message), addr)
				fmt.Println(len(clients))
				return
			}
		}
	}
}
