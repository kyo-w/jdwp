package jdwp

type Type interface {
	Mirror
	// GetSignature 获取类型的签名
	GetSignature() string
	// GetTypeName 获得类型的名称
	GetTypeName() string
}

type BooleanType struct {
	Vm VirtualMachine
}
type ByteType struct {
	Vm VirtualMachine
}

type CharType struct {
	Vm VirtualMachine
}
type DoubleType struct {
	Vm VirtualMachine
}
type FloatType struct {
	Vm VirtualMachine
}
type IntegerType struct {
	Vm VirtualMachine
}
type LongType struct {
	Vm VirtualMachine
}
type ShortType struct {
	Vm VirtualMachine
}
type VoidType struct {
	Vm VirtualMachine
}

func (b *BooleanType) GetVirtualMachine() VirtualMachine {
	return b.Vm
}

func (b *ByteType) GetVirtualMachine() VirtualMachine {
	return b.Vm
}

func (c *CharType) GetVirtualMachine() VirtualMachine {
	return c.Vm
}

func (d *DoubleType) GetVirtualMachine() VirtualMachine {
	return d.Vm
}

func (f *FloatType) GetVirtualMachine() VirtualMachine {
	return f.Vm
}

func (i *IntegerType) GetVirtualMachine() VirtualMachine {
	return i.Vm
}

func (l *LongType) GetVirtualMachine() VirtualMachine {
	return l.Vm
}

func (s *ShortType) GetVirtualMachine() VirtualMachine {
	return s.Vm
}

func (v *VoidType) GetVirtualMachine() VirtualMachine {
	return v.Vm
}

func (b *BooleanType) GetTypeName() string {
	return "boolean"
}

func (b *ByteType) GetTypeName() string {
	return "byte"
}

func (c *CharType) GetTypeName() string {
	return "char"
}

func (d *DoubleType) GetTypeName() string {
	return "double"
}

func (f *FloatType) GetTypeName() string {
	return "float"
}

func (i *IntegerType) GetTypeName() string {
	return "int"
}

func (l *LongType) GetTypeName() string {
	return "long"
}

func (s *ShortType) GetTypeName() string {
	return "short"
}

func (v *VoidType) GetTypeName() string {
	return "void"
}

func (v *VoidType) GetSignature() string {
	return string(rune(TagVoid))
}

func (s *ShortType) GetSignature() string {
	return string(rune(TagShort))
}

func (d *DoubleType) GetSignature() string {
	return string(rune(TagDouble))
}
func (f *FloatType) GetSignature() string {
	return string(rune(TagFloat))
}
func (i *IntegerType) GetSignature() string {
	return string(rune(TagInt))
}
func (l *LongType) GetSignature() string {
	return string(rune(TagLong))
}
func (c *CharType) GetSignature() string {
	return string(rune(TagChar))
}
func (b *BooleanType) GetSignature() string {
	return string(rune(TagBoolean))
}
func (b *ByteType) GetSignature() string {
	return string(rune(TagByte))
}

type ReferenceType interface {
	Mirror
	Type // GetGenericSignature 与GetSignature不同, GetGenericSignature适用于泛型引用对象
	GetGenericSignature() string
	// GetClassLoader 获取引用对象的ClassLoader引用
	GetClassLoader() ClassLoaderReference
	// GetModule 获取包含与此类型对应的类的模块对象。并非所有目标虚拟机都支持此操作。使用VirtualMachine.canGetModuleInfo（）确定是否支持该操作。
	GetModule() ModuleReference
	// IsStatic 引用对象是否为静态类
	IsStatic() bool
	// IsAbstract 引用对象是否为抽象类
	IsAbstract() bool
	// IsFinal 引用对象是否存在final关键字标识符
	IsFinal() bool
	// IsPrepared 引用对象的类是否已经加载(JVM加载字节码的一个状态)
	IsPrepared() bool
	// IsVerified 引用对象的类是否经过验证的流程(JVM加载字节码的一个状态)
	IsVerified() bool
	// IsInitialized 引用对象的类是否进行初始化(JVM加载字节码的一个状态)
	IsInitialized() bool
	// FailedToInitialize 引用对象的类初始化失败
	FailedToInitialize() bool
	// GetFields 获取所有的字段。继承字段不会计入!
	GetFields() []Field
	// GetAllVisibleFields 获取可见的字段(不包含private)。继承字段不会计入!
	GetAllVisibleFields() []Field
	// GetAllFields 获取所有的字段。包含继承/实现字段
	GetAllFields() []Field
	// GetFieldByName 根据字段名获取字段引用
	GetFieldByName(string) Field
	// GetMethods 获取所有的方法,继承方法不会计入!
	GetMethods() []Method
	// GetVisibleMethods 获取所有的方法(不包含private)，继承方法不会计入!
	GetVisibleMethods() []Method
	// GetAllMethods 获取所有的方法。包含继承/实现方法
	GetAllMethods() []Method
	// GetMethodsByName 根据方法名获取method引用
	GetMethodsByName(string) []Method
	// GetMethodsByNameAndSign 根据方法名和方法签名获取method引用
	GetMethodsByNameAndSign(string, string) []Method
	// GetNestedTypes 返回此ReferenceType加载时所需要的其他ReferenceType
	GetNestedTypes() []ReferenceType
	// GetValue 根据字段值(Field)返回字段的属性值，类似对象的反射获取字段值
	GetValue(Field) Value
	// GetValues GetValue只能对一个字段操作，GetValues提供多个字段获取，返回值是Map[Field]Value
	GetValues([]Field) map[Field]Value
	// GetClassObject 获取ReferenceType的Class对象引用。ClassObjectReference是Class类的引用
	GetClassObject() ClassObjectReference

	// GetAllLineLocations 返回类的所有行号
	GetAllLineLocations() []Location
	GetAllInterfaces() []InterfaceType

	// GetLocationsOfLine 通过行号获取Location
	GetLocationsOfLine(int) []Location
	// GetInstances 返回内存中所有ReferenceType的对象引用，注意：参数代表最多接受多少。比如max = 3，内存中有5个，此时返回值最多返回3个。max = 0， 不做任何限制
	GetInstances(max int64) []ObjectReference
	// GetMinorVersion 字节码版本
	GetMinorVersion() int
	GetUniqueID() ReferenceTypeID
	GetTypeTag() TypeTag
}

type ArrayType interface {
	ReferenceType
	NewInstance(length int) ArrayReference
	GetComponentSignature() string
	GetComponentTypeName() string
	GetComponentType() Type
}

type ClassType interface {
	ReferenceType
	GetSuperclass() ClassType
	GetOwnInterface() []InterfaceType
	GetSubclasses() []ClassType
	IsEnum() bool
	InvokeMethod(reference ThreadReference, method Method, args []Value, options InvokeOptions) (valueRef Value, error ObjectReference)
}

type InterfaceType interface {
	ReferenceType
	GetSuperInterfaces() []InterfaceType
	GetSubInterfaces() []InterfaceType
	GetAllImplementors() []ClassType
	InvokeMethod(thread ThreadReference, method Method, args []Value, options InvokeOptions) (value ObjectReference, error ObjectReference)
}
