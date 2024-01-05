package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type LocalVariableImpl struct {
	*MirrorImpl
	Name             string
	Signature        string
	MethodRef        jdi.Method
	GenericSignature string
	StartIndex       int
	IndexLength      int
	SlotIndex        int
}

func (l *LocalVariableImpl) GetName() string {
	return l.Name
}

func (l *LocalVariableImpl) GetTypeName() string {
	return jdi.TranslateSignatureToClassName(l.Signature)
}

func (l *LocalVariableImpl) GetType() jdi.Type {
	declaringType := l.MethodRef.GetDeclaringType()
	return findType(declaringType, l.Signature)
}

func (l *LocalVariableImpl) GetSignature() string {
	return l.Signature
}

func (l *LocalVariableImpl) GetGenericSignature() string {
	return l.GenericSignature
}

func (l *LocalVariableImpl) IsVisible(frame jdi.StackFrame) bool {
	return frame.GetLocation().GetMethod().GetUniqueID() == l.MethodRef.GetUniqueID()
}

func (l *LocalVariableImpl) IsArgument() bool {
	return len(l.MethodRef.GetVariables()) > l.SlotIndex
}
func (l *LocalVariableImpl) GetSlot() int {
	return l.SlotIndex
}
