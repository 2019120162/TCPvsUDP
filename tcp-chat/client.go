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
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("❌ Connection error:", err)
		return
	}
	defer conn.Close()

	start := time.Now()
	go logClientMetrics(start)

	// Receive name prompt
	serverMsg, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print(serverMsg)

	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	conn.Write([]byte(name))

	// Receive messages
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

	// Send messages with latency logging
	for {
		text, _ := reader.ReadString('\n')
		startSend := time.Now()
		conn.Write([]byte(text))
		fmt.Printf("📤 Message sent in %v\n", time.Since(startSend))
	}
}
