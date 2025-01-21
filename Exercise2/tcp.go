package main

import (
	"fmt"
	"net"
	"time"
)

func tcp_client() {
	addr := "10.100.23.204:33546"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	message := "Who is this diva\x00"
	conn.Write([]byte(message))
	fmt.Println("Message sendt: ", message)

	time.Sleep(1 * time.Second)

	buffer := make([]byte, 1024)
	fmt.Println("Listening for messages...")
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Message received: %s \n", string(buffer[:n]))

}

func tcp_server() {
	addr := "10.100.23.204:33546"

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	message := "Connect to: 10.100.23.13:33546\x00"
	conn.Write([]byte(message))
	fmt.Println("Connection message:", message)

	l, err := net.Listen("tcp", ":33546")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		conn2, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleClient(conn2)
	}

}

func handleClient(conn net.Conn) {
	defer conn.Close()

	for {

		message := "Who is this diva\x00"
		conn.Write([]byte(message))
		fmt.Println("Message sendt: ", message)

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Message received: %s \n", string(buffer[:n]))

		time.Sleep(1 * time.Second)
	}

}

func main() {
	//serverIP := "10.100.23.204:34933"
	//tcp_client()
	tcp_server()
}
