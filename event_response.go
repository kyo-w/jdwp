// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jdwp

type EventKind uint8

const (
	// SingleStep is the kind of event raised when a single-step has been completed.
	SingleStep = EventKind(1)
	// Breakpoint is the kind of event raised when a breakpoint has been hit.
	Breakpoint = EventKind(2)
	// FramePop is the kind of event raised when a stack-frame is popped.
	FramePop = EventKind(3)
	// Exception is the kind of event raised when an exception is thrown.
	Exception = EventKind(4)
	// UserDefined is the kind of event raised when a user-defind event is fired.
	UserDefined = EventKind(5)
	// ThreadStart is the kind of event raised when a new thread is started.
	ThreadStart = EventKind(6)
	// ThreadDeath is the kind of event raised when a thread is stopped.
	ThreadDeath = EventKind(7)
	// ClassPrepare is the kind of event raised when a class enters the prepared state.
	ClassPrepare = EventKind(8)
	// ClassUnload is the kind of event raised when a class is unloaded.
	ClassUnload = EventKind(9)
	// ClassLoad is the kind of event raised when a class enters the loaded state.
	ClassLoad = EventKind(10)
	// FieldAccess is the kind of event raised when a field is accessed.
	FieldAccess = EventKind(20)
	// FieldModification is the kind of event raised when a field is modified.
	FieldModification = EventKind(21)
	// ExceptionCatch is the kind of event raised when an exception is caught.
	ExceptionCatch = EventKind(30)
	// MethodEntry is the kind of event raised when a method has been entered.
	MethodEntry = EventKind(40)
	// MethodExit is the kind of event raised when a method has been exited.
	MethodExit = EventKind(41)
	// VMStart is the kind of event raised when the virtual machine is initialized.
	VMStart = EventKind(90)
	// VMDeath is the kind of event raised when the virtual machine is shutdown.
	VMDeath = EventKind(99)
)

// EventsResponse Events is a collection of events.
type EventsResponse struct {
	Policy SuspendPolicy
	Events []EventResponse
}

// EventResponse Event is the interface implemented by all events raised by the VM.
type EventResponse interface {
	GetRequest() EventRequestID
	Kind() EventKind
}

// EventVMStartResponse EventVMStart represents an event raised when the virtual machine is started.
type EventVMStartResponse struct {
	Request EventRequestID
	Thread  ThreadID
}

// EventVMDeathResponse EventVMDeath represents an event raised when the virtual machine is stopped.
type EventVMDeathResponse struct {
	Request EventRequestID
}

// EventSingleStepResponse EventSingleStep represents an event raised when a single-step has been completed.
type EventSingleStepResponse struct {
	Request  EventRequestID
	Thread   ThreadID
	Location LocationID
}

// EventBreakpointResponse EventBreakpoint represents an event raised when a breakpoint has been hit.
type EventBreakpointResponse struct {
	Request  EventRequestID
	Thread   ThreadID
	Location LocationID
}

// EventMethodEntryResponse EventMethodEntry represents an event raised when a method has been entered.
type EventMethodEntryResponse struct {
	Request  EventRequestID
	Thread   ThreadID
	Location LocationID
}

// EventMethodExitResponse EventMethodExit represents an event raised when a method has been exited.
type EventMethodExitResponse struct {
	Request  EventRequestID
	Thread   ThreadID
	Location LocationID
}

// EventExceptionResponse EventException represents an event raised when an exception is thrown.
type EventExceptionResponse struct {
	Request       EventRequestID
	Thread        ThreadID
	Location      LocationID
	Exception     TaggedObjectID
	CatchLocation LocationID
}

// EventThreadStartResponse EventThreadStart represents an event raised when a new thread is started.
type EventThreadStartResponse struct {
	Request EventRequestID
	Thread  ThreadID
}

// EventThreadDeathResponse EventThreadDeath represents an event raised when a thread is stopped.
type EventThreadDeathResponse struct {
	Request EventRequestID
	Thread  ThreadID
}

// EventClassPrepareResponse EventClassPrepare represents an event raised when a class enters the prepared state.
type EventClassPrepareResponse struct {
	Request   EventRequestID
	Thread    ThreadID
	ClassKind TypeTag
	ClassType ReferenceTypeID
	Signature string
	Status    ByteID
}

// EventClassUnloadResponse EventClassUnload represents an event raised when a class is unloaded.
type EventClassUnloadResponse struct {
	Request   EventRequestID
	Signature string
}

// EventFieldAccessResponse EventFieldAccess represents an event raised when a field is accessed.
type EventFieldAccessResponse struct {
	Request   EventRequestID
	Thread    ThreadID
	Location  LocationID
	FieldKind TypeTag
	FieldType ReferenceTypeID
	Field     FieldID
	Object    TaggedObjectID
}

// EventFieldModificationResponse EventFieldModification represents an event raised when a field is modified.
type EventFieldModificationResponse struct {
	Request   EventRequestID
	Thread    ThreadID
	Location  LocationID
	FieldKind TypeTag
	FieldType ReferenceTypeID
	Field     FieldID
	Object    TaggedObjectID
	NewValue  Value
}

func (e EventVMStartResponse) GetRequest() EventRequestID           { return e.Request }
func (e EventVMDeathResponse) GetRequest() EventRequestID           { return e.Request }
func (e EventSingleStepResponse) GetRequest() EventRequestID        { return e.Request }
func (e EventBreakpointResponse) GetRequest() EventRequestID        { return e.Request }
func (e EventMethodEntryResponse) GetRequest() EventRequestID       { return e.Request }
func (e EventMethodExitResponse) GetRequest() EventRequestID        { return e.Request }
func (e EventExceptionResponse) GetRequest() EventRequestID         { return e.Request }
func (e EventThreadStartResponse) GetRequest() EventRequestID       { return e.Request }
func (e EventThreadDeathResponse) GetRequest() EventRequestID       { return e.Request }
func (e EventClassPrepareResponse) GetRequest() EventRequestID      { return e.Request }
func (e EventClassUnloadResponse) GetRequest() EventRequestID       { return e.Request }
func (e EventFieldAccessResponse) GetRequest() EventRequestID       { return e.Request }
func (e EventFieldModificationResponse) GetRequest() EventRequestID { return e.Request }

// Kind returns VMStart
func (EventVMStartResponse) Kind() EventKind { return VMStart }

// Kind returns VMDeath
func (EventVMDeathResponse) Kind() EventKind { return VMDeath }

// Kind returns SingleStep
func (EventSingleStepResponse) Kind() EventKind { return SingleStep }

// Kind returns Breakpoint
func (EventBreakpointResponse) Kind() EventKind { return Breakpoint }

// Kind returns MethodEntry
func (EventMethodEntryResponse) Kind() EventKind { return MethodEntry }

// Kind returns MethodExit
func (EventMethodExitResponse) Kind() EventKind { return MethodExit }

// Kind returns Exception
func (EventExceptionResponse) Kind() EventKind { return Exception }

// Kind returns ThreadStart
func (EventThreadStartResponse) Kind() EventKind { return ThreadStart }

// Kind returns ThreadDeath
func (EventThreadDeathResponse) Kind() EventKind { return ThreadDeath }

// Kind returns ClassPrepare
func (EventClassPrepareResponse) Kind() EventKind { return ClassPrepare }

// Kind returns ClassUnload
func (EventClassUnloadResponse) Kind() EventKind { return ClassUnload }

// Kind returns FieldAccess
func (EventFieldAccessResponse) Kind() EventKind { return FieldAccess }

// Kind returns FieldModification
func (EventFieldModificationResponse) Kind() EventKind { return FieldModification }

// Event returns a default-initialized Event of the specified kind.
func (k EventKind) Event() EventResponse {
	switch k {
	case SingleStep:
		return &EventSingleStepResponse{}
	case Breakpoint:
		return &EventBreakpointResponse{}
	case Exception:
		return &EventExceptionResponse{}
	case ThreadStart:
		return &EventThreadStartResponse{}
	case ThreadDeath:
		return &EventThreadDeathResponse{}
	case ClassPrepare:
		return &EventClassPrepareResponse{}
	case ClassUnload:
		return &EventClassUnloadResponse{}
	case FieldAccess:
		return &EventFieldAccessResponse{}
	case FieldModification:
		return &EventFieldModificationResponse{}
	case ExceptionCatch:
		return &EventExceptionResponse{}
	case MethodEntry:
		return &EventMethodEntryResponse{}
	case MethodExit:
		return &EventMethodExitResponse{}
	case VMStart:
		return &EventVMStartResponse{}
	case VMDeath:
		return &EventVMDeathResponse{}
	default:
		return nil
	}
}
