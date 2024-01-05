package impl

import jdi "github.com/kyo-w/jdwp"

func translateEventToObject(response jdi.EventResponse, vm *VirtualMachineImpl) jdi.EventObject {
	eventObject := &eventObjectImpl{Response: response, vm: vm}
	switch response.(type) {
	case *jdi.EventVMStartResponse:
		return EventVMStartResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventThreadStartResponse:
		return EventThreadStartResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventThreadDeathResponse:
		return EventThreadDeathResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventSingleStepResponse:
		return EventSingleStepResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventBreakpointResponse:
		return EventBreakpointResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventMethodEntryResponse:
		return EventMethodEntryResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventMethodExitResponse:
		return EventMethodExitResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventExceptionResponse:
		return EventExceptionResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventClassPrepareResponse:
		return EventClassPrepareResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventFieldAccessResponse:
		return EventFieldAccessResponseObject{eventObjectImpl: eventObject}
	case *jdi.EventClassUnloadResponse:
		return EventClassUnloadResponseObject{eventObjectImpl: eventObject}
	default:
		panic("unknown event object")
	}
}

type eventObjectImpl struct {
	vm       *VirtualMachineImpl
	Response jdi.EventResponse
}
type EventVMStartResponseObject struct {
	*eventObjectImpl
	thread jdi.ThreadReference
}
type EventVMDeathResponseObject struct {
	*eventObjectImpl
}
type EventThreadStartResponseObject struct {
	*eventObjectImpl
	thread jdi.ThreadReference
}
type EventThreadDeathResponseObject struct {
	*eventObjectImpl
	thread jdi.ThreadReference
}
type EventSingleStepResponseObject struct {
	*eventObjectImpl
	thread   jdi.ThreadReference
	location jdi.Location
}
type EventBreakpointResponseObject struct {
	*eventObjectImpl
	thread   jdi.ThreadReference
	location jdi.Location
}
type EventMethodEntryResponseObject struct {
	*eventObjectImpl
	thread   jdi.ThreadReference
	location jdi.Location
	method   jdi.Method
}
type EventMethodExitResponseObject struct {
	*eventObjectImpl
	method      jdi.Method
	location    jdi.Location
	thread      jdi.ThreadReference
	returnValue jdi.Value
}
type EventExceptionResponseObject struct {
	*eventObjectImpl
	thread   jdi.ThreadReference
	location jdi.Location
}
type EventClassPrepareResponseObject struct {
	*eventObjectImpl
	thread  jdi.ThreadReference
	typeRef jdi.ReferenceType
}
type EventFieldAccessResponseObject struct {
	*eventObjectImpl
	field        jdi.Field
	location     jdi.Location
	thread       jdi.ThreadReference
	objectValue  jdi.ObjectReference
	currentValue jdi.Value
}
type EventClassUnloadResponseObject struct {
	*eventObjectImpl
	signature string
}

func (e *eventObjectImpl) GetRequest() jdi.EventResponse {
	return e.Response
}

func (e *EventVMStartResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventVMStartResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}

func (e *EventThreadStartResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventThreadStartResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}

func (e *EventThreadDeathResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventThreadDeathResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}

func (e *EventSingleStepResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventSingleStepResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}
func (e *EventSingleStepResponseObject) GetLocation() jdi.Location {
	if e.location == nil {
		step := e.GetRequest().(*jdi.EventSingleStepResponse)
		referenceTypeRef := e.vm.makeReferenceTypeMirror(jdi.ReferenceTypeID(step.Location.Class), step.Location.Type, &referenceTypeInfo{})
		methodRef := e.vm.makeMethodMirror(step.Location.Method, &typeComponentInfo{DeclaringType: referenceTypeRef})
		e.location = e.vm.makeLocationMirror(&locationInfo{
			Method:        methodRef,
			DeclaringType: referenceTypeRef,
			MethodId:      step.Location.Method,
		})
	}
	return e.location
}

func (e *EventBreakpointResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventBreakpointResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}
func (e *EventBreakpointResponseObject) GetLocation() jdi.Location {
	if e.location == nil {
		breakpoint := e.GetRequest().(*jdi.EventBreakpointResponse)
		referenceTypeRef := e.vm.makeReferenceTypeMirror(jdi.ReferenceTypeID(breakpoint.Location.Class), breakpoint.Location.Type, &referenceTypeInfo{})
		methodRef := e.vm.makeMethodMirror(breakpoint.Location.Method, &typeComponentInfo{DeclaringType: referenceTypeRef})
		e.location = e.vm.makeLocationMirror(&locationInfo{
			Method:        methodRef,
			DeclaringType: referenceTypeRef,
			MethodId:      breakpoint.Location.Method,
		})
	}
	return e.location
}

func (e *EventMethodEntryResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventMethodEntryResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}
func (e *EventMethodEntryResponseObject) GetLocation() jdi.Location {
	if e.location == nil {
		methodEntry := e.GetRequest().(*jdi.EventMethodEntryResponse)
		referenceTypeRef := e.vm.makeReferenceTypeMirror(jdi.ReferenceTypeID(methodEntry.Location.Class), methodEntry.Location.Type, &referenceTypeInfo{})
		e.location = e.vm.makeLocationMirror(&locationInfo{
			Method:        e.GetMethod(),
			DeclaringType: referenceTypeRef,
			MethodId:      methodEntry.Location.Method,
		})
	}
	return e.location
}
func (e *EventMethodEntryResponseObject) GetMethod() jdi.Method {
	if e.method == nil {
		methodEntry := e.GetRequest().(*jdi.EventMethodEntryResponse)
		referenceTypeRef := e.vm.makeReferenceTypeMirror(jdi.ReferenceTypeID(methodEntry.Location.Class), methodEntry.Location.Type, &referenceTypeInfo{})
		e.method = e.vm.makeMethodMirror(methodEntry.Location.Method, &typeComponentInfo{DeclaringType: referenceTypeRef})
	}
	return e.method
}

func (e *EventMethodExitResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventMethodExitResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}
func (e *EventMethodExitResponseObject) GetLocation() jdi.Location {
	if e.location == nil {
		methodExit := e.GetRequest().(*jdi.EventMethodExitResponse)
		referenceTypeRef := e.vm.makeReferenceTypeMirror(jdi.ReferenceTypeID(methodExit.Location.Class), methodExit.Location.Type, &referenceTypeInfo{})
		e.location = e.vm.makeLocationMirror(&locationInfo{
			Method:        e.GetMethod(),
			DeclaringType: referenceTypeRef,
			MethodId:      methodExit.Location.Method,
		})
	}
	return e.location
}
func (e *EventMethodExitResponseObject) GetMethod() jdi.Method {
	if e.method == nil {
		methodExit := e.GetRequest().(*jdi.EventMethodExitResponse)
		referenceTypeRef := e.vm.makeReferenceTypeMirror(jdi.ReferenceTypeID(methodExit.Location.Class), methodExit.Location.Type, &referenceTypeInfo{})
		e.method = e.vm.makeMethodMirror(methodExit.Location.Method, &typeComponentInfo{DeclaringType: referenceTypeRef})
	}
	return e.method
}
func (e *EventMethodExitResponseObject) GetReturnValue() jdi.Value {
	//TODO implement me
	panic("implement me")
}

func (e *EventExceptionResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventExceptionResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}
func (e *EventExceptionResponseObject) GetLocation() jdi.Location {
	if e.location == nil {
		exception := e.GetRequest().(*jdi.EventExceptionResponse)
		referenceTypeRef := e.vm.makeReferenceTypeMirror(jdi.ReferenceTypeID(exception.Location.Class), exception.Location.Type, &referenceTypeInfo{})
		methodRef := e.vm.makeMethodMirror(exception.Location.Method, &typeComponentInfo{DeclaringType: referenceTypeRef})
		e.location = e.vm.makeLocationMirror(&locationInfo{
			Method:        methodRef,
			DeclaringType: referenceTypeRef,
			MethodId:      exception.Location.Method,
		})
	}
	return e.location
}

func (e *EventClassPrepareResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventClassPrepareResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}
func (e *EventClassPrepareResponseObject) GetReferenceType() jdi.ReferenceType {
	if e.typeRef == nil {
		classPrepare := e.GetRequest().(*jdi.EventClassPrepareResponse)
		e.typeRef = e.vm.makeReferenceTypeMirror(classPrepare.ClassType, classPrepare.ClassKind, &referenceTypeInfo{
			Status:        jdi.ClassStatus(classPrepare.Status),
			SignatureName: classPrepare.Signature})
	}
	return e.typeRef
}

func (e *EventClassUnloadResponseObject) GetClasName() string {
	return jdi.TranslateSignatureToClassName(e.GetClassSignature())
}
func (e *EventClassUnloadResponseObject) GetClassSignature() string {
	if e.signature == "" {
		e.signature = e.GetRequest().(*jdi.EventClassUnloadResponse).Signature
	}
	return e.signature
}

func (e *EventFieldAccessResponseObject) GetThread() jdi.ThreadReference {
	if e.thread == nil {
		e.thread = e.vm.makeObjectMirror(jdi.ObjectID(e.GetRequest().(*jdi.EventFieldAccessResponse).Thread), jdi.THREAD).(jdi.ThreadReference)
	}
	return e.thread
}
func (e *EventFieldAccessResponseObject) GetLocation() jdi.Location {
	if e.location == nil {
		accessField := e.GetRequest().(*jdi.EventFieldAccessResponse)
		referenceTypeRef := e.vm.makeReferenceTypeMirror(jdi.ReferenceTypeID(accessField.Location.Class), accessField.Location.Type, &referenceTypeInfo{})
		methodRef := e.vm.makeMethodMirror(accessField.Location.Method, &typeComponentInfo{DeclaringType: referenceTypeRef})
		e.location = e.vm.makeLocationMirror(&locationInfo{
			Method:        methodRef,
			DeclaringType: referenceTypeRef,
			MethodId:      accessField.Location.Method,
		})
	}
	return e.location
}
func (e *EventFieldAccessResponseObject) GetField() jdi.Field {
	if e.field != nil {
		accessField := e.GetRequest().(*jdi.EventFieldAccessResponse)
		e.field = e.vm.makeFieldMirror(accessField.Field, &typeComponentInfo{
			DeclaringType: e.vm.makeReferenceTypeMirror(accessField.FieldType, accessField.FieldKind, &referenceTypeInfo{}),
		})
	}
	return e.field
}
func (e *EventFieldAccessResponseObject) GetObject() jdi.ObjectReference {
	if e.objectValue == nil {
		accessField := e.GetRequest().(*jdi.EventFieldAccessResponse)
		e.objectValue = e.vm.makeObjectMirror(accessField.Object.ObjectID, accessField.Object.TagID)
	}
	return e.objectValue
}
func (e *EventFieldAccessResponseObject) GetValueCurrent() jdi.Value {
	if e.objectValue == nil {
		accessField := e.GetRequest().(*jdi.EventFieldAccessResponse)
		e.currentValue = (*e.vm.objectReferenceGetValues(accessField.Object.ObjectID, []jdi.FieldID{accessField.Field}))[0]
	}
	return e.objectValue
}
