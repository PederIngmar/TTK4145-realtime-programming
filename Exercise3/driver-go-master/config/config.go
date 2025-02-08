package config

import (
	"time"
)

type ClearRequestVariant int

const (
	All ClearRequestVariant = iota
	InDir
)

const (
	NUM_FLOORS            = 4
	NUM_BUTTONS           = 3
	CLEAR_REQUEST_VARIANT = InDir
	DOOR_OPEN_TIME        = 3 * time.Second
	TRAVEL_TIME           = 2
)
