package elevator

const (
	N_FLOORS  = 4
	N_BUTTONS = 3
)

type ElevatorState int

const (
	Idle ElevatorState = iota
	Moving
	DoorOpen
	// unavailable = 3
)

type ElevatorDir int

const (
	Down ElevatorDir = -1
	Stop ElevatorDir = 0
	Up   ElevatorDir = 1
)

type Elevator struct {
	Floor int
	State ElevatorState
	Dir   ElevatorDir
	Queue [N_FLOORS][N_BUTTONS]bool
}
