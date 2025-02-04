package main

import (
	"Driver-go/elevator/elevio"
	"Driver-go/elevator"
	"fmt"
)

func main() {
	numFloors := 4
	elevio.Init("localhost:15657", numFloors)
	e Elevator := ElevatorInit()  
	elevator.runElevatorFsm(e)
}
