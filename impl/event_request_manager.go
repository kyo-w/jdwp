package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type EventRequestManagerImpl struct {
	vm                      *VirtualMachineImpl
	ClassPrepareRequest     []jdi.ClassPrepareRequest
	ClassUnloadRequest      []jdi.ClassUnloadRequest
	ThreadStartRequest      []jdi.ThreadStartRequest
	ThreadDeathRequest      []jdi.ThreadDeathRequest
	VMDeathRequest          []jdi.VMDeathRequest
	MethodExitRequest       []jdi.MethodExitRequest
	MethodEntryRequest      []jdi.MethodEntryRequest
	AccessWatchpointRequest []jdi.AccessWatchpointRequest
	BreakpointRequest       []jdi.BreakpointRequest
	ExceptionRequest        []jdi.ExceptionRequest
	StepRequest             []jdi.StepRequest
}

func (e *EventRequestManagerImpl) createRequestHook(kind jdi.EventKind) EventRequestImpl {
	return EventRequestImpl{vm: e.vm, EventKind: kind}
}

func (e *EventRequestManagerImpl) createClassRequestHook(kind jdi.EventKind) ClassVisibleEventRequestImpl {
	return ClassVisibleEventRequestImpl{ThreadVisibleEventRequestImpl: ThreadVisibleEventRequestImpl{e.createRequestHook(kind)}}
}

func (e *EventRequestManagerImpl) CreateClassPrepareRequest() jdi.ClassPrepareRequest {
	request := &ClassPrepareRequestImpl{ClassVisibleEventRequestImpl: e.createClassRequestHook(jdi.ClassPrepare)}
	e.ClassPrepareRequest = append(e.ClassPrepareRequest, request)
	return request
}

func (e *EventRequestManagerImpl) CreateClassUnloadRequest() jdi.ClassUnloadRequest {
	request := &ClassUnloadRequestImpl{ClassVisibleEventRequestImpl: e.createClassRequestHook(jdi.ClassUnload)}
	e.ClassUnloadRequest = append(e.ClassUnloadRequest, request)
	return request
}

func (e *EventRequestManagerImpl) CreateThreadStartRequest() jdi.ThreadStartRequest {
	request := &ThreadStartRequestImpl{ClassVisibleEventRequestImpl: e.createClassRequestHook(jdi.ThreadStart)}
	e.ThreadStartRequest = append(e.ThreadStartRequest, request)
	return request
}

func (e *EventRequestManagerImpl) CreateThreadDeathRequest() jdi.ThreadDeathRequest {
	request := &ThreadDeathRequestImpl{ClassVisibleEventRequestImpl: e.createClassRequestHook(jdi.ThreadDeath)}
	e.ThreadDeathRequest = append(e.ThreadDeathRequest, request)
	return request
}
func (e *EventRequestManagerImpl) CreateMethodEntryRequest() jdi.MethodEntryRequest {
	request := &MethodEntryRequestImpl{ClassVisibleEventRequestImpl: e.createClassRequestHook(jdi.MethodEntry)}
	e.MethodEntryRequest = append(e.MethodEntryRequest, request)
	return request
}

func (e *EventRequestManagerImpl) CreateMethodExitRequest() jdi.MethodExitRequest {
	request := &MethodExitRequestImpl{ClassVisibleEventRequestImpl: e.createClassRequestHook(jdi.MethodExit)}
	e.MethodExitRequest = append(e.MethodExitRequest, request)
	return request
}

func (e *EventRequestManagerImpl) CreateStepRequest(thread jdi.ThreadReference, size, depth int) jdi.StepRequest {
	request := &StepRequestImpl{ClassVisibleEventRequestImpl: e.createClassRequestHook(jdi.SingleStep)}
	request.filters = make([]jdi.EventModifier, 1)
	filter := jdi.StepEventModifier{
		Thread: jdi.ThreadID(thread.GetUniqueID()),
		Size:   size,
		Depth:  depth,
	}
	request.filters[0] = filter
	e.StepRequest = append(e.StepRequest, request)
	return request
}

func (e *EventRequestManagerImpl) CreateBreakpointRequest(location jdi.Location) jdi.BreakpointRequest {
	request := &BreakpointRequestImpl{ClassVisibleEventRequestImpl: e.createClassRequestHook(jdi.Breakpoint)}
	request.filters = make([]jdi.EventModifier, 1)
	request.filters[0] = jdi.LocationOnlyEventModifier{
		Method:   jdi.MethodID(location.GetMethod().GetUniqueID()),
		Type:     location.GetDeclaringType().GetTypeTag(),
		Class:    jdi.ClassObjectID(location.GetDeclaringType().GetUniqueID()),
		Location: uint64(location.GetCodeIndex()),
	}
	request.Location = location
	e.BreakpointRequest = append(e.BreakpointRequest, request)
	return request
}

func (e *EventRequestManagerImpl) CreateAccessWatchpointRequest(field jdi.Field) jdi.AccessWatchpointRequest {
	request := &AccessWatchpointRequestImpl{ClassVisibleEventRequestImpl: e.createClassRequestHook(jdi.FieldAccess)}
	request.filters = make([]jdi.EventModifier, 1)
	request.filters[0] = jdi.FieldOnlyEventModifier{
		Field: jdi.FieldID(field.GetUniqueID()),
		Type:  field.GetDeclaringType().GetUniqueID(),
	}
	request.Field = field
	e.AccessWatchpointRequest = append(e.AccessWatchpointRequest, request)
	return request
}

func (e *EventRequestManagerImpl) DeleteAllBreakpoints() {
	for _, value := range e.BreakpointRequest {
		value.Disable()
	}
	e.BreakpointRequest = []jdi.BreakpointRequest{}
}

func (e *EventRequestManagerImpl) GetStepRequests() []jdi.StepRequest {
	return e.StepRequest
}

func (e *EventRequestManagerImpl) GetClassPrepareRequests() []jdi.ClassPrepareRequest {
	return e.ClassPrepareRequest
}

func (e *EventRequestManagerImpl) GetClassUnloadRequests() []jdi.ClassUnloadRequest {
	return e.ClassUnloadRequest
}

func (e *EventRequestManagerImpl) GetThreadStartRequests() []jdi.ThreadStartRequest {
	return e.ThreadStartRequest
}

func (e *EventRequestManagerImpl) GetExceptionRequests() []jdi.ExceptionRequest {
	return e.ExceptionRequest
}

func (e *EventRequestManagerImpl) GetBreakpointRequests() []jdi.BreakpointRequest {
	return e.BreakpointRequest
}

func (e *EventRequestManagerImpl) GetAccessWatchpointRequests() []jdi.AccessWatchpointRequest {
	return e.AccessWatchpointRequest
}

func (e *EventRequestManagerImpl) GetMethodEntryRequests() []jdi.MethodEntryRequest {
	return e.MethodEntryRequest
}

func (e *EventRequestManagerImpl) GetMethodExitRequests() []jdi.MethodExitRequest {
	return e.MethodExitRequest
}

func (e *EventRequestManagerImpl) GetVmDeathRequests() []jdi.VMDeathRequest {
	return e.VMDeathRequest
}
