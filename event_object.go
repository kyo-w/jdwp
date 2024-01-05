package jdwp

type EventObject interface {
	GetRequest() EventResponse
}
type VMDeathEventObject EventObject

type LocatableEventObject interface {
	EventObject
	GetThread() ThreadReference
	GetLocation() Location
}
type BreakpointEventObject interface {
	LocatableEventObject
}
type StepEventObject LocatableEventObject
type MethodEntryEventObject interface {
	LocatableEventObject
	GetMethod() Method
}
type MethodExitEventObject interface {
	LocatableEventObject
	GetMethod() Method
	GetReturnValue() Value
}

type WatchpointEventObject interface {
	LocatableEventObject
	GetField() Field
	GetObject() ObjectReference
	GetValueCurrent() Value
}
type AccessWatchpointEventObject WatchpointEventObject

type ExceptionEventObject interface {
	LocatableEventObject
}

type ClassPrepareEventObject interface {
	EventObject
	GetThread() ThreadReference
	GetReferenceType() ReferenceType
}

type ClassUnloadEventObject interface {
	EventObject
	GetClasName() string
	GetClassSignature() string
}

type ThreadDeathEventObject interface {
	EventObject
	GetThread() ThreadReference
}
type ThreadStartEventObject interface {
	EventObject
	GetThread() ThreadReference
}
type VMStartEventObject interface {
	EventObject
	GetThread() ThreadReference
}
