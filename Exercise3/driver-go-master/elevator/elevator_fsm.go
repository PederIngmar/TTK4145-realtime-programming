package elevator

import (
	"Driver-go/config"
	"Driver-go/elevator/elevio"
)

// 3 states: Idle, Moving, Door open
// events Button press, Arrive at floor, Timer timed out

// 3 functions chooseDirection, shoushouldStopld_stop, clearRequestsAtFloor
/// mama a girl

func requestsAbove(e Elevator) bool {
	for f := e.Floor + 1; f < N_FLOORS; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Queue[f][btn] {
				return true
			}
		}
	}
	return false
}

func requestsBelow(e Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			if e.Queue[e.Floor][btn] {
				return true
			}
		}
	}
	return false
}

func requestsHere(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.Queue[e.Floor][btn] {
			return true
		}
	}
	return false
}

func chooseDirection(e Elevator) (ElevatorDir, ElevatorState) {
	switch e.Dir {
	case Up:
		if requestsAbove(e) {
			return Up, Moving
		} else if requestsHere(e) {
			return Down, DoorOpen
		} else if requestsBelow(e) {
			return Down, Moving
		} else {
			return Stop, Idle
		}
	case Down:
		if requestsBelow(e) {
			return Down, Moving
		} else if requestsHere(e) {
			return Up, DoorOpen
		} else if requestsAbove(e) {
			return Up, Moving
		} else {
			return Stop, Idle
		}
	case Stop:
		if requestsHere(e) {
			return Stop, DoorOpen
		} else if requestsAbove(e) {
			return Up, Moving
		} else if requestsBelow(e) {
			return Down, Moving
		} else {
			return Stop, Idle
		}
	default:
		return Stop, Idle
	}

}

func shouldStop(e Elevator) bool {
	switch e.Dir {
	case Down:
		return e.Queue[e.Floor][elevio.BT_HallDown] || e.Queue[e.Floor][elevio.BT_Cab] || !requestsBelow(e)
	case Up:
		return e.Queue[e.Floor][elevio.BT_HallUp] || e.Queue[e.Floor][elevio.BT_Cab] || !requestsAbove(e)
	case Stop:
		return true
	default:
		return true
	}
}

func clearRequestsAtFloor(e Elevator) {
	switch config.CLEAR_REQUEST_VARIANT {
	case config.All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.Queue[e.Floor][btn] = false
		}

	case config.InDirn:
		e.Queue[e.Floor][elevio.BT_Cab] = false
		switch e.Dir {
		case Up:
			if !requestsAbove(e) && !e.Queue[e.Floor][elevio.BT_HallUp] {
				e.Queue[e.Floor][elevio.BT_HallDown] = false
			}
			e.Queue[e.Floor][elevio.BT_HallUp] = false

		case Down:
			if !requestsBelow(e) && !e.Queue[e.Floor][elevio.BT_HallDown] {
				e.Queue[e.Floor][elevio.BT_HallUp] = false
			}
			e.Queue[e.Floor][elevio.BT_HallDown] = false

		case Stop:
			e.Queue[e.Floor][elevio.BT_HallUp] = false
			e.Queue[e.Floor][elevio.BT_HallDown] = false
		default:
			break
		}
	default:
		break
	}
}


func runElevatorFsm(e Elevator) {
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	elevio.Init("localhost:15657", 4)

	for {
		select {
			case a := <-drv_buttons: // button press
				// add to queue
				// 
			case a := <-drv_floors: // arrive at floor
			// should stop?
			// if yes: stop, clear requests at floor 
			//-> open door, start timer
			// if no: continue moving
			
			

			// in open door state: if timer timed out and no requests -> idle else -> moving
				
			case a := <-drv_obstr: // obstruction
				fmt.Printf("%+v\n", a)
				if a {
					elevio.SetMotorDirection(elevio.MD_Stop)
				} else {
					elevio.SetMotorDirection(d)
				}

			case a := <-drv_stop: // 
				fmt.Printf("%+v\n", a)
				for f := 0; f < numFloors; f++ {
					for b := elevio.ButtonType(0); b < 3; b++ {
						elevio.SetButtonLamp(b, f, false)
					}
				}
			}
	}
}