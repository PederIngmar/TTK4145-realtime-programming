package main

import (
	"fmt"
	"net"
	"os"
)

func findServerIP() (string, error) {
	port := ":30000" // Listen on port 30000 for broadcast messages
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return "", fmt.Errorf("failed to resolve UDP address: %w", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return "", fmt.Errorf("failed to listen on UDP port: %w", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	fmt.Println("Listening for server broadcasts...")
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read from UDP: %w", err)
	}

	serverIP := string(buffer[:n])
	fmt.Printf("Received broadcast: %s\n", serverIP)
	return serverIP, nil
}

func communicateWithServer(serverIP string, portOffset int) error {
	serverAddr := fmt.Sprintf("%s:%d", serverIP, 20000+portOffset)
	addr, err := net.ResolveUDPAddr("udp4", serverAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve server address: %w", err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	// Send a message
	message := "Hello, Server!"
	fmt.Printf("Sending message: %s\n", message)
	_, err = conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// Receive a response
	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}

	fmt.Printf("Received response: %s\n", string(buffer[:n]))
	return nil
}

func main() {
	serverIP_message, err := findServerIP()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding server IP: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Server IP: %s\n", serverIP_message)

	serverIP := "10.22.23.144"
	workspaceNumber := 9 // Replace with your workspace number

	err = communicateWithServer(serverIP, workspaceNumber)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error communicating with server: %v\n", err)
		os.Exit(1)
	}
}
