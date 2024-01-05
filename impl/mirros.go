package impl

import (
	jdi "github.com/kyo-w/jdwp"
	connect "github.com/kyo-w/jdwp/impl/internal"
	"log"
	"reflect"
)

type referenceTypeInfo struct {
	SignatureName    string
	GenericSignature string
	Status           jdi.ClassStatus
}
type typeComponentInfo struct {
	Name             string
	Modifiers        int
	Signature        string
	GenericSignature string
	DeclaringType    jdi.ReferenceType
}
type localVariableInfo struct {
	Name             string
	Signature        string
	MethodRef        jdi.Method
	GenericSignature string
	StartIndex       int
	IndexLength      int
	SlotIndex        int
}
type locationInfo struct {
	DeclaringType jdi.ReferenceType
	CodeIndex     jdi.Long
	LineNumber    jdi.Int
	Method        jdi.Method
	MethodId      jdi.MethodID
}
type lines struct {
	LineCodeIndex jdi.Long
	LineNumber    jdi.Int
}
type slotGenericSignature struct {
	CodeIndex        jdi.Long
	Name             string
	Signature        string
	GenericSignature string
	Length           jdi.Int
	Slot             jdi.Int
}
type slot struct {
	CodeIndex jdi.Long
	Name      string
	Signature string
	Length    jdi.Int
	Slot      jdi.Int
}

type MirrorImpl struct {
	vm *VirtualMachineImpl
	// 当freezeVm锁定时,会默认目标VM已经不再产生新的Class, 请自行确保运行的环境不会产生新的Class
	lockClasses     bool
	classTypesCache *[]jdi.ReferenceType
	typeRefMap      map[jdi.ObjectID]jdi.ReferenceType
}

func (m *MirrorImpl) FreezeVm() {
	m.lockClasses = true
}
func (m *MirrorImpl) UnFreezeVm() {
	m.lockClasses = false
}

func (m *MirrorImpl) hasLockClasses() bool {
	return m.lockClasses
}

func (m *MirrorImpl) runCmd(cmd connect.Cmd, req interface{}, out interface{}) {
	err := m.GetConnect().SendCommand(cmd, req, out)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
func (m *MirrorImpl) createEmptyMirror() *MirrorImpl {
	return m
}
func (m *MirrorImpl) readValueID(res *[]jdi.ValueID) *[]jdi.Value {
	out := make([]jdi.Value, len(*res))
	for index, value := range *res {
		valueRef := reflect.ValueOf(value)
		switch value.(type) {
		case jdi.ArrayID:
			out[index] = &ArrayReferenceImpl{ObjectReferenceImpl: &ObjectReferenceImpl{MirrorImpl: m.createEmptyMirror(), ObjectId: jdi.ObjectID(valueRef.Interface().(jdi.ArrayID))}}
		case jdi.ObjectID:
			out[index] = &ObjectReferenceImpl{MirrorImpl: m.createEmptyMirror(), ObjectId: valueRef.Interface().(jdi.ObjectID)}
		case jdi.StringID:
			out[index] = &StringReferenceImpl{ObjectReferenceImpl: &ObjectReferenceImpl{MirrorImpl: m.createEmptyMirror(), ObjectId: jdi.ObjectID(valueRef.Interface().(jdi.StringID))}}
		case jdi.ThreadID:
			out[index] = &ThreadReferenceImpl{ObjectReferenceImpl: &ObjectReferenceImpl{MirrorImpl: m.createEmptyMirror(), ObjectId: jdi.ObjectID(valueRef.Interface().(jdi.ThreadID))}}
		case jdi.ThreadGroupID:
			out[index] = &ThreadGroupReferenceImpl{ObjectReferenceImpl: &ObjectReferenceImpl{MirrorImpl: m.createEmptyMirror(), ObjectId: jdi.ObjectID(valueRef.Interface().(jdi.ThreadGroupID))}}
		case jdi.ClassLoaderID:
			out[index] = &ClassLoaderReferenceImpl{ObjectReferenceImpl: &ObjectReferenceImpl{MirrorImpl: m.createEmptyMirror(), ObjectId: jdi.ObjectID(valueRef.Interface().(jdi.ClassLoaderID))}}
		case jdi.ClassObjectID:
			out[index] = &ClassObjectReferenceImpl{ObjectReferenceImpl: &ObjectReferenceImpl{MirrorImpl: m.createEmptyMirror(), ObjectId: jdi.ObjectID(valueRef.Interface().(jdi.ClassObjectID))}}
		case byte:
			out[index] = &ByteValueImpl{MirrorImpl: m.createEmptyMirror(), value: valueRef.Interface().(byte)}
		case jdi.Char:
			out[index] = &CharValueImpl{MirrorImpl: m.createEmptyMirror(), value: valueRef.Interface().(jdi.Char)}
		case float32:
			out[index] = &FloatValueImpl{MirrorImpl: m.createEmptyMirror(), value: jdi.Float(valueRef.Interface().(float32))}
		case float64:
			out[index] = &DoubleValueImpl{MirrorImpl: m.createEmptyMirror(), value: jdi.Double(valueRef.Interface().(float64))}
		case int:
			out[index] = &IntegerValueImpl{MirrorImpl: m.createEmptyMirror(), value: jdi.Int(valueRef.Interface().(int))}
		case int16:
			out[index] = &ShortValueImpl{MirrorImpl: m.createEmptyMirror(), value: jdi.Short(valueRef.Interface().(int16))}
		case int64:
			out[index] = &LongValueImpl{MirrorImpl: m.createEmptyMirror(), value: jdi.Long(valueRef.Interface().(int64))}
		case bool:
			out[index] = &BooleanValueImpl{MirrorImpl: m.createEmptyMirror(), value: valueRef.Interface().(bool)}
		default:
			panic("unknown type :" + reflect.TypeOf(value).Name())
		}
	}
	return &out
}
func (m *MirrorImpl) GetVirtualMachine() jdi.VirtualMachine {
	return m.vm
}
func (m *MirrorImpl) validateMirrors(mirror jdi.Mirror) {
	if mirror != nil && m.GetVirtualMachine() != mirror.GetVirtualMachine() {
		panic("虚拟机异常")
	}
}
func (m *MirrorImpl) GetConnect() *connect.Connection {
	return m.vm.conn
}

func (m *MirrorImpl) makeObjectMirror(id jdi.ObjectID, tag jdi.Tag) jdi.ObjectReference {
	if id == 0 {
		return nil
	}
	var out jdi.ObjectReference
	objectImpl := &ObjectReferenceImpl{MirrorImpl: m.createEmptyMirror(), ObjectId: id}

	switch tag {
	case jdi.OBJECT:
		out = objectImpl
	case jdi.STRING:
		out = &StringReferenceImpl{ObjectReferenceImpl: objectImpl}
	case jdi.ARRAY:
		out = &ArrayReferenceImpl{ObjectReferenceImpl: objectImpl}
	case jdi.THREAD:
		out = &ThreadReferenceImpl{ObjectReferenceImpl: objectImpl}
	case jdi.ThreadGroup:
		out = &ThreadGroupReferenceImpl{ObjectReferenceImpl: objectImpl}
	case jdi.ClassLoader:
		out = &ClassLoaderReferenceImpl{ObjectReferenceImpl: objectImpl}
	case jdi.ClassObject:
		out = &ClassObjectReferenceImpl{ObjectReferenceImpl: objectImpl}
	default:
		panic("Invalid object tag: " + string(tag))
	}
	return out
}
func (m *MirrorImpl) makeReferenceTypeMirror(id jdi.ReferenceTypeID, tag jdi.TypeTag, info *referenceTypeInfo) jdi.ReferenceType {
	if id == 0 {
		return nil
	}
	var out jdi.ReferenceType
	refImpl := &ReferenceTypeImpl{
		Kind:             tag,
		MirrorImpl:       m.createEmptyMirror(),
		TypeID:           id,
		signatureName:    info.SignatureName,
		genericSignature: info.GenericSignature,
		status:           info.Status}
	switch tag {
	case jdi.ClassTypeTag:
		out = &ClassTypeImpl{ReferenceTypeImpl: refImpl}
	case jdi.InterfaceTypeTag:
		out = &InterfaceTypeImpl{ReferenceTypeImpl: refImpl}
	case jdi.ArrayTypeTag:
		out = &ArrayTypeImpl{ReferenceTypeImpl: refImpl}
	default:
		panic("unknown ReferenceType Tag: " + string(tag))
	}
	return out
}
func (m *MirrorImpl) makeFieldMirror(id jdi.FieldID, info *typeComponentInfo) jdi.Field {
	if id == 0 {
		return nil
	}
	out := &FieldImpl{TypeComponentImpl: &TypeComponentImpl{
		MirrorImpl:       m.createEmptyMirror(),
		RefId:            jdi.ObjectID(id),
		Name:             info.Name,
		signature:        info.Signature,
		genericSignature: info.GenericSignature,
		Modifiers:        info.Modifiers,
		DeclaringType:    info.DeclaringType,
	}}
	return out
}
func (m *MirrorImpl) makeMethodMirror(id jdi.MethodID, info *typeComponentInfo) jdi.Method {
	if id == 0 {
		return nil
	}
	out := &MethodImpl{TypeComponentImpl: &TypeComponentImpl{
		MirrorImpl:       m.createEmptyMirror(),
		RefId:            jdi.ObjectID(id),
		Name:             info.Name,
		signature:        info.Signature,
		genericSignature: info.GenericSignature,
		Modifiers:        info.Modifiers,
		DeclaringType:    info.DeclaringType,
	}}
	return out
}
func (m *MirrorImpl) makeLocalVariableMirror(localVar *localVariableInfo) jdi.LocalVariable {
	return &LocalVariableImpl{MirrorImpl: m.createEmptyMirror(),
		Name:             localVar.Name,
		Signature:        localVar.Signature,
		MethodRef:        localVar.MethodRef,
		GenericSignature: localVar.GenericSignature,
		StartIndex:       localVar.StartIndex,
		IndexLength:      localVar.IndexLength,
		SlotIndex:        localVar.SlotIndex,
	}
}
func (m *MirrorImpl) makeLocationMirror(location *locationInfo) jdi.Location {
	return &LocationImpl{MirrorImpl: m.createEmptyMirror(),
		DeclaringType: location.DeclaringType,
		CodeIndex:     location.CodeIndex,
		method:        location.Method,
		methodId:      location.MethodId,
		LineNumber:    location.LineNumber}
}

func (m *MirrorImpl) vmGetVersion() *jdi.VmVersion {
	out := &jdi.VmVersion{}
	m.runCmd(connect.CmdVirtualMachineVersion, struct{}{}, out)
	return out
}
func (m *MirrorImpl) vmClassesBySignature(signature string) *[]jdi.ReferenceType {
	var res []struct {
		Tag    jdi.TypeTag
		TypeID jdi.ReferenceTypeID
		Status jdi.ClassStatus
	}
	m.runCmd(connect.CmdVirtualMachineClassesBySignature, &signature, &res)
	out := make([]jdi.ReferenceType, len(res))
	for index, value := range res {
		out[index] = m.makeReferenceTypeMirror(value.TypeID, value.Tag, &referenceTypeInfo{Status: value.Status})
	}
	return &out
}
func (m *MirrorImpl) vmAllClasses() *[]jdi.ReferenceType {
	//if !m.freezeVm || m.classTypesCache == nil {
	var res []struct {
		Tag       jdi.TypeTag
		TypeID    jdi.ReferenceTypeID
		Signature string
		Status    jdi.ClassStatus
	}
	m.runCmd(connect.CmdVirtualMachineAllClasses, struct{}{}, &res)
	out := make([]jdi.ReferenceType, len(res))
	for index, value := range res {
		out[index] = m.makeReferenceTypeMirror(value.TypeID, value.Tag, &referenceTypeInfo{Status: value.Status, SignatureName: value.Signature})
	}
	m.classTypesCache = &out
	//}
	return m.classTypesCache
}
func (m *MirrorImpl) vmAllThreads() *[]jdi.ThreadReference {
	var res []jdi.ThreadID
	m.runCmd(connect.CmdVirtualMachineAllThreads, struct{}{}, &res)
	out := make([]jdi.ThreadReference, len(res))
	for index, value := range res {
		out[index] = m.makeObjectMirror(jdi.ObjectID(value), jdi.THREAD).(jdi.ThreadReference)
	}
	return &out
}
func (m *MirrorImpl) vmTopLevelThreadGroups() *[]jdi.ThreadGroupReference {
	var res []jdi.ThreadGroupID
	m.runCmd(connect.CmdVirtualMachineTopLevelThreadGroups, struct{}{}, &res)
	out := make([]jdi.ThreadGroupReference, len(res))
	for index, value := range res {
		out[index] = &ThreadGroupReferenceImpl{ObjectReferenceImpl: &ObjectReferenceImpl{ObjectId: jdi.ObjectID(value)}}
	}
	return &out
}
func (m *MirrorImpl) vmDispose() {
	m.runCmd(connect.CmdVirtualMachineDispose, struct{}{}, struct{}{})
}
func (m *MirrorImpl) vmExit(exitCode int) {
	m.runCmd(connect.CmdVirtualMachineExit, &exitCode, struct{}{})
}
func (m *MirrorImpl) vmIDSize() *jdi.IDSizes {
	out := jdi.IDSizes{}
	m.runCmd(connect.CmdVirtualMachineIDSizes, struct{}{}, out)
	return &out
}
func (m *MirrorImpl) vmSuspend() {
	m.runCmd(connect.CmdVirtualMachineSuspend, struct{}{}, struct{}{})
}
func (m *MirrorImpl) vmCreateString(str string) jdi.StringReference {
	res := jdi.StringID(0)
	m.runCmd(connect.CmdVirtualMachineCreateString, str, &res)

	return m.makeObjectMirror(jdi.ObjectID(res), jdi.STRING).(jdi.StringReference)
}
func (m *MirrorImpl) vmCapabilities() *jdi.Capabilities {
	return m.vmCapabilitiesNew()
}
func (m *MirrorImpl) vmClassPaths() *jdi.ClassPath {
	out := &jdi.ClassPath{}
	m.runCmd(connect.CmdVirtualMachineClassPaths, struct{}{}, out)
	return out
}
func (m *MirrorImpl) vmHoldEvents() {
	m.runCmd(connect.CmdVirtualMachineHoldEvents, struct{}{}, struct{}{})
}
func (m *MirrorImpl) vmReleaseEvents() {
	m.runCmd(connect.CmdVirtualMachineReleaseEvents, struct{}{}, struct{}{})
}
func (m *MirrorImpl) vmResume() {
	m.runCmd(connect.CmdVirtualMachineResume, struct{}{}, struct{}{})
}
func (m *MirrorImpl) vmCapabilitiesNew() *jdi.Capabilities {
	out := &jdi.Capabilities{}
	m.runCmd(connect.CmdVirtualMachineCapabilitiesNew, struct{}{}, out)
	return out
}
func (m *MirrorImpl) vmAllClassesWithGeneric() *[]jdi.ReferenceType {
	var res []struct {
		Tag              jdi.TypeTag
		TypeID           jdi.ReferenceTypeID
		Signature        string
		GenericSignature string
		Status           jdi.ClassStatus
	}
	m.runCmd(connect.CmdVirtualMachineAllClassesWithGeneric, struct{}{}, &res)
	out := make([]jdi.ReferenceType, len(res))
	for index, value := range res {
		out[index] = m.makeReferenceTypeMirror(value.TypeID, value.Tag, &referenceTypeInfo{
			Status:           value.Status,
			SignatureName:    value.Signature,
			GenericSignature: value.GenericSignature,
		})
	}
	return &out
}

func (m *MirrorImpl) referenceTypeSignature(id jdi.ReferenceTypeID) string {
	var signature string
	m.runCmd(connect.CmdReferenceTypeSignature, id, &signature)
	return signature
}
func (m *MirrorImpl) referenceTypeClassLoader(id jdi.ReferenceTypeID) jdi.ClassLoaderReference {
	var classLoaderId jdi.ClassLoaderID
	m.runCmd(connect.CmdReferenceTypeClassLoader, id, &classLoaderId)
	out := m.makeObjectMirror(jdi.ObjectID(classLoaderId), jdi.ClassLoader)
	if out == nil {
		return nil
	}
	return out.(jdi.ClassLoaderReference)
}
func (m *MirrorImpl) referenceTypeModifiers(id jdi.ReferenceTypeID) int {
	var out int
	m.runCmd(connect.CmdReferenceTypeModifiers, id, &out)
	return out
}
func (m *MirrorImpl) referenceTypeFields(thisType jdi.ReferenceType, id jdi.ReferenceTypeID) *[]jdi.Field {
	var res []struct {
		FieldID   jdi.FieldID
		Name      string
		Signature string
		ModBits   int
	}
	m.runCmd(connect.CmdReferenceTypeFields, id, &res)
	out := make([]jdi.Field, len(res))
	for index, value := range res {
		out[index] = m.makeFieldMirror(value.FieldID, &typeComponentInfo{
			Name:          value.Name,
			Signature:     value.Signature,
			Modifiers:     value.ModBits,
			DeclaringType: thisType,
		})
	}
	return &out
}
func (m *MirrorImpl) referenceTypeMethods(thisTypeRef jdi.ReferenceType, id jdi.ReferenceTypeID) *[]jdi.Method {
	var res []struct {
		MethodID  jdi.MethodID
		Name      string
		Signature string
		ModBits   int
	}
	m.runCmd(connect.CmdReferenceTypeMethods, id, &res)
	out := make([]jdi.Method, len(res))
	for index, value := range res {
		out[index] = m.makeMethodMirror(value.MethodID, &typeComponentInfo{
			Name:          value.Name,
			Signature:     value.Signature,
			Modifiers:     value.ModBits,
			DeclaringType: thisTypeRef,
		})
	}
	return &out
}
func (m *MirrorImpl) referenceTypeGetValues(id jdi.ReferenceTypeID, fields []jdi.FieldID) *[]jdi.Value {
	req := struct {
		ID     jdi.ReferenceTypeID
		Fields []jdi.FieldID
	}{ID: id, Fields: fields}
	var res []jdi.ValueID
	m.runCmd(connect.CmdReferenceTypeGetValues, &req, &res)
	return m.readValueID(&res)
}
func (m *MirrorImpl) referenceTypeSourceFile(id jdi.ReferenceTypeID) string {
	var out string
	m.runCmd(connect.CmdReferenceTypeSourceFile, id, &out)
	return out
}
func (m *MirrorImpl) referenceTypeNestedTypes(id jdi.ReferenceTypeID) *[]jdi.ReferenceType {
	var res []struct {
		Tag    jdi.TypeTag
		TypeId jdi.ReferenceTypeID
	}
	m.runCmd(connect.CmdReferenceTypeNestedTypes, id, &res)
	out := make([]jdi.ReferenceType, len(res))
	for index, value := range res {
		out[index] = m.makeReferenceTypeMirror(value.TypeId, value.Tag, &referenceTypeInfo{})
	}
	return &out
}
func (m *MirrorImpl) referenceTypeStatus(id jdi.ReferenceTypeID) jdi.ClassStatus {
	var out jdi.ClassStatus
	m.runCmd(connect.CmdReferenceTypeStatus, id, &out)
	return out
}
func (m *MirrorImpl) referenceTypeInterfaces(id jdi.ReferenceTypeID) *[]jdi.InterfaceType {
	var res []jdi.InterfaceID
	m.runCmd(connect.CmdReferenceTypeInterfaces, id, &res)
	out := make([]jdi.InterfaceType, len(res))
	for index, value := range res {
		out[index] = m.makeReferenceTypeMirror(jdi.ReferenceTypeID(value), jdi.InterfaceTypeTag, &referenceTypeInfo{}).(jdi.InterfaceType)
	}
	return &out
}
func (m *MirrorImpl) referenceClassObject(id jdi.ReferenceTypeID) *[]jdi.ClassObjectReference {
	var res []jdi.ClassObjectID
	m.runCmd(connect.CmdReferenceTypeInterfaces, id, &res)
	out := make([]jdi.ClassObjectReference, len(res))
	for index, value := range res {
		out[index] = m.makeObjectMirror(jdi.ObjectID(value), jdi.ClassObject).(jdi.ClassObjectReference)
	}
	return &out
}
func (m *MirrorImpl) referenceSignatureWithGeneric(id jdi.ReferenceTypeID) (string, string) {
	var res struct {
		Signature        string
		GenericSignature string
	}
	m.runCmd(connect.CmdReferenceTypeSignatureWithGeneric, id, &res)
	return res.Signature, res.GenericSignature
}
func (m *MirrorImpl) referenceSourceDebugExtension(id jdi.ReferenceTypeID) string {
	var out string
	m.runCmd(connect.CmdReferenceTypeSourceDebugExtension, id, &out)
	return out
}
func (m *MirrorImpl) referenceFieldsWithGeneric(thisTypeRef jdi.ReferenceType, id jdi.ReferenceTypeID) *[]jdi.Field {
	var res []struct {
		FieldID          jdi.FieldID
		Name             string
		Signature        string
		GenericSignature string
		ModBits          int
	}
	m.runCmd(connect.CmdReferenceTypeFieldsWithGeneric, id, &res)
	out := make([]jdi.Field, len(res))
	for index, value := range res {
		out[index] = m.makeFieldMirror(value.FieldID, &typeComponentInfo{
			Name:             value.Name,
			Signature:        value.Name,
			GenericSignature: value.Name,
			Modifiers:        value.ModBits,
			DeclaringType:    thisTypeRef,
		})
	}
	return &out
}
func (m *MirrorImpl) referenceMethodsWithGeneric(thisRef jdi.ReferenceType, id jdi.ReferenceTypeID) *[]jdi.Method {
	var res []struct {
		MethodID         jdi.MethodID
		Name             string
		Signature        string
		GenericSignature string
		ModBits          int
	}
	m.runCmd(connect.CmdReferenceTypeMethodsWithGeneric, id, &res)
	out := make([]jdi.Method, len(res))
	for index, value := range res {
		out[index] = m.makeMethodMirror(value.MethodID, &typeComponentInfo{
			Name:             value.Name,
			Signature:        value.Name,
			GenericSignature: value.Name,
			Modifiers:        value.ModBits,
			DeclaringType:    thisRef,
		})
	}
	return &out
}
func (m *MirrorImpl) referenceInstances(id jdi.ReferenceTypeID, max int) *[]jdi.ObjectReference {
	var req = struct {
		TypeID jdi.ReferenceTypeID
		Max    int
	}{id, max}
	var res []jdi.TaggedObjectID
	m.runCmd(connect.CmdReferenceTypeInstances, &req, &res)
	out := make([]jdi.ObjectReference, len(res))
	for index, value := range res {
		out[index] = m.makeObjectMirror(value.ObjectID, value.TagID)
	}
	return &out
}
func (m *MirrorImpl) referenceClassFileVersion(id jdi.ReferenceTypeID) (jdi.Int, jdi.Int) {
	var res struct {
		MajorVersion jdi.Int
		MinorVersion jdi.Int
	}
	m.runCmd(connect.CmdReferenceTypeClassFileVersion, id, &res)
	return res.MajorVersion, res.MinorVersion
}
func (m *MirrorImpl) referenceConstantPool(id jdi.ReferenceTypeID) []byte {
	var res struct {
		ConstantCount int
		Info          []byte
	}
	m.runCmd(connect.CmdReferenceTypeConstantPool, id, &res)
	return res.Info
}

// func (m *MirrorImpl) classTypeSetValues(id jdi.ClassID)
func (m *MirrorImpl) classTypeSuperclass(id jdi.ClassID) jdi.ClassType {
	var res jdi.ClassID
	m.runCmd(connect.CmdClassTypeSuperclass, id, &res)
	return m.makeReferenceTypeMirror(jdi.ReferenceTypeID(res), jdi.ClassTypeTag, &referenceTypeInfo{}).(jdi.ClassType)
}
func (m *MirrorImpl) classTypeInvokeMethod(classId jdi.ClassID, threadId jdi.ThreadID, methodId jdi.MethodID, args []jdi.ValueID, options jdi.InvokeOptions) (jdi.Value, jdi.ObjectReference) {
	req := struct {
		Class   jdi.ClassID
		Thread  jdi.ThreadID
		Method  jdi.MethodID
		Args    []jdi.ValueID
		Options jdi.InvokeOptions
	}{classId, threadId, methodId, args, options}

	var res struct {
		Result    jdi.ValueID
		Exception jdi.TaggedObjectID
	}
	m.runCmd(connect.CmdClassTypeInvokeMethod, &req, &res)
	valueOut := (*m.readValueID(&[]jdi.ValueID{res.Result}))[1]
	return valueOut, m.makeObjectMirror(res.Exception.ObjectID, res.Exception.TagID)
}
func (m *MirrorImpl) classTypeNewInstance(classId jdi.ClassID, threadId jdi.ThreadID, constructor jdi.MethodID, args []jdi.ValueID, options jdi.InvokeOptions) (jdi.ObjectReference, jdi.ObjectReference) {
	req := struct {
		Class       jdi.ClassID
		Thread      jdi.ThreadID
		Constructor jdi.MethodID
		Args        []jdi.ValueID
		Options     jdi.InvokeOptions
	}{classId, threadId, constructor, args, options}
	var res struct {
		NewObject jdi.TaggedObjectID
		Exception jdi.TaggedObjectID
	}
	m.runCmd(connect.CmdClassTypeNewInstance, &req, &res)
	return m.makeObjectMirror(res.NewObject.ObjectID, res.NewObject.TagID), m.makeObjectMirror(res.Exception.ObjectID, res.Exception.TagID)
}
func (m *MirrorImpl) arrayTypeNewInstance(id jdi.ArrayTypeID, length int) jdi.ArrayReference {
	req := struct {
		ArrayId jdi.ArrayTypeID
		Length  int
	}{id, length}
	var out jdi.TaggedObjectID
	m.runCmd(connect.CmdArrayTypeNewInstance, &req, &out)
	return m.makeObjectMirror(out.ObjectID, out.TagID).(jdi.ArrayReference)
}
func (m *MirrorImpl) interfaceTypeInvokeMethod(interfaceId jdi.InterfaceID, threadId jdi.ThreadID, methodId jdi.MethodID, args []jdi.ValueID, options jdi.InvokeOptions) (jdi.ObjectReference, jdi.ObjectReference) {
	req := struct {
		InterfaceId jdi.InterfaceID
		ThreadId    jdi.ThreadID
		MethodId    jdi.MethodID
		Args        []jdi.ValueID
		Options     jdi.InvokeOptions
	}{interfaceId, threadId, methodId, args, options}
	var res struct {
		ReturnValue jdi.TaggedObjectID
		Exception   jdi.TaggedObjectID
	}
	m.runCmd(connect.CmdInterfaceTypeInvokeMethod, &req, &res)
	return m.makeObjectMirror(res.ReturnValue.ObjectID, res.ReturnValue.TagID), m.makeObjectMirror(res.Exception.ObjectID, res.Exception.TagID)
}

// method实现不全!
// Bytecodes
// IsObsolete
// VariableTableWithGeneric
func (m *MirrorImpl) methodTypeLineTable(typeRef jdi.ReferenceType, method jdi.Method) (jdi.Long, jdi.Long, []jdi.Location) {
	var req = struct {
		ReferenceTypeID jdi.ReferenceTypeID
		MethodId        jdi.MethodID
	}{typeRef.GetUniqueID(), jdi.MethodID(method.GetUniqueID())}
	var res struct {
		Start jdi.Long
		End   jdi.Long
		Lines []lines
	}
	m.runCmd(connect.CmdMethodTypeLineTable, &req, &res)
	out := make([]jdi.Location, len(res.Lines))
	for index, value := range res.Lines {
		out[index] = m.makeLocationMirror(&locationInfo{
			DeclaringType: typeRef,
			MethodId:      jdi.MethodID(method.GetUniqueID()),
			Method:        method,
			CodeIndex:     value.LineCodeIndex,
			LineNumber:    value.LineNumber,
		})
	}
	return res.Start, res.End, out
}
func (m *MirrorImpl) methodTypeVariableTable(refTypeId jdi.ReferenceTypeID, method jdi.Method) (jdi.Int, []jdi.LocalVariable) {
	var req = struct {
		RefType  jdi.ReferenceTypeID
		MethodId jdi.MethodID
	}{refTypeId, jdi.MethodID(method.GetUniqueID())}
	var res struct {
		ArgCnt jdi.Int
		Slots  []slot
	}
	m.runCmd(connect.CmdMethodTypeVariableTable, &req, &res)
	out := make([]jdi.LocalVariable, len(res.Slots))
	for index, value := range res.Slots {
		out[index] = m.makeLocalVariableMirror(&localVariableInfo{
			MethodRef:   method,
			Name:        value.Name,
			Signature:   value.Signature,
			StartIndex:  int(value.CodeIndex),
			IndexLength: int(value.Length),
			SlotIndex:   int(value.Slot)})
	}
	return res.ArgCnt, out
}
func (m *MirrorImpl) methodTypeIsObsolete(refTypeId jdi.ReferenceTypeID, method jdi.MethodID) bool {
	var req = struct {
		RefType  jdi.ReferenceTypeID
		MethodId jdi.MethodID
	}{refTypeId, method}
	var out bool
	m.runCmd(connect.CmdMethodTypeIsObsolete, &req, &out)
	return out
}

func (m *MirrorImpl) VariableTableWithGeneric(refTypeId jdi.ReferenceTypeID, method jdi.Method) []jdi.LocalVariable {
	var req = struct {
		RefType  jdi.ReferenceTypeID
		MethodId jdi.MethodID
	}{refTypeId, jdi.MethodID(method.GetUniqueID())}
	var res struct {
		ArgCnt jdi.Int
		Slots  []slotGenericSignature
	}
	m.runCmd(connect.CmdMethodTypeVariableTable, &req, &res)
	out := make([]jdi.LocalVariable, len(res.Slots))
	for index, value := range res.Slots {
		out[index] = m.makeLocalVariableMirror(&localVariableInfo{
			MethodRef:        method,
			Name:             value.Name,
			Signature:        value.Signature,
			GenericSignature: value.GenericSignature,
			StartIndex:       int(value.CodeIndex),
			IndexLength:      int(value.Length),
			SlotIndex:        int(value.Slot),
		})
	}
	return out
}

// func (m *MirrorImpl) objectReferenceSetValues(id jdi.ObjectID) {}
// func (m *MirrorImpl)objectReferenceMonitorInfo(){}
func (m *MirrorImpl) objectReferenceReferenceType(id jdi.ObjectID) jdi.ReferenceType {
	var out struct {
		RefTypeTag jdi.TypeTag
		TypeID     jdi.ReferenceTypeID
	}
	m.runCmd(connect.CmdObjectReferenceReferenceType, id, &out)
	return m.makeReferenceTypeMirror(out.TypeID, out.RefTypeTag, &referenceTypeInfo{})
}
func (m *MirrorImpl) objectReferenceGetValues(id jdi.ObjectID, fields []jdi.FieldID) *[]jdi.Value {
	var req = struct {
		ObjectID jdi.ObjectID
		Fields   []jdi.FieldID
	}{id, fields}
	var res []jdi.ValueID
	m.runCmd(connect.CmdObjectReferenceGetValues, &req, &res)
	return m.readValueID(&res)
}
func (m *MirrorImpl) objectReferenceInvokeMethod(id jdi.ObjectID, threadId jdi.ThreadID, classId jdi.ClassID, methodId jdi.MethodID, args []jdi.ValueID, options jdi.InvokeOptions) (jdi.Value, jdi.ObjectReference) {
	var req = struct {
		ObjectId jdi.ObjectID
		ThreadId jdi.ThreadID
		ClassId  jdi.ClassID
		MethodId jdi.MethodID
		Args     []jdi.ValueID
		Options  jdi.InvokeOptions
	}{id, threadId, classId, methodId, args, options}
	var res struct {
		Result    jdi.ValueID
		Exception jdi.TaggedObjectID
	}
	m.runCmd(connect.CmdObjectReferenceInvokeMethod, &req, &res)
	valueOut := (*m.readValueID(&[]jdi.ValueID{res.Result}))[0]
	return valueOut, m.makeObjectMirror(res.Exception.ObjectID, res.Exception.TagID)
}
func (m *MirrorImpl) objectReferenceDisableCollection(id jdi.ObjectID) {
	m.runCmd(connect.CmdObjectReferenceDisableCollection, id, struct{}{})
}
func (m *MirrorImpl) objectReferenceEnableCollection(id jdi.ObjectID) {
	m.runCmd(connect.CmdObjectReferenceEnableCollection, id, struct{}{})
}
func (m *MirrorImpl) objectReferenceIsCollected(id jdi.ObjectID) bool {
	var out bool
	m.runCmd(connect.CmdObjectReferenceIsCollected, id, &out)
	return out
}
func (m *MirrorImpl) objectReferenceReferringObjects(id jdi.ObjectID, max int) *[]jdi.ObjectReference {
	var req = struct {
		ObjectId jdi.ObjectID
		Max      int
	}{id, max}
	var res []jdi.TaggedObjectID
	m.runCmd(connect.CmdObjectReferringObjects, &req, &res)
	out := make([]jdi.ObjectReference, len(res))
	for index, value := range res {
		out[index] = m.makeObjectMirror(value.ObjectID, value.TagID)
	}
	return &out
}

func (m *MirrorImpl) stringReferenceValue(id jdi.StringID) string {
	var out string
	m.runCmd(connect.CmdStringReferenceValue, id, &out)
	return out
}

// func (m *MirrorImpl) threadReferenceOwnedMonitorsStackDepthInfo(id jdi.ThreadID) {}
// func (m *MirrorImpl) threadReferenceForceEarlyReturn(id jdi.ThreadID) {}
// func (m *MirrorImpl) threadReferenceOwnedMonitors(id jdi.ThreadID){}
// func (m *MirrorImpl) threadReferenceCurrentContendedMonitor(id jdi.ThreadID) {}
func (m *MirrorImpl) threadReferenceName(id jdi.ThreadID) string {
	var out string
	m.runCmd(connect.CmdThreadReferenceName, id, &out)
	return out
}
func (m *MirrorImpl) threadReferenceSuspend(id jdi.ThreadID) {
	m.runCmd(connect.CmdThreadReferenceSuspend, id, struct{}{})
}
func (m *MirrorImpl) threadReferenceResume(id jdi.ThreadID) {
	m.runCmd(connect.CmdThreadReferenceResume, id, struct{}{})
}
func (m *MirrorImpl) threadReferenceStatus(id jdi.ThreadID) jdi.ThreadStatus {
	var out jdi.ThreadStatus
	m.runCmd(connect.CmdThreadReferenceStatus, id, &out)
	return out
}
func (m *MirrorImpl) threadReferenceThreadGroup(id jdi.ThreadID) jdi.ThreadGroupReference {
	var res jdi.ThreadGroupID
	m.runCmd(connect.CmdThreadReferenceThreadGroup, id, &res)
	return m.makeObjectMirror(jdi.ObjectID(res), jdi.ThreadGroup).(jdi.ThreadGroupReference)
}
func (m *MirrorImpl) threadReferenceFrames(thread jdi.ThreadReference, startFrame, length int) []jdi.StackFrame {
	var req = struct {
		ThreadId        jdi.ThreadID
		StackFrameIndex jdi.Int
		Length          jdi.Int
	}{jdi.ThreadID(thread.GetUniqueID()), jdi.Int(startFrame), jdi.Int(length)}
	var res []struct {
		FrameId   jdi.FrameID
		Tag       jdi.TypeTag
		ClassRef  jdi.Long
		MethodId  jdi.MethodID
		CodeIndex jdi.Long
	}
	m.runCmd(connect.CmdThreadReferenceFrames, &req, &res)
	out := make([]jdi.StackFrame, len(res))
	for index, value := range res {
		if value.FrameId == 0 {
			out[index] = nil
		} else {
			mirror := m.createEmptyMirror()

			out[index] = &StackFrameImpl{MirrorImpl: mirror,
				StackFrameId: value.FrameId,
				ThreadRef:    thread,
				Location: &LocationImpl{MirrorImpl: mirror,
					CodeIndex:     value.CodeIndex,
					methodId:      value.MethodId,
					DeclaringType: m.makeReferenceTypeMirror(jdi.ReferenceTypeID(value.ClassRef), value.Tag, &referenceTypeInfo{})}}
		}
	}
	return out
}
func (m *MirrorImpl) threadReferenceFrameCount(id jdi.ThreadID) int {
	var out jdi.Int
	m.runCmd(connect.CmdThreadReferenceFrameCount, id, &out)
	return int(out)
}
func (m *MirrorImpl) threadReferenceStop(threadId jdi.ThreadID, throwableId jdi.ObjectID) {
	var req = struct {
		ThreadId jdi.ThreadID
		ThrowId  jdi.ObjectID
	}{threadId, throwableId}
	m.runCmd(connect.CmdThreadReferenceStop, &req, struct{}{})
}
func (m *MirrorImpl) threadReferenceInterrupt(id jdi.ThreadID) {
	m.runCmd(connect.CmdThreadReferenceInterrupt, id, struct{}{})
}
func (m *MirrorImpl) threadReferenceSuspendCount(id jdi.ThreadID) int {
	var out jdi.Int
	m.runCmd(connect.CmdThreadReferenceSuspendCount, id, &out)
	return int(out)
}

func (m *MirrorImpl) threadGroupReferenceName(id jdi.ThreadGroupID) string {
	var out string
	m.runCmd(connect.CmdThreadGroupReferenceName, id, &out)
	return out
}
func (m *MirrorImpl) threadGroupReferenceParent(id jdi.ThreadGroupID) jdi.ThreadGroupReference {
	var res jdi.ThreadGroupID
	m.runCmd(connect.CmdThreadGroupReferenceParent, id, &res)
	return m.makeObjectMirror(jdi.ObjectID(res), jdi.ThreadGroup).(jdi.ThreadGroupReference)
}
func (m *MirrorImpl) threadGroupReferenceChildren(id jdi.ThreadGroupID) ([]jdi.ThreadReference, []jdi.ThreadGroupReference) {
	var res struct {
		ChildThreads []jdi.ThreadID
		ChildGroups  []jdi.ThreadGroupID
	}
	m.runCmd(connect.CmdThreadGroupReferenceChildren, id, &res)
	threadOut := make([]jdi.ThreadReference, len(res.ChildThreads))
	groupsOut := make([]jdi.ThreadGroupReference, len(res.ChildGroups))

	for index, value := range res.ChildThreads {
		threadOut[index] = m.makeObjectMirror(jdi.ObjectID(value), jdi.THREAD).(jdi.ThreadReference)
	}

	for index, value := range res.ChildGroups {
		groupsOut[index] = m.makeObjectMirror(jdi.ObjectID(value), jdi.ThreadGroup).(jdi.ThreadGroupReference)
	}
	return threadOut, groupsOut
}

// func (m *MirrorImpl)arrayReferenceSetValues(){}
func (m *MirrorImpl) arrayReferenceLength(id jdi.ArrayID) int {
	var out jdi.Int
	m.runCmd(connect.CmdArrayReferenceLength, id, &out)
	return int(out)
}
func (m *MirrorImpl) arrayReferenceGetValues(id jdi.ArrayID, first, length int) *[]jdi.Value {
	var req = struct {
		ArrayId    jdi.ArrayID
		FirstIndex jdi.Int
		Length     jdi.Int
	}{id, jdi.Int(first), jdi.Int(length)}
	var res []jdi.ValueID
	m.runCmd(connect.CmdArrayReferenceGetValues, &req, &res)
	return m.readValueID(&res)
}

func (m *MirrorImpl) classLoaderReferenceVisibleClasses(id jdi.ClassObjectID) *[]jdi.ReferenceType {
	var res []struct {
		RefTypeTag jdi.TypeTag
		TypeID     jdi.ReferenceTypeID
	}
	m.runCmd(connect.CmdClassLoaderReferenceVisibleClasses, id, &res)
	out := make([]jdi.ReferenceType, len(res))
	for index, value := range res {
		out[index] = m.makeReferenceTypeMirror(value.TypeID, value.RefTypeTag, &referenceTypeInfo{}).(jdi.ReferenceType)
	}
	return &out
}

func (m *MirrorImpl) eventRequestSet(kind jdi.EventKind, suspendPolicy jdi.SuspendPolicy, modifiers []jdi.EventModifier) jdi.EventRequestID {
	req := struct {
		Kind          jdi.EventKind
		SuspendPolicy jdi.SuspendPolicy
		Modifiers     []jdi.EventModifier
	}{kind, suspendPolicy, modifiers}
	var out jdi.EventRequestID
	m.runCmd(connect.CmdEventRequestSet, &req, &out)
	return out
}
func (m *MirrorImpl) eventRequestClear(kind jdi.EventKind, requestId jdi.EventRequestID) {
	var req = struct {
		Kind      jdi.EventKind
		RequestId jdi.EventRequestID
	}{kind, requestId}
	m.runCmd(connect.CmdEventRequestClear, &req, struct{}{})
}
func (m *MirrorImpl) eventRequestClearAllBreakpoints() {
	m.runCmd(connect.CmdEventRequestClearAllBreakpoints, struct{}{}, struct{}{})
}

// func (m *MirrorImpl) stackFrameSetValues(threadId jdi.ThreadID, frameId jdi.FrameID, stacks []jdi.StackFrameRequest) *[]jdi.ValueID {}
func (m *MirrorImpl) stackFrameGetValues(threadId jdi.ThreadID, frameId jdi.FrameID, variables []jdi.LocalVariable) map[jdi.LocalVariable]jdi.Value {
	stackRequest := make([]jdi.StackFrameRequest, len(variables))
	for index, value := range variables {
		stackRequest[index] = jdi.StackFrameRequest{Slot: jdi.Int(value.GetSlot()), SigByte: jdi.Tag(value.GetSignature()[0])}
	}
	var req = struct {
		ThreadId jdi.ThreadID
		FrameId  jdi.FrameID
		Stack    []jdi.StackFrameRequest
	}{threadId, frameId, stackRequest}
	var res []jdi.ValueID
	m.runCmd(connect.CmdStackFrameGetValues, &req, &res)
	result := m.readValueID(&res)
	out := make(map[jdi.LocalVariable]jdi.Value)
	for index, value := range *result {
		out[variables[index]] = value
	}
	return out
}
func (m *MirrorImpl) stackFrameThisObject(threadId jdi.ThreadID, frameId jdi.FrameID) jdi.ObjectReference {
	var req = struct {
		ThreadId jdi.ThreadID
		FrameId  jdi.FrameID
	}{threadId, frameId}
	var out jdi.TaggedObjectID
	m.runCmd(connect.CmdStackFrameThisObject, &req, &out)
	return m.makeObjectMirror(out.ObjectID, out.TagID)
}
func (m *MirrorImpl) stackFramePopFrames(threadId jdi.ThreadID, frameId jdi.FrameID) {
	var req = struct {
		ThreadId jdi.ThreadID
		FrameId  jdi.FrameID
	}{threadId, frameId}
	m.runCmd(connect.CmdStackFrameGetValues, &req, struct{}{})
}

func (m *MirrorImpl) classObjectReferenceReflectedType(classObjectID jdi.ClassObjectID) jdi.ReferenceType {
	var res struct {
		RefTypeTag jdi.TypeTag
		TypeID     jdi.ReferenceTypeID
	}
	m.runCmd(connect.CmdClassObjectReferenceReflectedType, classObjectID, &res)
	return m.makeReferenceTypeMirror(res.TypeID, res.RefTypeTag, &referenceTypeInfo{})
}
