package elevator

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)
// Elevator states: Idle, Moving, Door open
type ElevatorState int

const (
	Idle ElevatorState = iota
	Moving
	DoorOpen
	
	// unavailable = 3
)
// Elevator directions: Down, Stop, Up
type ElevatorDir int

const (
	Down ElevatorDir = -1
	Stop ElevatorDir = 0
	Up   ElevatorDir = 1
)

type Elevator struct {
	Floor 	int
	State 	ElevatorState
	Dir   	ElevatorDir
	Queue 	[N_FLOORS][N_BUTTONS]bool // 
}
