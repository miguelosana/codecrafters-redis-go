package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"example.com/redis/app/redis"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	dataStore := make(map[string]string)

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

		go requestHandler(conn, dataStore)
		log.Printf("connection closed")

	}
}

func requestHandler(conn net.Conn, dataStore map[string]string) {
	log.Print("Handling request")
	decoder := redis.NewDecoder(conn)
	for {
		respValue, _, err := decoder.Decode()
		if err != nil {
			if err == io.EOF {
				conn.Close()
			}
			log.Print(err)
			return

		}
		log.Print("Got something....")
		log.Print(string(respValue.Bytes()))
		if len(respValue.Array()) > 0 {
			arr := respValue.Array()
			if strings.ToUpper(arr[0].String()) == "PING" {
				conn.Write([]byte("+PONG\r\n"))
			} else if strings.ToUpper(arr[0].String()) == "ECHO" {
				var response []byte
				response = append(response, '+')
				response = append(response, arr[1].String()...)
				response = append(response, "\r\n"...)
				_, err := conn.Write(response)
				if err != nil {
					log.Fatal(err)
				}
			} else if strings.ToUpper(arr[0].String()) == "COMMAND" {
				conn.Write([]byte("+OK\r\n"))
			} else if strings.ToUpper(arr[0].String()) == "SET" {
				dataStore[arr[1].String()] = arr[2].String()
				conn.Write([]byte("+OK\r\n"))
			} else if strings.ToUpper(arr[0].String()) == "GET" {
				var response []byte
				response = append(response, '+')
				response = append(response, dataStore[arr[1].String()]...)
				response = append(response, "\r\n"...)
				conn.Write(response)
			}

		}
	}
}
