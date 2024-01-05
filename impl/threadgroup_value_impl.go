package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type ThreadGroupReferenceImpl struct {
	*ObjectReferenceImpl
	name           string
	parent         jdi.ThreadGroupReference
	initParent     bool
	subThread      []jdi.ThreadReference
	subThreadGroup []jdi.ThreadGroupReference
}

func (t *ThreadGroupReferenceImpl) GetName() string {
	return t.threadGroupReferenceName(jdi.ThreadGroupID(t.ObjectId))
}
func (t *ThreadGroupReferenceImpl) GetParent() jdi.ThreadGroupReference {
	if !t.initParent {
		t.parent = t.threadGroupReferenceParent(jdi.ThreadGroupID(t.ObjectId))
	}
	return t.parent
}
func (t *ThreadGroupReferenceImpl) Suspend() {
	for _, value := range t.GetAllThread() {
		value.Suspend()
	}
	for _, value := range t.GetThreadGroups() {
		value.Suspend()
	}
}
func (t *ThreadGroupReferenceImpl) Resume() {
	for _, value := range t.GetAllThread() {
		value.Resume()
	}
	for _, value := range t.GetThreadGroups() {
		value.Resume()
	}
}
func (t *ThreadGroupReferenceImpl) GetAllThread() []jdi.ThreadReference {
	var out []jdi.ThreadReference
	if t.subThread == nil && t.subThreadGroup == nil {
		t.subThread, t.subThreadGroup = t.threadGroupReferenceChildren(jdi.ThreadGroupID(t.ObjectId))
		out = append(out, t.subThread...)
	}
	if t.subThreadGroup != nil {
		for _, value := range t.subThreadGroup {
			out = append(out, value.GetAllThread()...)
		}
	}
	return out
}
func (t *ThreadGroupReferenceImpl) GetThreadGroups() []jdi.ThreadGroupReference {
	if t.subThread == nil && t.subThreadGroup == nil {
		t.subThread, t.subThreadGroup = t.threadGroupReferenceChildren(jdi.ThreadGroupID(t.ObjectId))
	}
	return t.subThreadGroup
}
