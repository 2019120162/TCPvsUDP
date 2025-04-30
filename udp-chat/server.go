package main

import (
	"fmt"
	"net"
	"strings"
)

var clients = make(map[string]*net.UDPAddr)

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":9001")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	buf := make([]byte, 1024)

	fmt.Println("UDP server started on :9001")

	for {
		n, clientAddr, _ := conn.ReadFromUDP(buf)
		msg := strings.TrimSpace(string(buf[:n]))

		if _, exists := clients[clientAddr.String()]; !exists {
			clients[clientAddr.String()] = clientAddr
			fmt.Printf("%s joined\n", clientAddr)
		}

		for _, c := range clients {
			if c.String() != clientAddr.String() {
				conn.WriteToUDP([]byte(fmt.Sprintf("%s: %s", clientAddr, msg)), c)
			}
		}
	}
}
