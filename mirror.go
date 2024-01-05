package jdwp

type Mirror interface {
	// GetVirtualMachine 获取镜像引用的JVM对象引用
	GetVirtualMachine() VirtualMachine
}

type VirtualMachine interface {
	GetVirtualMachine() VirtualMachine
	GetClassesByName(className string) []ReferenceType
	GetClassesBySignature(string) []ReferenceType
	GetAllModules() ([]ModuleReference, error)
	GetAllClasses() []ReferenceType
	RedefineClasses(map[ReferenceType][]byte)
	GetAllThread() []ThreadReference
	Suspend()
	Resume()
	GetTopLevelThreadGroups() []ThreadGroupReference
	GetEventRequestManager() EventRequestManager
	MirrorOfBool(bool) BooleanValue
	MirrorOfString(string) StringReference
	MirrorOfByte(byte) ByteValue
	MirrorOfChar(char int16) CharValue
	MirrorOfInt(int) IntegerValue
	MirrorOfLong(int64) LongValue
	MirrorOfFloat(float32) FloatValue
	MirrorOfDouble(float64) DoubleValue
	MirrorOfVoid() VoidValue
	// Dispose 不实现!
	//Process()
	Dispose()
	Exit(int)
	CanWatchFieldModification() bool
	CanWatchFieldAccess() bool
	CanGetBytecodes() bool
	CanGetSyntheticAttribute() bool
	CanGetOwnedMonitorInfo() bool
	CanGetCurrentContendedMonitor() bool
	CanGetMonitorInfo() bool
	CanUseInstanceFilters() bool
	CanRedefineClasses() bool
	CanAddMethod() bool
	CanUnrestrictedlyRedefineClasses() bool
	CanPopFrames() bool
	CanGetSourceDebugExtension() bool
	CanRequestVMDeathEvent() bool
	CanGetMethodReturnValues() bool
	CanGetInstanceInfo() bool
	CanUseSourceNameFilters() bool
	CanForceEarlyReturn() bool
	CanBeModified() bool
	CanRequestMonitorEvents() bool
	CanGetMonitorFrameInfo() bool
	CanGetClassFileVersion() bool
	CanGetConstantPool() bool
	CanGetModuleInfo() bool
	GetDescription() string
	GetVersion() (string, error)
	GetName() string
}
