package jdwp

type EventRequest interface {
	Mirror
	IsEnabled() bool
	SetEnabled(bool)
	Enable()
	Disable()
	SetSuspendPolicy(policy SuspendPolicy)
	GetSuspendPolicy() SuspendPolicy
	SetHandler(func(request EventObject) bool)
	GetKindType() EventKind
}

type AccessWatchpointRequest EventRequest
type ClassPrepareRequest interface {
	EventRequest
	AddClassFilter(referenceType ReferenceType)
	AddClassNameFilter(classPattern string)
	AddClassExclusionFilter(classPattern string)
}
type ClassUnloadRequest interface {
	EventRequest
	AddClassNameFilter(classPattern string)
	AddClassExclusionFilter(classPattern string)
}
type BreakpointRequest interface {
	EventRequest
	GetLocation() Location
	AddThreadFilter(reference ThreadReference)
	AddInstanceFilter(reference ObjectReference)
}
type ThreadStartRequest interface {
	EventRequest
	AddThreadFilter(ThreadReference)
}
type ThreadDeathRequest interface {
	EventRequest
	AddThreadFilter(ThreadReference)
}
type ExceptionRequest interface {
	EventRequest
	GetException() ReferenceType
	NotifyCaught() bool
	NotifyUncaught() bool
	AddThreadFilter(reference ThreadReference)
	AddClassFilter(referenceType ReferenceType)
	AddClassNameFilter(classPattern string)
	AddClassExclusionFilter(classPattern string)
	AddInstanceFilter(reference ObjectReference)
}
type MethodEntryRequest interface {
	EventRequest
	AddThreadFilter(reference ThreadReference)
	AddClassFilter(referenceType ReferenceType)
	AddClassNameFilter(classPattern string)
	AddClassExclusionFilter(classPattern string)
	AddInstanceFilter(reference ObjectReference)
}
type MethodExitRequest interface {
	EventRequest
	AddThreadFilter(reference ThreadReference)
	AddClassFilter(referenceType ReferenceType)
	AddClassNameFilter(classPattern string)
	AddClassExclusionFilter(classPattern string)
	AddInstanceFilter(reference ObjectReference)
}
type StepRequest interface {
	EventRequest
	GetThread() ThreadReference
	GetSize() int
	GetDepth() int
	AddClassFilter(referenceType ReferenceType)
	AddClassNameFilter(classPattern string)
	AddClassExclusionFilter(classPattern string)
	AddInstanceFilter(reference ObjectReference)
}
type VMDeathRequest interface {
	EventRequest
}
type WatchpointRequest interface {
	EventRequest
	GetField() Field
	AddThreadFilter(reference ThreadReference)
	AddClassFilter(referenceType ReferenceType)
	AddClassNameFilter(classPattern string)
	AddClassExclusionFilter(classPattern string)
	AddInstanceFilter(reference ObjectReference)
}
