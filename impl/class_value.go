package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type ClassObjectReferenceImpl struct {
	*ObjectReferenceImpl
}

func (c *ClassObjectReferenceImpl) GetReflectedType() jdi.ReferenceType {
	return c.classObjectReferenceReflectedType(jdi.ClassObjectID(c.GetUniqueID()))
}
func (c *ClassObjectReferenceImpl) GetTagType() jdi.Tag {
	return jdi.ClassObject
}
