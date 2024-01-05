package impl

import (
	jdi "github.com/kyo-w/jdwp"
	connect "github.com/kyo-w/jdwp/impl/internal"
)

type EventRequestImpl struct {
	EventKind     jdi.EventKind
	vm            *VirtualMachineImpl
	Id            jdi.EventRequestID
	filters       []jdi.EventModifier
	isEnabled     bool
	deleted       bool
	suspendPolicy jdi.SuspendPolicy
	handler       func(request jdi.EventObject) bool
}
type BreakpointRequestImpl struct {
	ClassVisibleEventRequestImpl
	Location jdi.Location
}
type ClassPrepareRequestImpl struct {
	ClassVisibleEventRequestImpl
}
type ClassUnloadRequestImpl struct {
	ClassVisibleEventRequestImpl
}
type ClassVisibleEventRequestImpl struct {
	ThreadVisibleEventRequestImpl
}
type MethodEntryRequestImpl struct {
	ClassVisibleEventRequestImpl
}
type MethodExitRequestImpl struct {
	ClassVisibleEventRequestImpl
}
type StepRequestImpl struct {
	ClassVisibleEventRequestImpl
	Thread jdi.ThreadReference
	Size   int
	depth  int
}
type ThreadDeathRequestImpl struct {
	ClassVisibleEventRequestImpl
}
type ThreadStartRequestImpl struct {
	ClassVisibleEventRequestImpl
}
type ThreadVisibleEventRequestImpl struct {
	EventRequestImpl
}
type AccessWatchpointRequestImpl struct {
	ClassVisibleEventRequestImpl
	Field jdi.Field
}

func (e *EventRequestImpl) SetHandler(f func(request jdi.EventObject) bool) {
	e.handler = f
}
func (e *EventRequestImpl) GetKindType() jdi.EventKind {
	return e.EventKind
}
func (e *EventRequestImpl) Disable() {
	e.SetEnabled(false)
}
func (e *EventRequestImpl) GetVirtualMachine() jdi.VirtualMachine {
	return e.vm
}
func (e *EventRequestImpl) delete() {
	if !e.deleted {
		e.Disable()
		e.deleted = true
	}
}
func (e *EventRequestImpl) IsEnabled() bool {
	return e.isEnabled
}
func (e *EventRequestImpl) SetEnabled(isEnable bool) {
	if e.handler == nil {
		panic("lost handler")
	}
	if e.deleted {
		panic("can't set event ")
	} else {
		if isEnable != e.isEnabled {
			if !isEnable {
				e.vm.eventRequestClear(e.GetKindType(), e.Id)
			} else {
				e.Id = e.vm.eventRequestSet(e.GetKindType(), e.suspendPolicy, e.filters)
				e.isEnabled = true
				go e.listenHandler()
			}
		}
	}
}
func (e *EventRequestImpl) Enable() {
	if e.handler == nil {
		panic("not handler func")
	}
	e.SetEnabled(true)
}
func (e *EventRequestImpl) AddCountFilter(i int) {
	e.filters = append(e.filters, jdi.CountEventModifier(i))
}
func (e *EventRequestImpl) SetSuspendPolicy(policy jdi.SuspendPolicy) {
	if e.IsEnabled() || e.deleted {
		panic("this event request has delete")
	}
	e.suspendPolicy = policy
}
func (e *EventRequestImpl) GetSuspendPolicy() jdi.SuspendPolicy {
	if e.IsEnabled() || e.deleted {
		panic("this event request has delete")
	}
	return e.suspendPolicy
}
func (e *EventRequestImpl) listenHandler() {
	events := make(chan jdi.EventResponse, 8)
	e.vm.conn.Lock()
	e.vm.conn.Events[e.Id] = events
	e.vm.conn.Unlock()
	defer func() {
		e.vm.conn.Lock()
		delete(e.vm.conn.Events, e.Id)
		e.vm.conn.Unlock()
	}()
	closeHandler := false
	for !closeHandler {
		e.vm.vmResume()
		select {
		case event := <-events:
			eventObject := translateEventToObject(event, e.vm)
			e.vm.vmSuspend()
			e.vm.FreezeVm()
			closeHandler = e.handler(eventObject)
			e.vm.UnFreezeVm()
			e.vm.vmResume()
		case <-connect.ShouldStop(e.vm.Context):
		}
	}
	e.vm.eventRequestClear(e.GetKindType(), e.Id)

	endEvent := false
	for !endEvent {
		select {
		case event := <-events:
			eventObject := translateEventToObject(event, e.vm)
			e.handler(eventObject)
		default:
			endEvent = true
		}
	}
}

func (b *BreakpointRequestImpl) GetLocation() jdi.Location {
	return b.Location
}

func (c *ClassVisibleEventRequestImpl) AddClassFilter(clazz jdi.ReferenceType) {
	if c.IsEnabled() || c.deleted {
		panic("event request has send")
	}
	c.filters = append(c.filters, jdi.ClassOnlyEventModifier(clazz.GetUniqueID()))
}
func (c *ClassVisibleEventRequestImpl) AddClassNameFilter(classPattern string) {
	if c.IsEnabled() || c.deleted {
		panic("event request has send")
	}
	c.filters = append(c.filters, jdi.ClassMatchEventModifier(classPattern))
}
func (c *ClassVisibleEventRequestImpl) AddClassExclusionFilter(classPattern string) {
	if c.IsEnabled() || c.deleted {
		panic("event request has send")
	}
	c.filters = append(c.filters, jdi.ClassExcludeEventModifier(classPattern))
}
func (c *ClassVisibleEventRequestImpl) AddInstanceFilter(instance jdi.ObjectReference) {
	if c.IsEnabled() || c.deleted {
		panic("event request has send")
	}
	c.filters = append(c.filters, jdi.InstanceOnlyEventModifier(instance.GetUniqueID()))
}

func (m *MethodExitRequestImpl) GetKindType() jdi.EventKind {
	return jdi.MethodExit
}

func (s *StepRequestImpl) GetThread() jdi.ThreadReference {
	return s.Thread
}
func (s *StepRequestImpl) GetSize() int {
	return s.Size
}
func (s *StepRequestImpl) GetDepth() int {
	return s.depth
}

func (t *ThreadVisibleEventRequestImpl) AddThreadFilter(thread jdi.ThreadReference) {
	if t.IsEnabled() || t.deleted {
		panic("event request has send")
	}
	t.filters = append(t.filters, jdi.ThreadOnlyEventModifier(thread.GetUniqueID()))
}

func (w *AccessWatchpointRequestImpl) GetField() jdi.Field {
	return w.Field
}
