package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type ArrayTypeImpl struct {
	*ReferenceTypeImpl
	initComponentType bool
	componentType     jdi.Type
}

func (a *ArrayTypeImpl) NewInstance(length int) jdi.ArrayReference {
	return a.arrayTypeNewInstance(jdi.ArrayTypeID(a.TypeID), length)
}
func (a *ArrayTypeImpl) GetComponentSignature() string {
	signatureName := a.GetSignature()
	return signatureName[1:]
}
func (a *ArrayTypeImpl) GetComponentTypeName() string {
	return jdi.TranslateSignatureToClassName(a.GetSignature()[1:])
}
func (a *ArrayTypeImpl) GetComponentType() jdi.Type {
	if !a.hasLockClasses() || !a.initComponentType {
		hasFind := false
		if isObjectTag(jdi.Tag(a.GetComponentSignature()[0])) {
			signatureTypes := a.vmClassesBySignature(a.GetComponentSignature())
			for _, value := range *signatureTypes {
				if a.GetClassLoader() == value.GetClassLoader() {
					a.componentType = value
					hasFind = true
				}
			}
			if !hasFind {
				panic(a.GetSignature() + " class has not yet been loaded")
			}
		} else {
			a.componentType = a.vm.primitiveTypeMirror(jdi.Tag(a.GetComponentSignature()[0]))
		}
		a.initComponentType = true
	}
	return a.componentType
}
func (a *ArrayTypeImpl) GetAllInterfaces() []jdi.InterfaceType {
	return []jdi.InterfaceType{}
}
