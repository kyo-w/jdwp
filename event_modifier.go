package jdwp

// EventModifier is the interface implemented by all event modifier types.
// These are filters on the events that are raised.
// See http://docs.oracle.com/javase/1.5.0/docs/guide/jpda/jdwp/jdwp-protocol.html#JDWP_EventRequest_Set
// for detailed descriptions and rules for each of the EventModifiers.
type EventModifier interface {
	ModKind() uint8
}

// CountEventModifier is an EventModifier that limits the number of times an
// event is fired. For example, using a CountEventModifier of 2 will only let
// two events fire.
type CountEventModifier int

// ThreadOnlyEventModifier is an EventModifier that filters the events to those
// that are raised on the specified thread.
type ThreadOnlyEventModifier ThreadID

// ClassOnlyEventModifier is an EventModifier that filters the events to those
// that are associated with the specified class.
type ClassOnlyEventModifier ClassID

// ClassMatchEventModifier is an EventModifier that filters the events to those
// that are associated with class names that match the pattern. The pattern can
// be an exact class name match, for use a '*' wildcard at the start or end of
// the string. Examples:
// • "java.lang.String"
// • "*.String"
// • "java.lang.*"
type ClassMatchEventModifier string

// ClassExcludeEventModifier is an EventModifier that filters the events to
// those that are not associated with class names that match the pattern.
// See ClassMatchEventModifier for the permitted patterns.
type ClassExcludeEventModifier string

// LocationOnlyEventModifier is an EventModifier that filters the events to
// those that only originate at the specified location.
type LocationOnlyEventModifier LocationID

// ExceptionOnlyEventModifier is an EventModifier that filters exception events.
// Can only be used for exception events.
type ExceptionOnlyEventModifier struct {
	ExceptionOrNull ReferenceTypeID // If not nil, only permit exceptions of this api.
	Caught          bool            // Report caught exceptions
	Uncaught        bool            // Report uncaught exceptions
}

// FieldOnlyEventModifier is an EventModifier that filters events to those
// relating to the specified field.
// Can only be used for field access or field modified events.
type FieldOnlyEventModifier struct {
	Type  ReferenceTypeID
	Field FieldID
}

// StepEventModifier is an EventModifier that filters step events to those which
// satisfy depth and size constraints.
// Can only be used with step events.
type StepEventModifier struct {
	Thread ThreadID
	Size   int
	Depth  int
}

// InstanceOnlyEventModifier is an EventModifier that filters events to those
// which have the specified 'this' object.
type InstanceOnlyEventModifier ObjectID

func (CountEventModifier) ModKind() uint8         { return 1 }
func (ThreadOnlyEventModifier) ModKind() uint8    { return 3 }
func (ClassOnlyEventModifier) ModKind() uint8     { return 4 }
func (ClassMatchEventModifier) ModKind() uint8    { return 5 }
func (ClassExcludeEventModifier) ModKind() uint8  { return 6 }
func (LocationOnlyEventModifier) ModKind() uint8  { return 7 }
func (ExceptionOnlyEventModifier) ModKind() uint8 { return 8 }
func (FieldOnlyEventModifier) ModKind() uint8     { return 9 }
func (StepEventModifier) ModKind() uint8          { return 10 }
func (InstanceOnlyEventModifier) ModKind() uint8  { return 11 }
