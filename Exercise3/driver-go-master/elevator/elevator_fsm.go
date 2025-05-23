package elevator

import (
	"Driver-go/config"
	"Driver-go/elevator/elevio"
	"fmt"
	"time"
)

// 3 states: Idle, Moving, Door open
// events Button press, Arrive at floor, Timer timed out

// 3 functions chooseDirection, should_stop, clearRequestsAtFloor

// Returns true if there are requests above the elevator's current floor
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

// Returns true if there are requests below the elevator's current floor
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

// Returns true if there are requests at the elevator's current floor
func requestsHere(e Elevator) bool {
	for btn := 0; btn < N_BUTTONS; btn++ {
		if e.Queue[e.Floor][btn] {
			return true
		}
	}
	return false
}

// Chooses the direction the elevator should move in, and the state it should be in
// func chooseDirection(e Elevator) (ElevatorDir, ElevatorState) {
// 	switch {
// 	case requestsAbove(e):
// 		return Up, Moving
// 	case requestsHere(e):
// 		return Stop, DoorOpen
// 	case requestsBelow(e):
// 		return Down, Moving
// 	default:
// 		return Stop, Idle
// 	}
// }

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

// Returns true if the elevator should stop at the current floor
//
//	func shouldStop(e Elevator) bool {
//		return e.Queue[e.Floor][elevio.BT_HallUp] ||
//			e.Queue[e.Floor][elevio.BT_HallDown] ||
//			e.Queue[e.Floor][elevio.BT_Cab]
//	}


// Må se mer på en bedre løsning for denne
func shouldStop(e Elevator) bool {
	switch config.CLEAR_REQUEST_VARIANT {
	case config.All:
		return e.Queue[e.Floor][elevio.BT_HallUp] ||
			e.Queue[e.Floor][elevio.BT_HallDown] ||
			e.Queue[e.Floor][elevio.BT_Cab]
	case config.InDir:
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
	default:
	}
	return false

}

// Clears all requests at the elevator's current floor
func clearRequestsAtFloor(e *Elevator) {
	switch config.CLEAR_REQUEST_VARIANT {
	case config.All:
		for btn := 0; btn < N_BUTTONS; btn++ {
			e.Queue[e.Floor][btn] = false
		}

	case config.InDir:
		e.Queue[e.Floor][elevio.BT_Cab] = false
		switch e.Dir {
		case Up:
			if !requestsAbove(*e) && !e.Queue[e.Floor][elevio.BT_HallUp] {
				e.Queue[e.Floor][elevio.BT_HallDown] = false
			}
			e.Queue[e.Floor][elevio.BT_HallUp] = false
		case Down:
			if !requestsBelow(*e) && !e.Queue[e.Floor][elevio.BT_HallDown] {
				e.Queue[e.Floor][elevio.BT_HallUp] = false
			}
			e.Queue[e.Floor][elevio.BT_HallDown] = false
		case Stop:
			e.Queue[e.Floor][elevio.BT_HallUp] = false
			e.Queue[e.Floor][elevio.BT_HallDown] = false
		default:
		}
	default:
	}
}

// Sets all lights in the elevator to match the elevator's queue
func setAllLights(e Elevator) {
	for floor := 0; floor < N_FLOORS; floor++ {
		for btn := 0; btn < N_BUTTONS; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, e.Queue[floor][btn])
		}
	}

}

// Returns true if the elevator should clear requests immediately when arriving at a floor
func shouldClearImmediately(e Elevator, floor int, button elevio.ButtonType) bool {
	fmt.Println("Should clear immediately?")
	switch config.CLEAR_REQUEST_VARIANT {
	case config.All:
		return e.Floor == floor
	case config.InDir:
		fmt.Println("CLEARING INDIR")
		return e.Floor == floor && ((e.Dir == Up && button == elevio.BT_HallUp) ||
			(e.Dir == Down && button == elevio.BT_HallDown) ||
			e.Dir == Stop ||
			button == elevio.BT_Cab)
	default:
		return false
	}
}

// Initializes the elevator
func ElevatorInit() Elevator {
	return Elevator{
		Floor: elevio.GetFloor(),
		//Floor: 0,
		State: Idle,
		Dir:   Stop,
		Queue: [N_FLOORS][N_BUTTONS]bool{},
	}
}

func RunElevatorFSM() {
	doorTimer := time.NewTimer(config.DOOR_OPEN_TIME) // New timer created with duration DOOR_OPEN_TIME
	if !doorTimer.Stop() {
		<-doorTimer.C
	}

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go elevio.PollButtons(drv_buttons)
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollStopButton(drv_stop)

	elevio.Init("localhost:15657", config.NUM_FLOORS) // Initialize elevator hardware
	e := ElevatorInit()                               // Initialize elevator struct
	setAllLights(e)
	if elevio.GetFloor() == -1 { // If the elevator is between floors
		elevio.SetMotorDirection(elevio.MD_Down) // Return the elevator to the nearest floor
		e.Dir = Down
		e.State = Moving
	initLoop:
		for {
			select {
			case <-drv_floors:
				elevio.SetMotorDirection(elevio.MD_Stop)
				e.Floor = elevio.GetFloor()
				elevio.SetFloorIndicator(e.Floor)
				fmt.Println("Ankommet etasje: ", e.Floor)
				e.State = Idle
				e.Dir, e.State = chooseDirection(e)
				elevio.SetMotorDirection(elevio.MotorDirection(e.Dir))
				break initLoop

			case <-time.After(5 * time.Second):
				fmt.Println("Venter på etasjesensor")
			}
		}
	}

	for {
		select {
		// Knappetrykk
		case btnPress := <-drv_buttons:
			fmt.Printf("\nKnappetrykk: %+v\n", btnPress)
			//e.Queue[btnPress.Floor][btnPress.Button] = true
			//setAllLights(e)
			switch e.State {

			case Idle:
				// Legger forespørselen til i køen og velger ny retning
				e.Queue[btnPress.Floor][btnPress.Button] = true
				e.Dir, e.State = chooseDirection(e) // Velger ny retning

				switch e.State {

				case DoorOpen:
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(config.DOOR_OPEN_TIME)
					clearRequestsAtFloor(&e)

				case Moving:
					elevio.SetMotorDirection(elevio.MotorDirection(e.Dir))

				case Idle:
				}

			case DoorOpen:
				fmt.Println("Button press while door is open")
				if shouldClearImmediately(e, btnPress.Floor, btnPress.Button) {
					doorTimer.Reset(config.DOOR_OPEN_TIME)
				} else {
					e.Queue[btnPress.Floor][btnPress.Button] = true
				}

			case Moving:
				e.Queue[btnPress.Floor][btnPress.Button] = true
			}

			setAllLights(e)

		case floor := <-drv_floors:
			fmt.Printf("Ankommet etasje: %d\n", floor)
			e.Floor = floor
			elevio.SetFloorIndicator(e.Floor)

			if e.State == Moving && shouldStop(e) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elevio.SetDoorOpenLamp(true)
				clearRequestsAtFloor(&e)
				doorTimer.Reset(config.DOOR_OPEN_TIME)
				setAllLights(e)
				e.State = DoorOpen
			}

		case obstr := <-drv_obstr:
			fmt.Printf("Hindring: %v", obstr)
			if e.State == DoorOpen {
				if elevio.GetObstruction() {
					elevio.SetMotorDirection(elevio.MD_Stop)
					elevio.SetDoorOpenLamp(true)
					fmt.Println("Obstruksjon aktiv. Døren forblir åpen")
					<-drv_obstr
					if !elevio.GetObstruction() {
						fmt.Println("Obstruksjon fjernet")
						doorTimer.Reset(config.DOOR_OPEN_TIME)
					}
				} else {
					fmt.Println("Obstruksjon registrert, men døren er lukket")
				}
			}
		// Stopp-knapp
		case <-drv_stop:
			fmt.Println("Stopp-knapp trykket")
			for f := 0; f < config.NUM_FLOORS; f++ {
				for b := elevio.ButtonType(0); b < config.NUM_BUTTONS; b++ {
					elevio.SetButtonLamp(b, f, false)
				}
			}
			//elevio.SetMotorDirection(elevio.MD_Stop)
			//elevio.SetDoorOpenLamp(false)

		case <-doorTimer.C:
			fmt.Println("Dør-timer utløpt")
			if e.State == DoorOpen {
				newDir, newState := chooseDirection(e)
				e.Dir, e.State = newDir, newState
				switch e.State {
				case DoorOpen:
					elevio.SetDoorOpenLamp(true)
					clearRequestsAtFloor(&e)
					setAllLights(e)
					doorTimer.Reset(config.DOOR_OPEN_TIME)
				case Moving:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(elevio.MotorDirection(e.Dir))
				case Idle:
					elevio.SetDoorOpenLamp(false)
					elevio.SetMotorDirection(elevio.MD_Stop)
				}
			} else {
				fmt.Printf("Door timer, heis state: %v", e.State)
			}
		}
	}
}

func requests_clearAtCurrentFloor(e_old Elevator, onCleared func(elevio.ButtonType, int)) Elevator {
	e := e_old
	for btn := elevio.ButtonType(0); btn < N_BUTTONS; btn++ {
		if e.Queue[e.Floor][btn] {
			e.Queue[e.Floor][btn] = false
			if onCleared != nil {
				onCleared(btn, e.Floor)
			}
		}
	}
	return e
}

// Returns the time it takes for the elevator to reach the floor of the button press
func TimeToServeRequest(e_copy Elevator, btnPress elevio.ButtonEvent) time.Duration {
	duration := 0 * time.Second

	elevatorArrival := 0

	e := e_copy
	e.Queue[btnPress.Floor][btnPress.Button] = true // Add the button press to the queue

	// Function to be called when a request is cleared
	// Sets elevatorArrival to 1 if the button press is cleared
	onCleared := func(btn elevio.ButtonType, floor int) {
		if btn == btnPress.Button && floor == btnPress.Floor {
			elevatorArrival = 1
		}
	}
	switch e.State {
	case Idle:
		e.Dir, e.State = chooseDirection(e)
		if e.Dir == Stop {
			return duration // Elevator is already at the floor
		}
	case Moving:
		duration += config.TRAVEL_TIME / 2
		e.Floor += int(e.Dir)
	case DoorOpen:
		duration -= config.DOOR_OPEN_TIME / 2
		if !requestsAbove(e) && !requestsBelow(e) {
			return duration
		}
	}
	for {
		if shouldStop(e) {
			e = requests_clearAtCurrentFloor(e, onCleared)
			if elevatorArrival == 1 {
				return duration
			}
			duration += config.DOOR_OPEN_TIME
			e.Dir, _ = chooseDirection(e)
		}
		e.Floor += int(e.Dir)
		duration += config.TRAVEL_TIME
	}
}
