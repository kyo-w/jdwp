package impl

import jdi "github.com/kyo-w/jdwp"

func removenullvalue(slice []jdi.Value) *[]jdi.Value {
	var output []jdi.Value
	for _, element := range slice {
		if element != nil { //if condition satisfies add the elements in new slice
			output = append(output, element)
		}
	}
	return &output //slice with no nil-values
}

func isObjectTag(tag jdi.Tag) bool {
	return (jdi.OBJECT == tag) ||
		(tag == jdi.ARRAY) ||
		(tag == jdi.STRING) ||
		(tag == jdi.THREAD) ||
		(tag == jdi.ThreadGroup) ||
		(tag == jdi.ClassLoader) ||
		(tag == jdi.ClassObject)
}
func findType(declareType jdi.ReferenceType, signature string) jdi.Type {
	vm := declareType.GetVirtualMachine().(*VirtualMachineImpl)
	var out jdi.Type
	if len(signature) == 1 {
		getBaseType(vm, jdi.Tag(signature[0]))
	} else {
		hasFind := false
		loader := declareType.GetClassLoader()
		classes := vm.GetClassesBySignature(signature)
		for _, value := range classes {
			if (value.GetSignature() == signature) && (value.GetClassLoader().GetUniqueID() == loader.GetUniqueID()) {
				out = value
				hasFind = true
			}
		}
		if !hasFind {
			panic("can't find class : " + signature)
		}
	}
	return out
}
func getBaseType(vm *VirtualMachineImpl, signatureTag jdi.Tag) jdi.Type {
	var out jdi.Type
	if signatureTag == jdi.VOID {
		out = &jdi.VoidType{Vm: vm}
	} else {
		out = vm.primitiveTypeMirror(signatureTag)
	}
	return out
}
func invokeStaticMethod(mirrorRoot *MirrorImpl, classId jdi.ClassID, thread jdi.ThreadReference, method jdi.Method, args []jdi.Value, options jdi.InvokeOptions) (jdi.Value, jdi.ObjectReference) {
	sendArgs := make([]jdi.TaggedAny, len(args))
	for index, value := range args {
		sendArgs[index] = translateValue(value)
	}
	return mirrorRoot.classTypeInvokeMethod(classId, jdi.ThreadID(thread.GetUniqueID()),
		jdi.MethodID(method.GetUniqueID()),
		sendArgs, options)
}
func invokeObjectMethod(mirrorRoot *MirrorImpl, objectId jdi.ObjectID, thread jdi.ThreadReference, classId jdi.ClassID, method jdi.Method, args []jdi.Value, options jdi.InvokeOptions) (jdi.Value, jdi.ObjectReference) {
	sendArgs := make([]jdi.TaggedAny, len(args))
	for index, value := range args {
		sendArgs[index] = translateValue(value)
	}
	return mirrorRoot.objectReferenceInvokeMethod(objectId, jdi.ThreadID(thread.GetUniqueID()), classId, jdi.MethodID(method.GetUniqueID()),
		sendArgs, options)
}

func translateValue(value jdi.Value) jdi.TaggedAny {
	var out jdi.TaggedAny
	out.TagID = value.GetTagType()
	reference, isObject := value.(jdi.ObjectReference)
	if isObject {
		out.Value = reference.GetUniqueID()
	}
	switch value.GetTagType() {
	case jdi.BYTE:
		out.Value = value.(jdi.ByteValue).GetValue()
	case jdi.CHAR:
		out.Value = value.(jdi.CharValue).GetValue()
	case jdi.FLOAT:
		out.Value = value.(jdi.CharValue).GetValue()
	case jdi.DOUBLE:
		out.Value = value.(jdi.DoubleValue).GetValue()
	case jdi.INT:
		out.Value = value.(jdi.IntegerValue).GetValue()
	case jdi.LONG:
		out.Value = value.(jdi.LongValue).GetValue()
	case jdi.SHORT:
		out.Value = value.(jdi.ShortValue).GetValue()
	case jdi.BOOLEAN:
		out.Value = value.(jdi.BooleanValue).GetValue()
	}
	return out
}
