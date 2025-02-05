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
	switch {
	case requestsAbove(e):
		return Up, Moving
	case requestsHere(e):
		return Stop, DoorOpen
	case requestsBelow(e):
		return Down, Moving
	default:
		return Stop, Idle
	}
}

func shouldStop(e Elevator) bool {
	return 	e.Queue[e.Floor][elevio.BT_HallUp] || 
			e.Queue[e.Floor][elevio.BT_HallDown] || 
			e.Queue[e.Floor][elevio.BT_Cab]
}


func clearRequestsAtFloor(e *Elevator) {
	for btn := 0; btn < N_BUTTONS; btn++ {
		e.Queue[e.Floor][btn] = false
	}
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


func ElevatorInit() Elevator {
	return Elevator{
		Floor: elevio.GetFloor(),
		//Floor: 0,
		State: Idle,
		Dir:   Stop,
		Queue: [N_FLOORS][N_BUTTONS]bool{},
	}
}
func EFsm() {
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

	elevio.Init("localhost:15657", config.NUM_FLOORS)
	e := ElevatorInit()
	setAllLights(e)
	if elevio.GetFloor() == -1 {
		elevio.SetMotorDirection(elevio.MD_Down)
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

				case <-time.After(5*time.Second):
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
			}
		}
	}
}

