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
	heartbeatFile     = "heartbeat.txt" 
	heartbeatInterval = 1 * time.Second         
	timeoutDuration   = 2500 * time.Millisecond 
	backupCheckDelay  = 500 * time.Millisecond  
)

func spawnBackup() {
	fmt.Println("Spawner backup🔥")
	// exePath, err := os.Executable()
	// if err != nil {
	// 	fmt.Println("Feil ved å hente kjørbar bane:", err)
	// 	return
	// }

	cmd := exec.Command("gnome-terminal", "--", "go", "run", "main.go", "--backup")
	if err := cmd.Start(); err != nil {
		fmt.Println("Feil ved spawning av backup:", err)
	}
}
func writeHeartbeat(counter int) {
	
	content := fmt.Sprintf("%d %f\n", counter, float64(time.Now().UnixNano())/1e9)
	err := os.WriteFile(heartbeatFile, []byte(content), 0644)
	if err != nil {
		fmt.Println("Feil ved skriving av heartbeat:", err)
	}
}

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

func primaryMode(startCounter int) {
	fmt.Printf("Jeg er PRIMARY🔥: starter telling fra %d\n", startCounter)
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
func backupMode() {
	fmt.Println("Jeg er BACKUP🔥 sjekker om primæren er i live...")
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
