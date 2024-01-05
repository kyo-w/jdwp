package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type StackFrameImpl struct {
	*MirrorImpl
	StackFrameId     jdi.FrameID
	ThreadRef        jdi.ThreadReference
	Location         jdi.Location
	thisObject       jdi.ObjectReference
	visibleVariables []jdi.LocalVariable
}

func (s *StackFrameImpl) GetArgumentValues() []jdi.Value {
	panic("implement me")
}

func (s *StackFrameImpl) GetLocation() jdi.Location {
	return s.Location
}

func (s *StackFrameImpl) GetThread() jdi.ThreadReference {
	return s.ThreadRef
}

func (s *StackFrameImpl) GetThisObject() jdi.ObjectReference {
	if s.thisObject == nil {
		s.thisObject = s.stackFrameThisObject(jdi.ThreadID(s.ThreadRef.GetUniqueID()), s.StackFrameId)
	}
	return s.thisObject
}

func (s *StackFrameImpl) GetVisibleVariables() []jdi.LocalVariable {
	if s.visibleVariables == nil {
		var visibleVar []jdi.LocalVariable
		variables := s.Location.GetMethod().GetVariables()
		for _, varValue := range variables {
			if varValue.IsVisible(s) {
				visibleVar = append(visibleVar, varValue)
			}
		}
		s.visibleVariables = visibleVar
	}
	return s.visibleVariables
}

func (s *StackFrameImpl) GetVisibleVariableByName(name string) jdi.LocalVariable {
	for _, varValue := range s.GetVisibleVariables() {
		if varValue.GetName() == name {
			return varValue
		}
	}
	return nil
}

func (s *StackFrameImpl) GetValue(variable jdi.LocalVariable) jdi.Value {
	return s.GetValues([]jdi.LocalVariable{variable})[variable]
}

func (s *StackFrameImpl) GetValues(variables []jdi.LocalVariable) map[jdi.LocalVariable]jdi.Value {
	for _, varValue := range variables {
		if !varValue.IsVisible(s) {
			panic(varValue.GetName() + " is not valid at this frame location")
		}
	}
	return s.stackFrameGetValues(jdi.ThreadID(s.ThreadRef.GetUniqueID()), s.StackFrameId, variables)
}
