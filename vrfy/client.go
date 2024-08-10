package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	User      string `json:"text"`
	Path      string `json:"path"`
	HasAccess bool   `json:"hasAccess"`
}

type Response struct {
	User      string `json:"text"`
	Path      string `json:"path"`
	HasAccess bool   `json:"hasAccess"`
}

func main() {
	socketPath := "/tmp/unix_socket"

	// Create a Unix domain socket
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		fmt.Println("Error connecting to Unix socket:", err)
		return
	}
	defer conn.Close()

	msg := Message{
		User: "testuser",
		Path: "/home/mstack/go/src/github.com/mstackhouse/verify/testfile.txt",
	}

	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&msg); err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	// Receive response from the socket
	decoder := json.NewDecoder(conn)
	var response Message
	if err := decoder.Decode(&response); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println("Response received:", response)

}
