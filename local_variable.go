package jdwp

type LocalVariable interface {
	Mirror
	GetTypeName() string
	GetType() Type
	GetSignature() string
	GetGenericSignature() string
	// IsVisible 对于StackFrame来说是否可方法
	IsVisible(StackFrame) bool
	IsArgument() bool
	GetName() string
	GetSlot() int
}
