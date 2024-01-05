package jdwp

type EventRequestManager interface {
	CreateClassPrepareRequest() ClassPrepareRequest
	CreateClassUnloadRequest() ClassUnloadRequest
	CreateThreadStartRequest() ThreadStartRequest
	CreateThreadDeathRequest() ThreadDeathRequest
	CreateMethodEntryRequest() MethodEntryRequest
	CreateMethodExitRequest() MethodExitRequest
	CreateStepRequest(thread ThreadReference, size, depth int) StepRequest
	CreateBreakpointRequest(location Location) BreakpointRequest
	CreateAccessWatchpointRequest(field Field) AccessWatchpointRequest
	DeleteAllBreakpoints()
	GetStepRequests() []StepRequest
	GetClassPrepareRequests() []ClassPrepareRequest
	GetClassUnloadRequests() []ClassUnloadRequest
	GetThreadStartRequests() []ThreadStartRequest
	GetExceptionRequests() []ExceptionRequest
	GetBreakpointRequests() []BreakpointRequest
	GetAccessWatchpointRequests() []AccessWatchpointRequest
	GetMethodEntryRequests() []MethodEntryRequest
	GetMethodExitRequests() []MethodExitRequest
	GetVmDeathRequests() []VMDeathRequest
}
