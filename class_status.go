package jdwp

const (
	// StatusVerified is used to describe a class in the verified state.
	StatusVerified = ClassStatus(1)
	// StatusPrepared is used to describe a class in the prepared state.
	StatusPrepared = ClassStatus(2)
	// StatusInitialized is used to describe a class in the initialized state.
	StatusInitialized = ClassStatus(4)
	// StatusError is used to describe a class in the error state.
	StatusError = ClassStatus(8)
)
