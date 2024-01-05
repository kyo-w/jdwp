package jdwp

const (
	// SuspendNone suspends no threads when a event is raised.
	SuspendNone = SuspendPolicy(0)
	// SuspendEventThread suspends only the event's thread when a event is raised.
	SuspendEventThread = SuspendPolicy(1)
	// SuspendAll suspends all threads when a event is raised.
	SuspendAll = SuspendPolicy(2)
)
