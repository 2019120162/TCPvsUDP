package main

import (
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"
)

var (
	clients       = make(map[string]*net.UDPAddr)
	clientNames   = make(map[string]string)
	totalMessages int
	totalClients  int
	startTime     = time.Now()
)

func logMetrics() {
	var memStats runtime.MemStats
	for {
		time.Sleep(5 * time.Second)
		runtime.ReadMemStats(&memStats)
		elapsed := time.Since(startTime).Seconds()

		fmt.Printf("\n--- SERVER METRICS ---\n")
		fmt.Printf("Uptime: %.1fs\n", elapsed)
		fmt.Printf("Total Clients: %d\n", totalClients)
		fmt.Printf("Messages: %d (%.2f msg/sec)\n", totalMessages, float64(totalMessages)/elapsed)
		fmt.Printf("CPU Cores: %d\n", runtime.NumCPU())
		fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
		fmt.Printf("RAM Usage: %.2f MB\n", float64(memStats.Alloc)/1024/1024)
		fmt.Printf("----------------------\n")
	}
}

func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":9001")
	conn, _ := net.ListenUDP("udp", addr)
	defer conn.Close()

	fmt.Println("✅ UDP server started on :9001")
	go logMetrics()

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, _ := conn.ReadFromUDP(buffer)
		clientID := clientAddr.String()
		msg := strings.TrimSpace(string(buffer[:n]))

		if _, exists := clients[clientID]; !exists {
			clients[clientID] = clientAddr
			clientNames[clientID] = msg // first message is name
			totalClients++
			fmt.Printf("👤 %s joined as '%s'\n", clientID, msg)
			continue
		}

		senderName := clientNames[clientID]
		timestamp := time.Now().Format("15:04:05")
		fullMsg := fmt.Sprintf("[%s] %s: %s", timestamp, senderName, msg)

		for id, addr := range clients {
			if id != clientID {
				conn.WriteToUDP([]byte(fullMsg), addr)
			}
		}

		totalMessages++
	}
}
