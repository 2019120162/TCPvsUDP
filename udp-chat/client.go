package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr, _ := net.ResolveUDPAddr("udp", "localhost:9001")
	conn, _ := net.DialUDP("udp", nil, serverAddr)
	defer conn.Close()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, _, _ := conn.ReadFromUDP(buf)
			fmt.Print(string(buf[:n]))
		}
	}()

	scanner := bufio.NewReader(os.Stdin)
	for {
		text, _ := scanner.ReadString('\n')
		conn.Write([]byte(text))
	}
}
