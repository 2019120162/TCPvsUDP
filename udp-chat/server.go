package main

import (
	"expvar"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	clients         = make(map[string]*net.UDPAddr)
	totalMessages   = expvar.NewInt("totalMessages")
	totalClients    = expvar.NewInt("totalClients")
	startTime       = time.Now()
)

func logMetrics() {
	var memStats runtime.MemStats
	for {
		time.Sleep(5 * time.Second)

		runtime.ReadMemStats(&memStats)
		elapsed := time.Since(startTime).Seconds()

		fmt.Printf("\n--- SERVER METRICS ---\n")
		fmt.Printf("Uptime: %.1fs\n", elapsed)
		fmt.Printf("Total Clients: %d\n", totalClients.Value())
		fmt.Printf("Messages: %d (%.2f msg/sec)\n", totalMessages.Value(), float64(totalMessages.Value())/elapsed)
		fmt.Printf("CPU Cores: %d\n", runtime.NumCPU())
		fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
		fmt.Printf("RAM Usage: %.2f MB\n", float64(memStats.Alloc)/1024/1024)
		fmt.Printf("----------------------\n")
	}
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":9001")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Listen error:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("UDP server started on :9001")

	go logMetrics()

	buffer := make([]byte, 1024)

	for {
		n, clientAddr, _ := conn.ReadFromUDP(buffer)
		msg := strings.TrimSpace(string(buffer[:n]))

		if _, exists := clients[clientAddr.String()]; !exists {
			clients[clientAddr.String()] = clientAddr
			totalClients.Add(1)
			fmt.Printf("New client joined: %s\n", clientAddr.String())
		}

		timestamp := time.Now().Format("15:04:05")
		fullMsg := fmt.Sprintf("[%s] %s: %s", timestamp, clientAddr, msg)

		for _, c := range clients {
			if c.String() != clientAddr.String() {
				conn.WriteToUDP([]byte(fullMsg), c)
			}
		}

		totalMessages.Add(1)
	}
}
