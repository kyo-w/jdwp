package impl

import (
	"context"
	jdi "github.com/kyo-w/jdwp"
	"log"
	"testing"
)

func TestVm(t *testing.T) {
	vm, _ := Attach(context.Background(), ":5005")
	signature := vm.GetClassesByName("org.apache.catalina.mapper.Mapper")[0]
	for _, value := range signature.GetAllMethods() {
		if value.GetName() == "internalMap" {
			location := value.GetLocation()
			request := vm.GetEventRequestManager().CreateBreakpointRequest(location)
			request.SetSuspendPolicy(jdi.SuspendEventThread)
			var i = 0
			request.SetHandler(func(eventObject jdi.EventObject) bool {
				object := eventObject.(EventBreakpointResponseObject)
				thread := object.GetThread()
				index := thread.GetFrameByIndex(0)
				thisObject := index.GetThisObject()
				getType := thisObject.GetType()
				log.Println(getType.GetTypeName())
				if i == 1 {
					i = i + 1
					return false
				}
				return true
			})
			request.Enable()
		}
	}
	select {}
}
