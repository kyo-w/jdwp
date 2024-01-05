package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type ClassTypeImpl struct {
	*ReferenceTypeImpl
	initSuperClass   bool
	initInterface    bool
	initAllInterface bool
	initSubClass     bool
	ownInterfaces    []jdi.InterfaceType
	allInterfaces    []jdi.InterfaceType
	superClass       jdi.ClassType
	subClass         []jdi.ClassType
}

func (c *ClassTypeImpl) GetSuperclass() jdi.ClassType {
	if !c.hasLockClasses() || !c.initSuperClass {
		c.superClass = c.classTypeSuperclass(jdi.ClassID(c.GetUniqueID()))
		c.initSuperClass = true
	}
	return c.superClass
}
func (c *ClassTypeImpl) GetOwnInterface() []jdi.InterfaceType {
	if !c.hasLockClasses() || !c.initInterface {
		c.ownInterfaces = *c.referenceTypeInterfaces(c.TypeID)
		c.initInterface = true
	}
	return c.ownInterfaces
}
func (c *ClassTypeImpl) GetAllInterfaces() []jdi.InterfaceType {
	if !c.hasLockClasses() || !c.initAllInterface {
		var out []jdi.InterfaceType
		currentInterface := c.GetOwnInterface()
		currentSuperClass := c.GetSuperclass()
		for _, iValue := range currentInterface {
			out = append(out, iValue.GetSubInterfaces()...)
		}
		if currentSuperClass != nil {
			out = append(out, currentSuperClass.GetAllInterfaces()...)
		}
		c.initAllInterface = true
		c.allInterfaces = out
	}
	return c.allInterfaces
}
func (c *ClassTypeImpl) GetSubclasses() []jdi.ClassType {
	if !c.hasLockClasses() || !c.initSubClass {
		var out []jdi.ClassType
		for _, value := range *c.vmAllClasses() {
			classTypeObject, isClassType := value.(jdi.ClassType)
			if isClassType {
				superclass := classTypeObject.GetSuperclass()
				if superclass != nil && superclass.GetUniqueID() == c.GetUniqueID() {
					out = append(out, superclass)
				}
			}
		}
		c.subClass = out
	}
	return c.subClass
}
func (c *ClassTypeImpl) IsEnum() bool {
	superclass := c.GetSuperclass()
	if superclass != nil && jdi.TranslateSignatureToClassName(superclass.GetSignature()) == "java.lang.Enum" {
		return true
	}
	return false
}
func (c *ClassTypeImpl) InvokeMethod(reference jdi.ThreadReference, method jdi.Method, args []jdi.Value, options jdi.InvokeOptions) (jdi.Value, jdi.ObjectReference) {
	argsValueId := make([]jdi.ValueID, len(args))
	for index, value := range args {
		argsValueId[index] = value
	}
	return c.classTypeInvokeMethod(jdi.ClassID(c.TypeID), jdi.ThreadID(reference.GetUniqueID()), jdi.MethodID(method.GetUniqueID()), argsValueId, options)
}
