package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients   = make(map[net.Conn]string)
	mutex     = &sync.Mutex{}
	broadcast = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	go handleBroadcast()

	fmt.Println("TCP server started on :9000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Enter your name: "))
	nameBuf, _ := bufio.NewReader(conn).ReadString('\n')
	name := strings.TrimSpace(nameBuf)

	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()

	broadcast <- fmt.Sprintf("%s joined the chat", name)

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		broadcast <- fmt.Sprintf("%s: %s", name, msg)
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()

	broadcast <- fmt.Sprintf("%s left the chat", name)
}

func handleBroadcast() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			client.Write([]byte(msg))
		}
		mutex.Unlock()
	}
}
