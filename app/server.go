package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
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

	for {
		go func() {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err.Error())
				os.Exit(1)
			}
			defer conn.Close()

			req, err := http.ReadRequest(bufio.NewReader(conn))
			if err != nil {
				fmt.Println("Error reading request. ", err.Error())
				return
			}

			path := req.URL.Path

			switch {
			case strings.HasPrefix(path, "/echo"):
				body := strings.Split(path, "/")[2]
				writeResponse(conn, fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(body), string(body)))
			case strings.HasPrefix(path, "/files/"):
				fileName := strings.Split(path, "/")[2]
				dir := os.Args[2]
				filePath := fmt.Sprintf("%s/%s", dir, fileName)
				file, err := os.Open(filePath)
				if err != nil {
					if os.IsNotExist(err) {
						writeResponse(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
					} else {
						writeResponse(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
					}
				}
				buf := make([]byte, 1024)
				n, err := file.Read(buf)
				if err != nil {
					writeResponse(conn, "HTTP/1.1 500 Internal Server Error\r\n\r\n")
				}
				writeResponse(conn, fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", n, string(buf)))
			case path == "/user-agent":
				ua := req.UserAgent()
				writeResponse(conn, fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(ua), ua))
			case path == "/":
				writeResponse(conn, "HTTP/1.1 200 OK\r\n\r\n")
			default:
				writeResponse(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
			}
		}()
	}
}

func writeResponse(conn net.Conn, res string) {
	if _, err := conn.Write([]byte(res)); err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
