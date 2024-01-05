package impl

import jdi "github.com/kyo-w/jdwp"

func isObjectTag(signature string) bool {
	tag := jdi.Tag(signature[0])
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
