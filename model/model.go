package model

type CommandStatus string

const (
	CommandStatusUnknown  CommandStatus = "Unknown"
	CommandStatusStarting               = "Starting"
	CommandStatusStarted                = "Started"
	CommandStatusRunning                = "Stopping"
	CommandStatusStopping               = "Stopping"
	CommandStatusStopped                = "Stopped"
)
