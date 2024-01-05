package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type ArrayReferenceImpl struct {
	*ObjectReferenceImpl
	Length       int
	hasGetLength bool
}

func (a *ArrayReferenceImpl) GetLength() int {
	if !a.hasGetLength {
		a.Length = a.arrayReferenceLength(jdi.ArrayID(a.GetUniqueID()))
	}
	return a.Length
}
func (a *ArrayReferenceImpl) GetArrayValue(index int) jdi.Value {
	return (*a.arrayReferenceGetValues(jdi.ArrayID(a.GetUniqueID()), index, index+1))[0]
}
func (a *ArrayReferenceImpl) GetArrayValues() []jdi.Value {
	arrayLen := a.GetLength()
	return *a.arrayReferenceGetValues(jdi.ArrayID(a.GetUniqueID()), 0, arrayLen)
}
func (a *ArrayReferenceImpl) GetArraySlice(index, length int) []jdi.Value {
	return *a.arrayReferenceGetValues(jdi.ArrayID(a.GetUniqueID()), index, length)
}
