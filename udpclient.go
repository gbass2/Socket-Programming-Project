package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	//Connect udp
	address := "localhost:8000"
	RemoteAddr, err := net.ResolveUDPAddr("udp", address)
	conn, err := net.DialUDP("udp", nil, RemoteAddr)
	if err != nil {
		log.Println("Error")
	}

	log.Printf("Established connection to %s \n", address)
	log.Printf("Remote UDP address : %s \n", conn.RemoteAddr().String())
	log.Printf("Local UDP client address : %s \n", conn.LocalAddr().String())
	fmt.Print("\n")

	defer conn.Close()

	fmt.Print("Welcome to the server \nEnter 'quit' to leave the server", "\n")
	fmt.Print("\nWhat is your name : ")
	name := ""

	for {
		// Getting user's message from input
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		if name == "" {
			name = scanner.Text()
			go readFromServer(name, conn)
		}
		clientMessage := scanner.Text()
		message := []byte(clientMessage)
		conn.Write([]byte(message)) // Sending message to the server

		if clientMessage == "quit" { // Suspending the program upon leaving
			os.Exit(0)
		}
	}
}

func readFromServer(name string, conn *net.UDPConn) { // Reading a message from the server
	for {
		buffer := make([]byte, 1024)
		n, _, _ := conn.ReadFromUDP(buffer)
		messageFromServer := string(buffer[:n])
		words := strings.Fields(messageFromServer)

		// We don't want to print the this client's message
		if words[0] != name {
			fmt.Print(messageFromServer, "\n")
		}
	}
}
