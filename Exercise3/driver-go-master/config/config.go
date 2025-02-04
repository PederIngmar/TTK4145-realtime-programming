package config

import("time")

type ClearRequestVariant int 
const (
	All ClearRequestVariant = iota
	InDirn
)

const (
	NUM_FLOORS = 4
	NUM_BUTTONS = 3
	CLEAR_REQUEST_VARIANT = All
	DOOR_OPEN_TIME = 3 * time.Second
)