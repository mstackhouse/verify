package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
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

	// Get User and Path values from command line arguments
	user := os.Args[1]
	path := os.Args[2]

	// Normalize the path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	msg := Message{
		User: user,
		Path: absPath,
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
