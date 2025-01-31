package elevator
import (
	. "Driver-go/elevator/elevio"
)

// 3 states: Idle, Moving, Door open
// events Button press, Arrive at floor, Timer timed out

// 3 functions chooseDirection, shoushouldStopld_stop, clearRequestsAtFloor


func requests_above(e Elevator) bool{
    for f := e.Floor+1; f < N_FLOORS; f++ {
        for btn := 0; btn < N_BUTTONS; btn++{
            if(e.Queue[f][btn]){
                return true
            }
        }
    }
    return false
}


func requests_below(e Elevator) bool{
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < N_BUTTONS; btn++{
			if(e.Queue[e.Floor][btn]){
				return true
			}
		}
	}
    return false
}

func requests_here(e Elevator) bool{
	for btn := 0; btn < N_BUTTONS; btn++{
		if(e.Queue[e.Floor][btn]){
			return true
		}
	}
    return false
}



func chooseDirection(e Elevator) (ElevatorDir, ElevatorState){
	switch e.Dir {
	case Up:
		if requests_above(e) {
			return Up, Moving
		} else if requests_here(e) {
			return Down, DoorOpen
		} else if requests_below(e) {
			return Down, Moving
		} else {
			return Stop, Idle
		}
	case Down:
		if requests_below(e) {
			return Down, Moving
		} else if requests_here(e) {
			return Up, DoorOpen
		} else if requests_above(e) {
			return Up, Moving
		} else {
			return Stop, Idle
		}
	case Stop:
		if requests_here(e) {
			return Stop, DoorOpen
		} else if requests_above(e) {
			return Up, Moving
		} else if requests_below(e) {
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
		return e.Queue[e.Floor][BT_HallDown] || e.Queue[e.Floor][BT_Cab] || !requests_below(e)
	case Up:
		return e.Queue[e.Floor][BT_HallUp] || e.Queue[e.Floor][BT_Cab] || !requests_above(e)
	case Stop:
		return true
	default:
		return true
	}
}

func clearRequestsAtFloor(e Elevator) {

}
