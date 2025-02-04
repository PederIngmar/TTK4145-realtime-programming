package main

import (
	"Driver-go/elevator"
)

func main() {
	e := elevator.Elevator{}

	elevator.RunElevatorFsm(e)
}
