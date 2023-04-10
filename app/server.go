package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"example.com/redis/app/redis"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go requestHandler(conn)

	}
}

func requestHandler(conn net.Conn) {
	log.Print("Handling request")
	decoder := redis.NewDecoder(conn)
	for {
		respValue, _, err := decoder.Decode()
		if err != nil {
			log.Fatal(err)

		}
		log.Print("Got something....")
		log.Print(string(respValue.Bytes()))
		if len(respValue.Array()) > 0 {
			arr := respValue.Array()
			if strings.ToUpper(arr[0].String()) == "PING" {
				conn.Write([]byte("+PONG\r\n"))
			}

		}
	}
}
