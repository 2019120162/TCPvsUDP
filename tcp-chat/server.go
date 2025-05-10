package main

import (
	"bufio"
	"expvar"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	clients   = make(map[net.Conn]string)
	mutex     = &sync.Mutex{}
	broadcast = make(chan string)

	totalMessages = expvar.NewInt("totalMessages")
	totalClients  = expvar.NewInt("totalClients")
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
		fmt.Printf("Clients: %d\n", totalClients.Value())
		fmt.Printf("Messages: %d (%.2f msg/sec)\n", totalMessages.Value(), float64(totalMessages.Value())/elapsed)
		fmt.Printf("CPU Cores: %d\n", runtime.NumCPU())
		fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
		fmt.Printf("RAM Usage: %.2f MB\n", float64(memStats.Alloc)/1024/1024)
		fmt.Printf("----------------------\n")
	}
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	go logMetrics()
	go handleBroadcast()

	fmt.Println(" TCP server started on :9000")

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
	totalClients.Add(1)
	mutex.Unlock()

	broadcast <- fmt.Sprintf("📢 %s joined the chat\n", name)

	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		totalMessages.Add(1)
		broadcast <- fmt.Sprintf("%s: %s", name, msg)
	}

	mutex.Lock()
	delete(clients, conn)
	mutex.Unlock()

	broadcast <- fmt.Sprintf("🚪 %s left the chat\n", name)
	totalClients.Add(-1)
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
