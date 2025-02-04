package main

import (
	. "Driver-go/elevator"
)

func main() {
	e := ElevatorInit()  
	runElevatorFsm(e)
}
