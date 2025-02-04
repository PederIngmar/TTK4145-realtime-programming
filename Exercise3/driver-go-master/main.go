package main

import (
	. "Driver-go/elevator"
)

func main() {
	e := ElevatorInit()

	RunElevatorFsm(e)
}
