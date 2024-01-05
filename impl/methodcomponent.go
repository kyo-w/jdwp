package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type MethodImpl struct {
	*TypeComponentImpl
	initLocation bool
	initVar      bool
	startIndex   int
	endIndex     int
	argCount     int
	vars         []jdi.LocalVariable
	locations    []jdi.Location
}

func (m *MethodImpl) GetLocation() jdi.Location {
	if !m.initLocation {
		m.initLocations()
	}
	return m.locations[0]
}

func (m *MethodImpl) GetReturnTypeName() string {
	return jdi.TranslateSignatureToClassName(m.signature)
}

func (m *MethodImpl) GetReturnType() jdi.Type {
	return findType(m.GetDeclaringType(), m.GetSignature())
}

func (m *MethodImpl) GetArgumentTypeNames() []string {
	if !m.initVar {
		m.initVars()
	}
	out := make([]string, m.argCount)
	for index, varValue := range m.vars[:m.argCount-1] {
		out[index] = varValue.GetName()
	}
	return out
}

func (m *MethodImpl) GetArgumentTypes() []jdi.Type {
	if !m.initVar {
		m.initVars()
	}
	out := make([]jdi.Type, m.argCount)
	for index, varValue := range m.vars[:m.argCount-1] {
		out[index] = varValue.GetType()
	}
	return out
}

func (m *MethodImpl) IsAbstract() bool {
	return m.isModifierSet(ABSTRACT)
}

func (m *MethodImpl) IsDefault() bool {
	_, isInterface := m.GetDeclaringType().(jdi.InterfaceType)
	return !m.isModifierSet(ABSTRACT) &&
		!m.isModifierSet(STATIC) &&
		!m.isModifierSet(PRIVATE) && isInterface
}

func (m *MethodImpl) IsSynchronized() bool {
	return m.isModifierSet(SYNCHRONIZED)
}

func (m *MethodImpl) IsNative() bool {
	return m.isModifierSet(NATIVE)
}

func (m *MethodImpl) IsVarArgs() bool {
	return m.isModifierSet(VARARGS)
}

func (m *MethodImpl) IsBridge() bool {
	return m.isModifierSet(BRIDGE)
}

func (m *MethodImpl) IsConstructor() bool {
	return m.Name == "<init>"
}

func (m *MethodImpl) IsStaticInitializer() bool {
	return m.Name == "<clinit>"
}

func (m *MethodImpl) IsObsolete() bool {
	return m.methodTypeIsObsolete(m.GetDeclaringType().GetUniqueID(), jdi.MethodID(m.GetUniqueID()))
}

func (m *MethodImpl) GetAllLineLocation() []jdi.Location {
	if !m.initLocation {
		m.initLocations()
	}
	return m.locations
}

func (m *MethodImpl) GetLocationsOfLine(lineNumber int) []jdi.Location {
	if !m.initLocation {
		m.initLocations()
	}
	var out []jdi.Location
	for _, value := range m.locations {
		if lineNumber == value.GetLineNumber() {
			out = append(out, value)
		}
	}
	return out
}

func (m *MethodImpl) GetLocationOfCodeIndex(codeIndex int64) jdi.Location {
	if !m.initLocation {
		m.initLocations()
	}
	for _, value := range m.locations {
		if codeIndex == value.GetCodeIndex() {
			return value
		}
	}
	return nil
}

func (m *MethodImpl) GetVariablesByName(name string) []jdi.LocalVariable {
	if !m.initVar {
		m.initVars()
	}
	var out []jdi.LocalVariable
	for _, value := range m.vars {
		if name == value.GetName() {
			out = append(out, value)
		}
	}
	return out
}

func (m *MethodImpl) GetArguments() []jdi.LocalVariable {
	if !m.initVar {
		m.initVars()
	}
	return m.vars[0:m.argCount]
}

func (m *MethodImpl) GetVariables() []jdi.LocalVariable {
	if !m.initVar {
		m.initVars()
	}
	return m.vars
}

func (m *MethodImpl) initVars() {
	argCount, vars := m.methodTypeVariableTable(m.GetDeclaringType().GetUniqueID(), m)
	m.argCount = int(argCount)
	m.vars = vars
	m.initVar = true
}
func (m *MethodImpl) initLocations() {
	StartIndex, EndIndex, locations := m.methodTypeLineTable(m.GetDeclaringType(), m)
	m.startIndex = int(StartIndex)
	m.endIndex = int(EndIndex)
	m.locations = locations
	m.initLocation = true
}
