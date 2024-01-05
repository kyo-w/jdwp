package impl

import (
	"context"
	"fmt"
	jdi "github.com/kyo-w/jdwp"
	connect "github.com/kyo-w/jdwp/impl/internal"
	"net"
)

func Attach(ctx context.Context, hostname string) (jdi.VirtualMachine, error) {
	// netConn资源由VmConn处理释放
	netConn, err := net.Dial("tcp", hostname)
	if err != nil {
		return nil, err
	}
	vmConn, err := connect.Open(ctx, netConn)
	if err != nil {
		return nil, err
	}
	vm := &VirtualMachineImpl{conn: vmConn, Context: ctx}
	eventManager := &EventRequestManagerImpl{vm: vm}
	mirrorRoot := &MirrorImpl{vm: vm}
	vm.MirrorImpl = mirrorRoot
	vm.EventManager = eventManager
	return vm, nil
}

type VirtualMachineImpl struct {
	*MirrorImpl
	Context        context.Context
	conn           *connect.Connection
	EventManager   jdi.EventRequestManager
	version        *jdi.VmVersion
	theVoidType    *jdi.VoidType
	theByteType    *jdi.ByteType
	theBooleanType *jdi.BooleanType
	theCharType    *jdi.CharType
	theShortType   *jdi.ShortType
	theFloatType   *jdi.FloatType
	theDoubleType  *jdi.DoubleType
	theIntType     *jdi.IntegerType
	theLongType    *jdi.LongType
	capabilities   *jdi.Capabilities
	// ReferenceType缓存
	objectIdMap map[jdi.ObjectID]jdi.ObjectReference
	typeIdMap   map[jdi.ReferenceTypeID]jdi.ReferenceType
}

func (vm *VirtualMachineImpl) MirrorOfBool(b bool) jdi.BooleanValue {
	return &BooleanValueImpl{MirrorImpl: &MirrorImpl{vm: vm}, value: b}
}

func (vm *VirtualMachineImpl) MirrorOfString(s string) jdi.StringReference {
	var out jdi.StringID
	err := vm.conn.SendCommand(connect.CmdVirtualMachineCreateString, &s, &out)
	if err != nil {
		panic(err)
	}
	return &StringReferenceImpl{value: s, ObjectReferenceImpl: &ObjectReferenceImpl{ObjectId: jdi.ObjectID(out)}}
}

func (vm *VirtualMachineImpl) MirrorOfByte(b byte) jdi.ByteValue {
	return &ByteValueImpl{MirrorImpl: vm.createEmptyMirror(), value: b}
}

func (vm *VirtualMachineImpl) MirrorOfChar(char int16) jdi.CharValue {
	return &CharValueImpl{MirrorImpl: vm.createEmptyMirror(), value: jdi.Char(char)}
}

func (vm *VirtualMachineImpl) MirrorOfInt(i int) jdi.IntegerValue {
	return &IntegerValueImpl{MirrorImpl: vm.createEmptyMirror(), value: jdi.Int(i)}
}

func (vm *VirtualMachineImpl) MirrorOfLong(i int64) jdi.LongValue {
	return &LongValueImpl{MirrorImpl: vm.createEmptyMirror(), value: jdi.Long(i)}
}

func (vm *VirtualMachineImpl) MirrorOfFloat(f float32) jdi.FloatValue {
	return &FloatValueImpl{MirrorImpl: vm.createEmptyMirror(), value: jdi.Float(f)}
}

func (vm *VirtualMachineImpl) MirrorOfDouble(f float64) jdi.DoubleValue {
	return &DoubleValueImpl{MirrorImpl: vm.createEmptyMirror(), value: jdi.Double(f)}
}

func (vm *VirtualMachineImpl) MirrorOfVoid() jdi.VoidValue {
	return &VoidValueImpl{MirrorImpl: vm.createEmptyMirror()}
}

func (vm *VirtualMachineImpl) voidType() *jdi.VoidType {
	if vm.theVoidType == nil {
		vm.theVoidType = &jdi.VoidType{Vm: vm}
	}
	return vm.theVoidType
}

func (vm *VirtualMachineImpl) byteType() *jdi.ByteType {
	if vm.theByteType == nil {
		vm.theByteType = &jdi.ByteType{Vm: vm}
	}
	return vm.theByteType
}

func (vm *VirtualMachineImpl) booleanType() *jdi.BooleanType {
	if vm.theBooleanType == nil {
		vm.theBooleanType = &jdi.BooleanType{Vm: vm}
	}
	return vm.theBooleanType
}

func (vm *VirtualMachineImpl) charType() *jdi.CharType {
	if vm.theCharType == nil {
		vm.theCharType = &jdi.CharType{Vm: vm}
	}
	return vm.theCharType
}

func (vm *VirtualMachineImpl) shortType() *jdi.ShortType {
	if vm.theShortType == nil {
		vm.theShortType = &jdi.ShortType{Vm: vm}
	}
	return vm.theShortType
}

func (vm *VirtualMachineImpl) intType() *jdi.IntegerType {
	if vm.theIntType == nil {
		vm.theIntType = &jdi.IntegerType{Vm: vm}
	}
	return vm.theIntType
}

func (vm *VirtualMachineImpl) longType() *jdi.LongType {
	if vm.theLongType == nil {
		vm.theLongType = &jdi.LongType{Vm: vm}
	}
	return vm.theLongType
}

func (vm *VirtualMachineImpl) floatType() *jdi.FloatType {
	if vm.theFloatType == nil {
		vm.theFloatType = &jdi.FloatType{Vm: vm}
	}
	return vm.theFloatType
}

func (vm *VirtualMachineImpl) doubleType() *jdi.DoubleType {
	if vm.theDoubleType == nil {
		vm.theDoubleType = &jdi.DoubleType{Vm: vm}
	}
	return vm.theDoubleType
}

func (vm *VirtualMachineImpl) GetVersion() (string, error) {
	vm.initVersion()
	return vm.version.Version, nil
}

func (vm *VirtualMachineImpl) primitiveTypeMirror(tag jdi.Tag) jdi.Type {
	switch tag {
	case jdi.BOOLEAN:
		return vm.booleanType()
	case jdi.BYTE:
		return vm.byteType()
	case jdi.CHAR:
		return vm.charType()
	case jdi.SHORT:
		return vm.shortType()
	case jdi.INT:
		return vm.intType()
	case jdi.LONG:
		return vm.longType()
	case jdi.FLOAT:
		return vm.longType()
	case jdi.DOUBLE:
		return vm.doubleType()
	default:
		panic("unknown type : " + string(tag))
	}
}

func (vm *VirtualMachineImpl) GetClassesByName(className string) []jdi.ReferenceType {
	return vm.GetClassesBySignature(jdi.TranslateClassNameToSignature(className))
}
func (vm *VirtualMachineImpl) GetClassesBySignature(signature string) []jdi.ReferenceType {
	return *vm.vmClassesBySignature(signature)
}

// GetAllClasses 这里不能做Cache处理，GetAllClasses API在JDWP通信中是获取当前已经加载的Class，而不是系统全部的Class/**
func (vm *VirtualMachineImpl) GetAllClasses() []jdi.ReferenceType {
	return *vm.vmAllClasses()
}

func (vm *VirtualMachineImpl) RedefineClasses(m map[jdi.ReferenceType][]byte) {
	//TODO implement me
	panic("implement me")
}

func (vm *VirtualMachineImpl) GetAllThread() []jdi.ThreadReference {
	var out []jdi.ThreadID
	err := vm.conn.SendCommand(connect.CmdVirtualMachineAllThreads, struct{}{}, &out)
	if err != nil {
		return nil
	}
	return nil
}

func (vm *VirtualMachineImpl) Suspend() {
	err := vm.conn.SendCommand(connect.CmdVirtualMachineSuspend, struct{}{}, struct{}{})
	if err != nil {
		panic(err)
	}
}

func (vm *VirtualMachineImpl) Resume() {
	err := vm.conn.SendCommand(connect.CmdVirtualMachineResume, struct{}{}, struct{}{})
	if err != nil {
		panic(err)
	}
}

func (vm *VirtualMachineImpl) GetTopLevelThreadGroups() []jdi.ThreadGroupReference {
	return *vm.vmTopLevelThreadGroups()
}

func (vm *VirtualMachineImpl) GetEventRequestManager() jdi.EventRequestManager {
	return vm.EventManager
}
func (vm *VirtualMachineImpl) Dispose() {
	err := vm.conn.SendCommand(connect.CmdVirtualMachineDispose, struct{}{}, struct{}{})
	if err != nil {
		panic(err)
	}
}

func (vm *VirtualMachineImpl) Exit(exitCode int) {
	vm.vmExit(exitCode)
}

func (vm *VirtualMachineImpl) CanWatchFieldModification() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanWatchFieldModification
}

func (vm *VirtualMachineImpl) CanWatchFieldAccess() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanWatchFieldAccess
}

func (vm *VirtualMachineImpl) CanGetBytecodes() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanGetBytecodes
}

func (vm *VirtualMachineImpl) CanGetSyntheticAttribute() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanGetSyntheticAttribute
}

func (vm *VirtualMachineImpl) CanGetOwnedMonitorInfo() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanGetOwnedMonitorInfo
}

func (vm *VirtualMachineImpl) CanGetCurrentContendedMonitor() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanGetCurrentContendedMonitor
}

func (vm *VirtualMachineImpl) CanGetMonitorInfo() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanGetMonitorFrameInfo
}

func (vm *VirtualMachineImpl) CanUseInstanceFilters() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanUseInstanceFilters
}

func (vm *VirtualMachineImpl) CanRedefineClasses() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanRedefineClasses
}

func (vm *VirtualMachineImpl) CanAddMethod() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanAddMethod
}

func (vm *VirtualMachineImpl) CanUnrestrictedlyRedefineClasses() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanUnrestrictedlyRedefineClasses
}

func (vm *VirtualMachineImpl) CanPopFrames() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanPopFrames
}

func (vm *VirtualMachineImpl) CanGetSourceDebugExtension() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanGetSourceDebugExtension
}

func (vm *VirtualMachineImpl) CanRequestVMDeathEvent() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanRequestVMDeathEvent
}

func (vm *VirtualMachineImpl) CanGetMethodReturnValues() bool {
	vm.initVersion()
	return vm.version.JDWPMajor > 1 || vm.version.JDWPMinor >= 6
}

func (vm *VirtualMachineImpl) CanGetInstanceInfo() bool {
	vm.initVersion()
	if vm.version.JDWPMajor > 1 || vm.version.JDWPMinor >= 6 {
		vm.capabilitiesNew()
		return vm.capabilities.CanGetInstanceInfo
	}
	return false
}

func (vm *VirtualMachineImpl) CanUseSourceNameFilters() bool {
	vm.initVersion()
	return vm.version.JDWPMajor > 1 || vm.version.JDWPMinor >= 6
}

func (vm *VirtualMachineImpl) CanForceEarlyReturn() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanForceEarlyReturn
}

func (vm *VirtualMachineImpl) CanBeModified() bool {
	return true
}

func (vm *VirtualMachineImpl) CanRequestMonitorEvents() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanRequestMonitorEvents
}

func (vm *VirtualMachineImpl) CanGetMonitorFrameInfo() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanGetMonitorFrameInfo
}

func (vm *VirtualMachineImpl) CanGetClassFileVersion() bool {
	vm.initVersion()
	return vm.version.JDWPMajor > 1 || vm.version.JDWPMinor >= 6
}

func (vm *VirtualMachineImpl) CanGetConstantPool() bool {
	vm.capabilitiesNew()
	return vm.capabilities.CanGetConstantPool
}

func (vm *VirtualMachineImpl) CanGetModuleInfo() bool {
	vm.initVersion()
	return vm.version.JDWPMajor >= 9
}

func (vm *VirtualMachineImpl) GetDescription() string {
	vm.initVersion()
	return fmt.Sprintf("Java JVM %d %d: %s", vm.version.JDWPMajor, vm.version.JDWPMinor, vm.version.Description)
}

func (vm *VirtualMachineImpl) GetName() string {
	vm.initVersion()
	return vm.version.Name
}

func (vm *VirtualMachineImpl) capabilitiesNew() {
	if vm.capabilities == nil {
		vm.capabilities = vm.vmCapabilitiesNew()
	}
}

func (vm *VirtualMachineImpl) initVersion() {
	if vm.version == nil {
		vm.version = vm.vmGetVersion()
	}
}

func (vm *VirtualMachineImpl) GetVirtualMachine() jdi.VirtualMachine {
	return vm
}
func (vm *VirtualMachineImpl) GetAllModules() ([]jdi.ModuleReference, error) {
	return nil, nil
}
