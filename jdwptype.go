// Package jdwp Package internal
// 本代码只记录所有的JDWP通信使用的数据类型
//
// /**
package jdwp

type IntID int
type LongID int64
type ObjectID uint64
type StringID uint64
type Tag uint8
type ThreadID uint64
type ThreadGroupID uint64
type ClassLoaderID uint64
type ClassObjectID uint64
type ArrayID uint64
type ReferenceTypeID uint64
type ClassID uint64
type InterfaceID uint64
type ArrayTypeID uint64
type MethodID uint64
type FieldID uint64
type FrameID uint64
type ValueID interface{}
type ThreadStatus int

type InvokeOptions int
type ByteID byte
type TypeTag uint8
type ClassStatus int

// EventRequestID is an identifier of an event request.
type EventRequestID int
type TaggedObjectID struct {
	TagID    Tag
	ObjectID ObjectID
}
type SuspendPolicy byte

type LocationID struct {
	Type     TypeTag
	Class    ClassID
	Method   MethodID
	Location uint64
}

type IDSizes struct {
	FieldIDSize         int32 // FieldID size in bytes
	MethodIDSize        int32 // MethodID size in bytes
	ObjectIDSize        int32 // ObjectID size in bytes
	ReferenceTypeIDSize int32 // ReferenceTypeID size in bytes
	FrameIDSize         int32 // FrameID size in bytes
}

const (
	ARRAY       = Tag(91)
	BYTE        = Tag(66)
	CHAR        = Tag(67)
	OBJECT      = Tag(76)
	FLOAT       = Tag(70)
	DOUBLE      = Tag(68)
	INT         = Tag(73)
	LONG        = Tag(74)
	SHORT       = Tag(83)
	VOID        = Tag(86)
	BOOLEAN     = Tag(90)
	STRING      = Tag(115)
	THREAD      = Tag(116)
	ThreadGroup = Tag(103)
	ClassLoader = Tag(108)
	ClassObject = Tag(99)

	ClassTypeTag     = TypeTag(1) // Type is a class.
	InterfaceTypeTag = TypeTag(2) // Type is an interface.
	ArrayTypeTag     = TypeTag(3) // Type is an array.
)
