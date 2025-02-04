package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

const (
	statusFile  = "process_status.txt"
	aliveSignal = "alive"
)

var (
	mutex       sync.Mutex
	counterFile = "counter.txt"
)

func main() {
	args := os.Args

	if len(args) == 2 && args[1] == "backup" {
		runBackup()
	} else {
		runPrimary()
	}
}

func runPrimary() {
	fmt.Println("Primary process started")
	initializeCounter()
	go broadcastAliveSignal()

	for {
		mutex.Lock()
		counter := readCounter()
		counter++
		writeCounter(counter)
		mutex.Unlock()

		fmt.Printf("Primary counting: %d\n", counter)
		time.Sleep(1 * time.Second)

		if counter%10 == 0 {
			// Simulate a crash
			fmt.Println("Primary process simulating a crash...")
			os.Exit(1)
		}
	}
}

func runBackup() {
	fmt.Println("Backup process started")

	for {
		time.Sleep(2 * time.Second)

		if !isPrimaryAlive() {
			fmt.Println("Backup taking over as primary...")
			spawnBackup()
			runPrimary()
		}
	}
}

func broadcastAliveSignal() {
	for {
		time.Sleep(1 * time.Second)
		mutex.Lock()
		writeStatusFile(aliveSignal)
		mutex.Unlock()
	}
}

func isPrimaryAlive() bool {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := os.ReadFile(statusFile)
	if err != nil || string(data) != aliveSignal {
		return false
	}
	return true
}

func spawnBackup() {
	fmt.Println("Spawning new backup...")
	cmd := exec.Command(os.Args[0], "backup")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to spawn backup:", err)
		os.Exit(1)
	}
}

func initializeCounter() {
	_, err := os.Stat(counterFile)
	if os.IsNotExist(err) {
		writeCounter(0)
	}
}

func readCounter() int {
	data, err := os.ReadFile(counterFile)
	if err != nil {
		fmt.Println("Error reading counter:", err)
		return 0
	}
	counter, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("Error converting counter to int:", err)
		return 0
	}
	return counter
}

func writeCounter(counter int) {
	err := os.WriteFile(counterFile, []byte(strconv.Itoa(counter)), 0644)
	if err != nil {
		fmt.Println("Error writing counter:", err)
	}
}

func writeStatusFile(status string) {
	err := os.WriteFile(statusFile, []byte(status), 0644)
	if err != nil {
		fmt.Println("Error writing status file:", err)
	}
}
