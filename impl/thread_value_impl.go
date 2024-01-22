package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type ThreadReferenceImpl struct {
	*ObjectReferenceImpl
	name                 string
	ThreadGroup          jdi.ThreadGroupReference
	suspendedZombieCount int
	frameCount           int
}

//func (t *ThreadReferenceImpl) IsAtBreakpoint() bool {
//	status := t.threadReferenceStatus(jdi.ThreadID(t.ObjectId))
//}

func (t *ThreadReferenceImpl) GetName() string {
	if t.name == "" {
		t.name = t.threadReferenceName(jdi.ThreadID(t.ObjectId))
	}
	return t.name
}

func (t *ThreadReferenceImpl) Suspend() {
	t.threadReferenceSuspend(jdi.ThreadID(t.ObjectId))
}

func (t *ThreadReferenceImpl) Resume() {
	if t.suspendedZombieCount > 0 {
		t.suspendedZombieCount--
		return
	}
	t.threadReferenceResume(jdi.ThreadID(t.ObjectId))
}

func (t *ThreadReferenceImpl) SuspendCount() int {
	if t.suspendedZombieCount > 0 {
		return t.suspendedZombieCount
	}
	return t.threadReferenceSuspendCount(jdi.ThreadID(t.ObjectId))
}

func (t *ThreadReferenceImpl) Status() jdi.ThreadStatus {
	return t.threadReferenceStatus(jdi.ThreadID(t.ObjectId))
}

func (t *ThreadReferenceImpl) IsSuspended() bool {
	status := t.threadReferenceStatus(jdi.ThreadID(t.ObjectId))
	return (t.suspendedZombieCount > 0) || (status.SuspendStatus&0x1) != 0
}

func (t *ThreadReferenceImpl) GetThreadGroup() jdi.ThreadGroupReference {
	if t.ThreadGroup == nil {
		t.ThreadGroup = t.threadReferenceThreadGroup(jdi.ThreadID(t.ObjectId))
	}
	return t.ThreadGroup
}

func (t *ThreadReferenceImpl) GetFrameCount() int {
	if t.frameCount == 0 {
		t.frameCount = t.threadReferenceFrameCount(jdi.ThreadID(t.ObjectId))
	}
	return t.frameCount
}

func (t *ThreadReferenceImpl) GetFrames() []jdi.StackFrame {
	return t.threadReferenceFrames(t, 0, -1)
}

func (t *ThreadReferenceImpl) GetFrameByIndex(i int) jdi.StackFrame {
	frames := t.threadReferenceFrames(t, i, 1)
	return frames[i]

}

func (t *ThreadReferenceImpl) GetFrameSlice(start, length int) []jdi.StackFrame {
	return t.threadReferenceFrames(t, start, length)
}
func (t *ThreadReferenceImpl) GetTagType() jdi.Tag {
	return jdi.THREAD
}
