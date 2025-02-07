package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	UDP_IP           = "127.0.0.1"
	UDP_PORT         = 5005
	ALIVE_INTERVAL   = 1 * time.Second   // Hvor ofte alive-meldinger sendes
	TIMEOUT_THRESHOLD = 3 * time.Second  // Hvor lenge backup venter før den promoterer seg selv
)

// spawnBackup starter en ny backup-prosess i et eget gnome-terminal-vindu
func spawnBackup() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("[PRIMARY] Feil ved henting av eksekverbar sti:", err)
		return
	}

	// Bygg kommandoen: gnome-terminal -- <eksekverbar sti> backup
	cmd := exec.Command("gnome-terminal", "--", exePath, "backup")
	if err := cmd.Start(); err != nil {
		fmt.Println("[PRIMARY] Feil ved spawning av backup:", err)
	} else {
		fmt.Println("[PRIMARY] Spawnet backup-prosess.")
	}
}

// primaryLoop er den primære prosessen som teller og sender alive-meldinger via UDP
func primaryLoop(startCount int) {
	count := startCount

	// Opprett en UDP-tilkobling for sending
	addr := net.UDPAddr{
		IP:   net.ParseIP(UDP_IP),
		Port: UDP_PORT,
	}
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		fmt.Println("[PRIMARY] Feil ved oppretting av UDP-forbindelse:", err)
		return
	}
	defer conn.Close()

	// Spawn backup-prosess ved oppstart
	spawnBackup()

	fmt.Println("[PRIMARY] Starter opp-telling fra", count)
	for {
		fmt.Println(count)
		msg := fmt.Sprintf("ALIVE:%d", count)
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("[PRIMARY] Feil ved sending av heartbeat:", err)
		}
		count++
		time.Sleep(ALIVE_INTERVAL)
	}
}

// backupLoop lytter etter alive-meldinger og promoterer seg hvis primæren stopper å sende
func backupLoop() {
	lastCount := 0

	// Lytt på UDP-porten for alive-meldinger
	addr := net.UDPAddr{
		IP:   net.ParseIP(UDP_IP),
		Port: UDP_PORT,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("[BACKUP] Feil ved binding til UDP-port:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("[BACKUP] Starter backup-modus. Venter på alive-meldinger fra primær... 👀")
	lastTime := time.Now()

	buf := make([]byte, 1024)
	for {
		// Sett en read deadline for å sjekke for timeout
		conn.SetReadDeadline(time.Now().Add(ALIVE_INTERVAL))
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Sjekk om timeout har overskredet terskelen
				if time.Since(lastTime) > TIMEOUT_THRESHOLD {
					fmt.Printf("[BACKUP] Ingen heartbeat mottatt på over %.0f sekunder.\n", TIMEOUT_THRESHOLD.Seconds())
					fmt.Printf("[BACKUP] Promoterer backup til primær med start-teller %d.\n", lastCount+1)
					conn.Close() // Lukk socketen slik at ny backup kan binde seg til porten
					primaryLoop(lastCount + 1)
					return
				}
				// Timeout, men ikke lenge nok til å promotere – fortsett lyttingen
				continue
			} else {
				fmt.Println("[BACKUP] Uventet feil under lesing:", err)
				continue
			}
		}

		// Melding mottatt
		msg := string(buf[:n])
		if strings.HasPrefix(msg, "ALIVE:") {
			parts := strings.Split(msg, ":")
			if len(parts) == 2 {
				num, err := strconv.Atoi(parts[1])
				if err == nil {
					lastCount = num
					lastTime = time.Now()
					// Debug: fmt.Printf("[BACKUP] Mottok heartbeat: %d\n", num)
				}
			}
		}
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "backup" {
		fmt.Println("[SYSTEM] Starter i BACKUP-modus. 👀")
		backupLoop()
	} else {
		fmt.Println("[SYSTEM] Starter i PRIMARY-modus. 🚀")
		primaryLoop(1)
	}
}
