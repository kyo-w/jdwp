package jdwp

// Location
//
//	JVM Debug中每一个断点的逻辑位置。CodeIndex并不代指当前断点在代码中的位置
//	CodeIndex代指当前断点在字节码中的方法区的位置
//
// /*
type Location interface {
	Mirror
	GetDeclaringType() ReferenceType
	GetMethod() Method
	GetCodeIndex() int64
	GetLineNumber() int
}
