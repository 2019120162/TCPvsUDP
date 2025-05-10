package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

func logClientMetrics(startTime time.Time) {
	var memStats runtime.MemStats
	for {
		time.Sleep(5 * time.Second)

		runtime.ReadMemStats(&memStats)
		elapsed := time.Since(startTime).Seconds()

		fmt.Printf("\n--- CLIENT METRICS ---\n")
		fmt.Printf("Uptime: %.1fs\n", elapsed)
		fmt.Printf("CPU Cores: %d\n", runtime.NumCPU())
		fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
		fmt.Printf("RAM Usage: %.2f MB\n", float64(memStats.Alloc)/1024/1024)
		fmt.Printf("----------------------\n")
	}
}

func main() {
	serverAddr, _ := net.ResolveUDPAddr("udp", "localhost:9001")
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	start := time.Now()
	go logClientMetrics(start)

	// Receive messages
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buffer)
			if err == nil {
				fmt.Print(string(buffer[:n]))
			}
		}
	}()

	// Send messages
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		startSend := time.Now()
		conn.Write([]byte(text))
		fmt.Printf("📤 Message sent in %v\n", time.Since(startSend))
	}
}
