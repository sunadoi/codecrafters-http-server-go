package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	req := make([]byte, 1024)
	if _, err := conn.Read(req); err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	path := strings.Split(string(req), " ")[1]

	switch {
	case strings.HasPrefix(path, "/echo"):
		body := strings.Split(path, "/")[2]
		writeResponse(conn, fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), string(body)))
	case path == "/user-agent":
		uaHeader := strings.Split(string(req), "\r\n")[2]
		ua := strings.Split(uaHeader, ": ")[1]
		writeResponse(conn, fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(ua), ua))
	case path == "/":
		writeResponse(conn, "HTTP/1.1 200 OK\r\n\r\n")
	default:
		writeResponse(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
	}
}

func writeResponse(conn net.Conn, res string) {
	if _, err := conn.Write([]byte(res)); err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
