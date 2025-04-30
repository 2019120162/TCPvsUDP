package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server.")

	// Handle server prompt for name
	serverMsg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read server prompt:", err)
		return
	}
	fmt.Print(serverMsg) // Expect: "Enter your name: "

	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read name input:", err)
		return
	}
	conn.Write([]byte(name))

	// Start goroutine to listen for server messages
	go func() {
		for {
			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("🔌 Disconnected from server.")
				return
			}
			fmt.Print(msg)
		}
	}()

	// Send messages
	for {
		text, _ := reader.ReadString('\n')
		start := time.Now()
		conn.Write([]byte(text))
		fmt.Printf("📤 Sent in %v\n", time.Since(start))
	}
}
