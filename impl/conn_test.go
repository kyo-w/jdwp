package impl

import (
	"context"
	"github.com/kyo-w/jdwp"
	"testing"
)

func TestVm(t *testing.T) {
	attach, err := Attach(context.Background(), "192.168.116.100:5005")
	if err != nil {
		panic(err)
	}

	name := attach.GetClassesByName("sun.misc.Launcher$AppClassLoader")
	instances := name[0].GetInstances(0)
	names := instances[0].GetValuesByFieldNames("ucp", "path")
	data := names.(jdwp.ObjectReference).GetValuesByFieldNames("elementData")
	values := data.(jdwp.ArrayReference).GetArrayValues()
}
