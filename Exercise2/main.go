package main

import (
	"fmt"
	"net"
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
	//serverIP := "10.100.23.204"
}

func send_udp_server(serverIP string, portOffset int) error {
	serverAddr := fmt.Sprintf("%s:%d", serverIP, 20000+portOffset) // formats the server address from the IP and port
	addr, err := net.ResolveUDPAddr("udp", serverAddr)             //
	if err != nil {
		return fmt.Errorf("failed to resolve server address: %w", err)
	} else {
		fmt.Println("udp address", addr)
	}

	conn, err := net.DialUDP("udp", nil, addr) // DialUDP instead of ListenUDP
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	// Send a message
	message := "Skibidi rizzzz, aha!"
	fmt.Printf("Sending message: %s\n", message)
	code, err := conn.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	} else {
		fmt.Println("status code:", code)
	}
	return nil
}

func receive_udp_server(portOffset int) error {
	port := 20000 + portOffset
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %w", err)
	} else {
		fmt.Println("udp address:", addr)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP port: %w", err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	fmt.Println("Listening for messages...")
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return fmt.Errorf("failed to read from UDP: %w", err)
	}

	message := string(buffer[:n])
	fmt.Printf("Received message: %s\n", message)
	return nil
}

//findServerIP()
//serverIP := "10.100.23.204"

