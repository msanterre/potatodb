package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func HandleCommand(conn net.Conn, body string) {
	var response string
	parts := strings.SplitN(body, " ", 3)

	if len(parts) >= 2 {
		command := strings.ToLower(parts[0])
		key := parts[1]

		switch command {
		case "set":
			if len(parts) > 2 {
				arg := parts[2]
				DB.Data[key] = arg
				response = "1"
			} else {
				response = "!ERR: missing argument for 'set'"
			}
		case "get":
			data := DB.Data[key]
			response = data
		case "delete":
			delete(DB.Data, key)
			response = "1"
		default:
			log.Println("Invalid command: ", command)
			response = "!ERR: invalid command"
		}

		_, err := conn.Write([]byte(response + "\n"))
		if err != nil {
			log.Printf("ERROR (w): ")
			log.Println(err)
		}
	}
}

func HandleConnection(connection net.Conn) {
	reader := bufio.NewReader(connection)

	for {
		body, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println(err)
			connection.Close()
			return
		}

		cleanBody := strings.TrimSpace(string(body))
		HandleCommand(connection, cleanBody)
	}
}

func SocketStart(listenAddress string) {
	server, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Println(err)
			}
			go HandleConnection(conn)
		}
	}()

	log.Println("Socket listening on", listenAddress)
}
