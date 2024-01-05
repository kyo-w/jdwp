package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type ClassLoaderReferenceImpl struct {
	*ObjectReferenceImpl
	initDefineClasses bool
	defineClasses     []jdi.ReferenceType
}

func (c ClassLoaderReferenceImpl) GetVisibleClasses() []jdi.ReferenceType {
	return *c.classLoaderReferenceVisibleClasses(jdi.ClassObjectID(c.GetUniqueID()))
}

func (c ClassLoaderReferenceImpl) GetDefinedClasses() []jdi.ReferenceType {
	if !c.hasLockClasses() || !c.initDefineClasses {
		var out []jdi.ReferenceType
		classes := c.vmAllClasses()
		for _, value := range *classes {
			if value.IsPrepared() && value.GetClassLoader().GetUniqueID() == c.GetUniqueID() {
				out = append(out, value)
			}
		}
		c.defineClasses = out
	}
	return c.defineClasses
}
