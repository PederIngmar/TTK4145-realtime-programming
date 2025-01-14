package main

import (
	"fmt"
	"net"
	"sync"
	"time"
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

func send_tcp_server(wg *sync.WaitGroup) error {
	defer wg.Done()
	port := "34933" // 33546
	addr := fmt.Sprintf(":%s", port)
	fmt.Println("TCP address", addr)

	conn, err := net.Dial("tcp", addr)
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

	time.Sleep(1 * time.Second)
	return nil
}

func receive_tcp_server() error {
	port := "34933" // 33546
	addr := fmt.Sprintf(":%s", port)
	fmt.Println("tcp port address", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on TCP port: %w", err)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept connection: %w", err)
	}

	defer conn.Close()

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetNoDelay(false)
	}

	buffer := make([]byte, 1024)
	fmt.Println("Listening for messages... + \x00")
	_, err = conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to read from TCP: %w", err)
	}

	// TCP_NODELAY
	return nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	//_, err := findServerIP()
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Error finding server IP: %v\n", err)
	//	os.Exit(1)
	//}
	//send_tcp_server()
	receive_tcp_server()

	// serverIP := "10.100.23.204"
	//workspaceNumber := 3 // Replace with your workspace number

	// for j := 0; j < 3; j++ {
	// 	err = send_udp_server(serverIP, workspaceNumber)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error communicating with server: %v\n", err)
	// 		os.Exit(1)
	// 	}
	// }

	// for j := 0; j < 3; j++ {
	// 	err = receive_udp_server(workspaceNumber)
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "Error communicating with server: %v\n", err)
	// 		os.Exit(1)
	// 	}
	// }
}
