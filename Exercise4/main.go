package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	"strings"
	"strconv"
)

const (
	heartbeatFile     = "heartbeat.txt" // Filen som brukes for "heartbeat"
	heartbeatInterval = 1 * time.Second         // Primæren oppdaterer heartbeat hvert sekund
	timeoutDuration   = 2500 * time.Millisecond // Backup venter 2.5 sekunder før den antar at primæren er død
	backupCheckDelay  = 500 * time.Millisecond  // Backup sjekker heartbeat hvert 500ms
)

// spawnBackup spawner en backup-instans i et nytt terminalvindu ved hjelp av gnome-terminal
func spawnBackup() {
	fmt.Println("Spawner backup")
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Feil ved å hente kjørbar bane:", err)
		return
	}

	// Start backup: Kjør det samme programmet med argumentet "--backup"
	cmd := exec.Command("gnome-terminal", "--", exePath, "--backup")
	if err := cmd.Start(); err != nil {
		fmt.Println("Feil ved spawning av backup:", err)
	}
}

// writeHeartbeat skriver ut den nåværende telleverdien og tidsstempelet til heartbeat-filen.
func writeHeartbeat(counter int) {
	// Format: "counter timestamp"
	content := fmt.Sprintf("%d %f\n", counter, float64(time.Now().UnixNano())/1e9)
	err := os.WriteFile(heartbeatFile, []byte(content), 0644)
	if err != nil {
		fmt.Println("Feil ved skriving av heartbeat:", err)
	}
}

// readHeartbeat leser heartbeat-filen og returnerer (counter, timestamp)
func readHeartbeat() (int, time.Time, error) {
	data, err := os.ReadFile(heartbeatFile)
	if err != nil {
		return 0, time.Time{}, err
	}
	parts := strings.Fields(string(data))
	if len(parts) < 2 {
		return 0, time.Time{}, fmt.Errorf("ugyldig format i heartbeat-filen")
	}

	counter, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, time.Time{}, err
	}

	tsFloat, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, time.Time{}, err
	}
	sec := int64(tsFloat)
	nsec := int64((tsFloat - float64(sec)) * 1e9)
	ts := time.Unix(sec, nsec)
	return counter, ts, nil
}

// primaryMode er hovedløkken til primæren: Teller opp, skriver heartbeat og spawner backup én gang.
func primaryMode(startCounter int) {
	fmt.Printf("Jeg er PRIMARY: starter telling fra %d\n", startCounter)
	backupSpawned := false
	counter := startCounter

	for {
		if !backupSpawned {
			spawnBackup()
			backupSpawned = true
		}

		writeHeartbeat(counter)
		fmt.Println(counter)
		counter++
		time.Sleep(heartbeatInterval)
	}
}

// backupMode sjekker heartbeat-filen, og hvis den ikke oppdateres innen timeout, tar den over som primær.
func backupMode() {
	fmt.Println("Jeg er BACKUP,,,,sjekker om primæren er i live...")
	for {
		// Hvis heartbeat-filen ikke finnes, vent litt og prøv igjen.
		if _, err := os.Stat(heartbeatFile); os.IsNotExist(err) {
			time.Sleep(backupCheckDelay)
			continue
		}

		counter, ts, err := readHeartbeat()
		if err != nil {
			fmt.Println("Feil ved lesing av heartbeat:", err)
			time.Sleep(backupCheckDelay)
			continue
		}

		// Hvis for mye tid har gått siden siste heartbeat, tar backupen over.
		if time.Since(ts) > timeoutDuration {
			fmt.Println("Ingen heartbeat på en stund! Tar over som PRIMARY 🔥")
			newStart := counter + 1
			// Spawn en ny backup før vi tar over
			//spawnBackup()
			primaryMode(newStart)
			return
		}
		time.Sleep(backupCheckDelay)
	}
}

func main() {
	// Sjekk om programmet kjører i backup-modus
	if len(os.Args) > 1 && os.Args[1] == "--backup" {
		backupMode()
	} else {
		primaryMode(1)
	}
}
