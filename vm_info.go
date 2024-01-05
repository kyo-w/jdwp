package jdwp

type VmVersion struct {
	Description string //		JVM描述
	JDWPMajor   int    //		Major版本号
	JDWPMinor   int    //		Minor版本号
	Version     string //		Target VM 的Java版本号
	Name        string //		Target VM name
}
type Capabilities struct {
	CanWatchFieldModification        bool
	CanWatchFieldAccess              bool
	CanGetBytecodes                  bool
	CanGetSyntheticAttribute         bool
	CanGetOwnedMonitorInfo           bool
	CanGetCurrentContendedMonitor    bool
	CanGetMonitorInfo                bool
	CanRedefineClasses               bool
	CanAddMethod                     bool
	CanUnrestrictedlyRedefineClasses bool
	CanPopFrames                     bool
	CanUseInstanceFilters            bool
	CanGetSourceDebugExtension       bool
	CanRequestVMDeathEvent           bool
	CanSetDefaultStratum             bool
	CanGetInstanceInfo               bool
	CanRequestMonitorEvents          bool
	CanGetMonitorFrameInfo           bool
	CanUseSourceNameFilters          bool
	CanGetConstantPool               bool
	CanForceEarlyReturn              bool
	R22                              bool
	R23                              bool
	R24                              bool
	R25                              bool
	R26                              bool
	R27                              bool
	R28                              bool
	R29                              bool
	R30                              bool
	R31                              bool
	R32                              bool
}

type ClassPath struct {
	BaseDir        string
	ClassPaths     []string
	BootClassPaths []string
}
