package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

var clients sync.Map

func main() {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("TCP server started on :9000")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	conn.Write([]byte("Enter your name: \n")) // Prompt client for name

	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read name:", err)
		return
	}
	name = strings.TrimSpace(name)
	clients.Store(conn, name)
	fmt.Printf("👤 %s joined the chat\n", name)

	broadcast(fmt.Sprintf("!!! %s joined the chat !!!\n", name), conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("!!! %s disconnected !!!\n", name)
			clients.Delete(conn)
			broadcast(fmt.Sprintf("!!! %s left the chat !!!\n", name), conn)
			return
		}
		broadcast(fmt.Sprintf("%s: %s", name, message), conn)
	}
}

func broadcast(message string, sender net.Conn) {
	clients.Range(func(key, value interface{}) bool {
		conn := key.(net.Conn)
		if conn != sender {
			conn.Write([]byte(message))
		}
		return true
	})
}
