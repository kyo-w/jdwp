package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type InterfaceTypeImpl struct {
	*ReferenceTypeImpl
	initCurrentSuperInterface bool
	initAllSuperInterface     bool
	initSubInterface          bool
	initImpl                  bool
	currentSuperInterfaces    []jdi.InterfaceType
	subInterfaces             []jdi.InterfaceType
	allSuperInterfaces        []jdi.InterfaceType
	allImpl                   []jdi.ClassType
}

func (i *InterfaceTypeImpl) GetSuperInterfaces() []jdi.InterfaceType {
	if !i.hasLockClasses() || !i.initCurrentSuperInterface {
		i.currentSuperInterfaces = *i.referenceTypeInterfaces(i.TypeID)
		i.initCurrentSuperInterface = true
	}
	return i.currentSuperInterfaces
}

func (i *InterfaceTypeImpl) GetSubInterfaces() []jdi.InterfaceType {
	if !i.hasLockClasses() || !i.initSubInterface {
		var out []jdi.InterfaceType
		for _, value := range *i.vmAllClasses() {
			refType, isInterface := value.(jdi.InterfaceType)
			if isInterface {
				if refType.IsPrepared() {
					allInterface := refType.GetSuperInterfaces()
					for _, interfaceElem := range allInterface {
						if interfaceElem.GetUniqueID() == i.GetUniqueID() {
							out = append(out, refType)
						}
					}
				}
			}
		}
		i.initSubInterface = true
		i.subInterfaces = out
	}
	return i.subInterfaces
}

func (i *InterfaceTypeImpl) GetAllImplementors() []jdi.ClassType {
	if !i.hasLockClasses() || !i.initImpl {
		var out []jdi.ClassType
		for _, value := range *i.vmAllClasses() {
			refType, isClassType := value.(jdi.ClassType)
			if isClassType {
				if refType.IsPrepared() {
					allInterface := refType.GetAllInterfaces()
					for _, interfaceElem := range allInterface {
						if interfaceElem.GetUniqueID() == i.GetUniqueID() {
							out = append(out, refType)
						}
					}
				}
			}
		}
		i.initImpl = true
		i.allImpl = out
	}
	return i.allImpl
}

func (i *InterfaceTypeImpl) InvokeMethod(thread jdi.ThreadReference, method jdi.Method, args []jdi.Value, options jdi.InvokeOptions) (jdi.ObjectReference, jdi.ObjectReference) {
	reqValue := make([]jdi.ValueID, len(args))
	for index, value := range args {
		reqValue[index] = value
	}
	return i.interfaceTypeInvokeMethod(jdi.InterfaceID(i.TypeID), jdi.ThreadID(thread.GetUniqueID()), jdi.MethodID(method.GetUniqueID()), reqValue, options)
}
func (i *InterfaceTypeImpl) GetAllInterfaces() []jdi.InterfaceType {
	if !i.hasLockClasses() || !i.initAllSuperInterface {
		var out []jdi.InterfaceType
		interfaces := i.GetSuperInterfaces()
		out = append(out, interfaces...)
		for _, value := range interfaces {
			if value != nil {
				out = append(out, value.GetAllInterfaces()...)
			}
		}
		i.initAllSuperInterface = true
		i.allSuperInterfaces = out
	}
	return i.allSuperInterfaces
}
