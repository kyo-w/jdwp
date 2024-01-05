package jdwp

type Location interface {
	Mirror
	GetDeclaringType() ReferenceType
	GetMethod() Method
	GetCodeIndex() int64
	GetLineNumber() int
}
