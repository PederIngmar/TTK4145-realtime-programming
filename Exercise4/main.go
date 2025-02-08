package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

const (
	aliveFile = "alive.txt"
	aliveInterval = 1 * time.Second
	backupTakeover = 2500 * time.Millisecond
	backupCheckup = 500 * time.Millisecond
)


func startBackup() {
	fmt.Println("Backup-prosess starter")

	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Feil ved å hente exePath\n", err)
		return
	}

	cmd := exec.Command("gnome-terminal", "--", exePath, "--backup")
	if err := cmd.Start(); err != nil {
		fmt.Println("Klarte ikke spawne backup brir")
	}
}

func runBackup() {
	fmt.Println("Bro er backup frrrr")
}

func runPrimary() { 

}

func main() {
	startBackup()
}