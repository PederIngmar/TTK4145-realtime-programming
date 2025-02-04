package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const (
	port           = ":8080"       // UDP-port for kommunikasjon mellom prosesser
	heartbeatFreq  = time.Second   // Frekvens for "I am alive"-meldinger
	heartbeatLimit = 3             // Hvor mange "I am alive"-meldinger som kan mistes før backup tar over
)

func main() {
	args := os.Args
	var counter int
	if len(args) > 1 {
		counter, _ = strconv.Atoi(args[1]) // Start counter fra eksisterende verdi
	}

	if len(args) == 1 {
		// Lag backup-prosess
		createBackupProcess(counter)
	}

	go listenForHeartbeat() // Backup lytter etter heartbeat-meldinger

	for {
		counter++
		fmt.Println("Counting:", counter)

		sendHeartbeat() // Hovedprosess sender heartbeat
		time.Sleep(time.Second)

		if len(args) == 1 { // Hovedprosess lager backup
			createBackupProcess(counter)
		}
	}
}

// Lager en ny prosess (backup)
func createBackupProcess(counter int) {
	cmd := exec.Command("gnome-terminal", "--", "go", "run", "main.go", strconv.Itoa(counter))
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to create backup process:", err)
		os.Exit(1)
	}
}

// Sender en UDP-melding som indikerer at prosessen er i live
func sendHeartbeat() {
	conn, err := net.Dial("udp", "localhost"+port)
	if err != nil {
		fmt.Println("Error sending heartbeat:", err)
		return
	}
	defer conn.Close()
	conn.Write([]byte("I am alive"))
}

// Lytter etter heartbeat-meldinger fra hovedprosessen
func listenForHeartbeat() {
	addr, _ := net.ResolveUDPAddr("udp", port)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error listening for heartbeat:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	missedBeats := 0

	for {
		conn.SetReadDeadline(time.Now().Add(heartbeatFreq * heartbeatLimit))
		_, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			missedBeats++
			if missedBeats >= heartbeatLimit {
				fmt.Println("Primary process has died. Taking over...")
				main() // Backup tar over som hovedprosessen
			}
		} else {
			missedBeats = 0 // Reset på hver mottatt melding
		}
	}
}
