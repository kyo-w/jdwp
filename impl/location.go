package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type LocationImpl struct {
	*MirrorImpl
	DeclaringType jdi.ReferenceType
	CodeIndex     jdi.Long
	LineNumber    jdi.Int
	method        jdi.Method
	methodId      jdi.MethodID
}

func (l *LocationImpl) GetDeclaringType() jdi.ReferenceType {
	return l.DeclaringType
}

func (l *LocationImpl) GetMethod() jdi.Method {
	return l.method
}

func (l *LocationImpl) GetCodeIndex() int64 {
	return int64(l.CodeIndex)
}

func (l *LocationImpl) GetLineNumber() int {
	return int(l.LineNumber)
}
