package main

import (
	"log"
	"net"
)

func readWrite(conn *net.UDPConn, clients map[string]string, slice []*net.UDPAddr) []*net.UDPAddr {
	//simple Read and Write to client
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	clientMessage := string(buffer[:n])
	clientAddr := addr.String()

	dontPrint := false // to prevent the printing the welcome message whenever someone joins

	if val, ok := clients[clientAddr]; ok { // Check to see if the client is in the map
		if clientMessage == "quit" { // Removing client from the list is requested
			delete(clients, clientAddr)

			for i := 0; i < len(slice); i++ {
				if slice[i] == addr {
					slice[i] = slice[len(slice)-1]
					return slice[:len(slice)-1]
				}
			}
			log.Println("******", val, "has disconnected ******")
			dontPrint = true
		} else {
			log.Println(val, ":", clientMessage) // Printing the message to the server
		}
	} else { // If not in the map then add the client
		clients[clientAddr] = clientMessage
		slice = append(slice, addr)
		dontPrint = true
		log.Println("******", clientMessage, "has established a connection ******")
	}

	if err != nil {
		log.Println(err)
	}

	// Displaying the messages to the clients
	message := []byte(clients[clientAddr] + " : " + clientMessage)
	for i := 0; i < len(slice); i++ {
		if dontPrint == false {
			_, err = conn.WriteToUDP(message, slice[i])
		}
	}

	if err != nil {
		log.Println(err)

	}

	return slice

}

func main() {

	clients := make(map[string]string) // all connected clients
	slice := []*net.UDPAddr{}

	// listen to incoming udp packets
	dst, err := net.ResolveUDPAddr("udp4", "localhost:8000")
	if err != nil {
		log.Fatal(err)
		log.Println("Error")
	}

	ln, err := net.ListenUDP("udp", dst)
	if err != nil {
		log.Fatal(err)
		log.Println("Error")
	}
	defer ln.Close()

	for {
		// Wait for client to send message
		slice = readWrite(ln, clients, slice)
	}
}
