package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
)

type Message struct {
	User string `json:"text"`
	Path string `json:"path"`
}

type Response struct {
	User      string `json:"text"`
	Path      string `json:"path"`
	HasAccess bool   `json:"hasAccess"`
}

func main() {
	socketPath := "/tmp/unix_socket"

	// Error if the program is not run as root
	if os.Getuid() != 0 {
		fmt.Println("This program must be run as root")
		return
	}

	// Remove the socket if it already exists
	if err := os.RemoveAll(socketPath); err != nil {
		fmt.Println("Error removing existing socket:", err)
		return
	}

	// Create a Unix domain socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Error creating Unix socket:", err)
		return
	}
	defer listener.Close()

	if err := os.Chmod(socketPath, 0666); err != nil {
		fmt.Println("Error setting socket permissions:", err)
		return
	}

	fmt.Println("Server is listening on", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	var msg Message
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&msg); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Check that msg.Path is a path that exists
	if _, err := os.Stat(msg.Path); err != nil {
		fmt.Println("Error validating path:", err)
		return
	}

	// system call to check if the user has access to the path
	// Execute the system command "test" to check if the user has access to the path
	// TODO: Toggle read/write check as optional input
	cmd := exec.Command("su", msg.User, "-c", "test -w "+msg.Path)
	out, _ := cmd.CombinedOutput()

	var response Response
	response.User = msg.User
	response.Path = msg.Path

	// If the command exits with a non-zero status code, the user does not have access
	if cmd.ProcessState.ExitCode() != 0 {
		response.HasAccess = false
		fmt.Println("User does not have access to", msg.Path)
		if len(out) > 0 {
			fmt.Println("Output:", string(out))
		}
	} else {
		response.HasAccess = true
		fmt.Println("User has access to", msg.Path)
	}

	// Send the response back to the client
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&response); err != nil {
		fmt.Println("Error encoding response:", err)
	}

}
