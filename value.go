package jdwp

/*
*
下面的所有接口均是对象引用的所有类型
*/
const (
	INVOKE_SINGLE_THREADED = 0x1
	INVOKE_NONVIRTUAL      = 0x2

	THREAD_STATUS_UNKNOWN     = -1
	THREAD_STATUS_ZOMBIE      = 0
	THREAD_STATUS_RUNNING     = 1
	THREAD_STATUS_SLEEPING    = 2
	THREAD_STATUS_MONITOR     = 3
	THREAD_STATUS_WAIT        = 4
	THREAD_STATUS_NOT_STARTED = 5
)

type Int int
type Float float32
type Double float64
type Char int16
type Long int64
type Short int16

type Value interface {
	Mirror
	GetType() Type
	GetTagType() Tag
}

type BooleanValue interface {
	Value
	GetValue() bool
}
type ByteValue interface {
	Value
	GetValue() byte
}
type CharValue interface {
	Value
	GetValue() Char
}
type DoubleValue interface {
	Value
	GetValue() Double
}
type FloatValue interface {
	Value
	GetValue() Float
}
type IntegerValue interface {
	Value
	GetValue() Int
}
type LongValue interface {
	Value
	GetValue() Long
}
type ShortValue interface {
	Value
	GetValue() Short
}
type StringReference interface {
	ObjectReference
	GetStringValue() string
}

type ObjectReference interface {
	Value
	GetUniqueID() ObjectID
	// GetReferenceType 返回对象引用的ReferenceType
	GetReferenceType() ReferenceType
	GetValuesByFieldNames(fields ...string) Value
	// GetValueByField 通过字段返回Value值
	GetValueByField(Field) Value
	GetValuesByFields([]Field) map[Field]Value

	// InvokeMethod
	//ThreadReference 线程引用
	//Method 方法引用
	//[]Value 参数列表
	//int 参数
	///**
	InvokeMethod(reference ThreadReference, method Method, args []Value, options InvokeOptions) (value Value, error ObjectReference)

	// DisableCollection 停止垃圾回收机制
	DisableCollection()
	// EnableCollection 开启垃圾回收机制
	EnableCollection()
	// IsCollected 是否被垃圾回收
	IsCollected() bool

	// GetReferringObjects /**
	GetReferringObjects(maxReferrers int) []ObjectReference
}

type ArrayReference interface {
	ObjectReference
	// GetLength 返回数组的长度
	GetLength() int
	GetArrayValue(index int) Value
	GetArrayValues() []Value
	// GetArraySlice
	//index表示起始位置
	//length表示要获取的个数
	///**
	GetArraySlice(index, length int) []Value
}

type ClassLoaderReference interface {
	ObjectReference

	// GetDefinedClasses
	//返回当前ClassLoaderReference引用加载的所有字节码ReferenceType引用
	///**
	GetDefinedClasses() []ReferenceType

	// GetVisibleClasses
	//返回当前ClassLoaderReference引用加载并且初始化的所有字节码ReferenceType引用
	///**
	GetVisibleClasses() []ReferenceType
}

type ClassObjectReference interface {
	ObjectReference
	// GetReflectedType 返回Class类对应的类ReferenceType
	GetReflectedType() ReferenceType
}

type ModuleReference interface {
	ObjectReference
	GetName() string
	GetClassLoader() ClassLoaderReference
}

type ThreadGroupReference interface {
	ObjectReference
	GetName() string
	// GetParent 返回此线程组的父级。
	GetParent() ThreadGroupReference
	// Suspend 挂起线程组所有的线程
	Suspend()
	// Resume 将所有的挂起线程启动起来
	Resume()
	// GetAllThread 返回线程组所有的线程引用
	GetAllThread() []ThreadReference
	// GetThreadGroups 返回一个列表，其中包含此线程组中的每个活动ThreadGroupReference。仅返回此直接线程组中的活动线程组（而不是其子组）。有关“活动”线程组的信息，请参阅线程组
	GetThreadGroups() []ThreadGroupReference
}
type ThreadReference interface {
	ObjectReference
	GetName() string
	Suspend()
	Resume()
	SuspendCount() int
	Status() ThreadStatus
	IsSuspended() bool
	GetThreadGroup() ThreadGroupReference
	GetFrameCount() int
	GetFrames() []StackFrame
	GetFrameByIndex(int) StackFrame
	GetFrameSlice(start, length int) []StackFrame
}
type VoidValue interface {
	Value
	Void()
}
