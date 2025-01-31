package main

import (
	"Driver-go/elevator/elevio"
	"fmt"
)

func main() {

	numFloors := 4

	elevio.Init("localhost:15657", numFloors)

	var d elevio.MotorDirection = elevio.MD_Up
	//elevio.SetMotorDirection(d)

	current_order := 1
	current_floor := -1

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			elevio.SetButtonLamp(a.Button, a.Floor, true)
			current_order = a.Floor
			current_floor = elevio.GetFloor()

			if current_order == current_floor {
				d = elevio.MD_Stop
			} else if current_order > current_floor {
				d = elevio.MD_Up
			} else if current_order < current_floor {
				d = elevio.MD_Down
			}
			elevio.SetMotorDirection(d)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if current_order == a {
				d = elevio.MD_Stop
			} else if current_order > a {
				d = elevio.MD_Up
			} else if current_order < a {
				d = elevio.MD_Down
			}
			elevio.SetMotorDirection(d)
			//elevio.SetFloorIndicator(elevio.GetFloor())

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(d)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			for f := 0; f < numFloors; f++ {
				for b := elevio.ButtonType(0); b < 3; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}
