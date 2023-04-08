package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	requestHandler(conn)
}

func requestHandler(conn net.Conn) {
	log.Print("Handling request")
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Read %d bytes", n)
	command := string(buffer[:n-1])
	if command == "ping" {
		log.Println("PONG")
		conn.Write([]byte("+PONG\r\n"))
	}
	log.Printf("%v", command)
	conn.Close()
}
