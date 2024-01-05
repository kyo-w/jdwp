package jdwp

type StackFrameRequest struct {
	Slot    Int
	SigByte Tag
}

type StackFrame interface {
	Mirror
	// GetLocation 堆栈其实就是方法的堆栈调用返回当前堆栈方法的Location
	GetLocation() Location
	// GetThread 获取栈帧的线程引用
	GetThread() ThreadReference
	// GetThisObject 获取当前栈帧的this对象
	GetThisObject() ObjectReference

	// GetVisibleVariables Throws:
	//AbsentInformationException – if there is no local variable information for this method.
	//InvalidStackFrameException – if this stack frame has become invalid. Once the frame's thread is resumed, the stack frame is no longer valid.
	//NativeMethodException – if the current method is native.
	GetVisibleVariables() []LocalVariable

	// GetVisibleVariableByName Throws: AbsentInformationException – if there is no local variable information for this method.
	//InvalidStackFrameException – if this stack frame has become invalid. Once the frame's thread is resumed, the stack frame is no longer valid.
	//NativeMethodException – if the current method is native.
	GetVisibleVariableByName(name string) LocalVariable

	// GetValue 堆栈中的LocalVariable取出它的对象引用，这个LocalVariable必须存在与当前的堆栈之中
	GetValue(LocalVariable) Value

	GetValues([]LocalVariable) map[LocalVariable]Value
	// GetArgumentValues 返回此帧中所有参数的值。即使不存在局部变量信息，也会返回值。
	GetArgumentValues() []Value
}
