package elevator

import (
	"Driver-go/config"
	"Driver-go/elevator/elevio"
	"fmt"
	"time"
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
			if e.Queue[f][btn] {
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
			fmt.Println("Jeg er i denne etasjen")
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
			fmt.Println("Jeg er i denne etasjen")
			return Up, DoorOpen
		} else if requestsAbove(e) {
			return Up, Moving
		} else {
			return Stop, Idle
		}
	case Stop:
		if requestsHere(e) {
			fmt.Println("Jeg er i denne etasjen")
			return Stop, DoorOpen
		} else if requestsAbove(e) {
			fmt.Println("Jeg beveger meg opp")
			return Up, Moving
		} else if requestsBelow(e) {
			fmt.Println("Jeg beveger meg ned")
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
	//case Stop:
	//return true
	default:
		return true
	}
}

func clearRequestsAtFloor(e Elevator) Elevator {
	switch config.CLEAR_REQUEST_VARIANT {
	case config.All:
		fmt.Printf("Clearing all requests at floor: %v\n", e.Floor)
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.Queue[e.Floor][btn] = false
		}
		break

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
	return e
}

func setAllLights(e Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, e.Queue[floor][btn])
		}
	}

}

func shouldClearImmediately(e Elevator, floor int, button elevio.ButtonType) bool {
	switch config.CLEAR_REQUEST_VARIANT {
	case config.All:
		return e.Floor == floor
	case config.InDirn:
		return e.Floor == floor && ((e.Dir == Up && button == elevio.BT_HallUp) ||
			(e.Dir == Down && button == elevio.BT_HallDown) ||
			e.Dir == Stop ||
			button == elevio.BT_Cab)
	default:
		return false
	}
}

func ElevatorInit(elev Elevator, drv_floors <-chan int) Elevator {
	if (elev == Elevator{}) {
		elev = Elevator{
			Floor: elevio.GetFloor(),
			State: Idle,
			Dir:   Stop,
			Queue: [N_FLOORS][N_BUTTONS]bool{},
		}
	}
	for floor := range elev.Queue {
		elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
		elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
	}

	return elev
}

func RunElevatorFsm(elev Elevator) {
	doorTimer := time.NewTimer(config.DOOR_OPEN_TIME)
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	elevio.Init("localhost:15657", config.NUM_FLOORS)
	e := ElevatorInit(elev, drv_floors)
	elevio.SetStopLamp(false)
	go setAllLights(elev)

	for {
		select {
		case a := <-drv_buttons: // button press
			// add to queue
			fmt.Printf("\nbutton request at: %+v", a)
			//print e.State
			//fmt.Printf("Current state: %v\n", e.State)

			switch e.State {
			case Idle:
				fmt.Println("IDLE 1")
				e.Queue[a.Floor][a.Button] = true // fiks kø
				e.Dir, e.State = chooseDirection(e)

				switch e.State {
				case DoorOpen:
					fmt.Println("DoorOpen 3")
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(config.DOOR_OPEN_TIME)
					clearRequestsAtFloor(e)
					e.State = DoorOpen
					fmt.Printf("Door timer: %+v", doorTimer)

				case Moving:
					elevio.SetMotorDirection(elevio.MotorDirection(e.Dir))
					fmt.Printf("Moving 3: set motor dir to: %+v", elevio.MotorDirection(e.Dir))
				case Idle:
					fmt.Println("Idle again 3")
					break
				}

			case DoorOpen:
				fmt.Println("DoorOpen 1")
				if shouldClearImmediately(e, a.Floor, a.Button) {
					time.NewTimer(config.DOOR_OPEN_TIME)
				} else {
					e.Queue[a.Floor][a.Button] = true
				}
			case Moving:
				fmt.Println("Moving 1")
				e.Queue[a.Floor][a.Button] = true
				//print queue
				fmt.Printf("Queue: %+v\n", e.Queue)
			default:
				fmt.Printf("Default")

			}

			setAllLights(e)
			fmt.Printf("\nNew state: %+v\n", e.State)

		case a := <-drv_floors: // arrive at floor
			// should stop?
			// if yes: stop, clear requests at floor
			//-> open door, start timer
			// if no: continue moving
			fmt.Printf("Case 1: Arrived at floor: %v", a)
			e.Floor = a

			switch e.State {
			case Moving:
				fmt.Printf("Case 2: Moving")
				if shouldStop(e) {
					fmt.Printf("VI MÅ STOPPE")
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					e := clearRequestsAtFloor(e)
					//time.NewTimer(config.DOOR_OPEN_TIME)
					doorTimer.Reset(config.DOOR_OPEN_TIME)
					setAllLights(e)
					e.State = DoorOpen //////// Stays open
					fmt.Printf("Case 3: Should stop")
					fmt.Printf("Current direction: %v\n", e.Dir)
					//e.State = Idle
				}
			default:
				fmt.Printf("Default: No button press")
				//break
			}

		// in open door state: if timer timed out and no requests -> idle else -> moving

		case a := <-drv_obstr: // obstruction
			fmt.Printf("Case 1: Obstruction")
			fmt.Printf("%+v\n", a)
			if a {
				elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				break
			}

		case a := <-drv_stop:
			fmt.Printf("Stopping %v\n", a)
			for f := 0; f < config.NUM_FLOORS; f++ {
				for b := elevio.ButtonType(0); b < config.NUM_BUTTONS; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
		case <-doorTimer.C:
			fmt.Println("Door timer timed out")
			elevio.SetDoorOpenLamp(false)
			prevDir := e.Dir
			e.Dir, e.State = chooseDirection(e)
			if prevDir != e.Dir && (e.Queue[e.Floor][elevio.BT_HallUp] || e.Queue[e.Floor][elevio.BT_HallDown]) {
				elevio.SetDoorOpenLamp(true)
				doorTimer.Reset(config.DOOR_OPEN_TIME)
				clearRequestsAtFloor(e)
				continue
			}
			if e.Dir == Stop {
				e.State = Idle
				//elevio.SetMotorDirection(elevio.MD_Stop)
			} else {
				elevio.SetMotorDirection(elevio.MotorDirection(e.Dir))
				e.State = Moving
				elevio.SetDoorOpenLamp(false)
			}
		}
	}
}
