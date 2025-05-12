# TCP vs UDP Chat App (Go)

Simple terminal chat programs written in Go to compare **TCP** (reliable) and **UDP** (fast) transport layers.  
Both versions log basic metrics every 5 s: CPU cores, goroutines, RAM (MB), uptime, message rate, and message-send latency.

---

## Run – TCP version

## 1. start server
cd tcp-chat
go run server.go

## 2. in another terminal start client(s)
cd tcp-chat
go run client.go        # enter your name when prompted

# Run – UDP version
## 1. start server
cd udp-chat
go run server.go

## 2. in another terminal start client(s)
cd udp-chat
go run client.go        # you’ll be asked for a name first


# What does the code do?

| File                 | Role                                                  |
| -------------------- | ----------------------------------------------------- |
| `tcp-chat/server.go` | TCP server, broadcasts to all, logs metrics           |
| `tcp-chat/client.go` | TCP client, shows prompt, logs latency & RAM          |
| `udp-chat/server.go` | UDP server, first packet is client name, logs metrics |
| `udp-chat/client.go` | UDP client, sends name first, logs latency & RAM      |


# Video

Link here: https://youtu.be/ouP9oG6HVfU
