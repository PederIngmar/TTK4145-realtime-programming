package config

type ClearRequestVariant int 
const (
	All ClearRequestVariant = iota
	InDirn
)

const (
	CLEAR_REQUEST_VARIANT = All
)