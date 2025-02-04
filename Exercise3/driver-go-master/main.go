package main

import (
	. "Driver-go/elevator"
)

func main() {
	e := Elevator{}

	RunElevatorFsm(e)
}
