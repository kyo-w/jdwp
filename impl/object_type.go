package impl

import (
	jdi "github.com/kyo-w/jdwp"
	"strings"
)

const (
	MAX_VALUE = 0x7fffffff
)

type ReferenceTypeImpl struct {
	*MirrorImpl
	Kind             jdi.TypeTag
	TypeID           jdi.ReferenceTypeID
	status           jdi.ClassStatus
	signatureName    string
	genericSignature string
	classLoader      jdi.ClassLoaderReference
	modifiers        int
	hasGetModifiers  bool
	fieldsRef        []jdi.Field
	module           jdi.ModuleReference
	hasModule        bool
	hasStatus        bool
	initSuperClass   bool
	superClass       jdi.ClassType
	methods          []jdi.Method
}

func (r *ReferenceTypeImpl) GetTypeTag() jdi.TypeTag {
	return r.Kind
}

func (r *ReferenceTypeImpl) GetAllInterfaces() []jdi.InterfaceType {
	interfaces := r.referenceTypeInterfaces(r.GetUniqueID())
	for _, value := range *interfaces {
		*interfaces = append(*interfaces, value.GetAllInterfaces()...)
	}
	return *interfaces
}
func (r *ReferenceTypeImpl) GetUniqueID() jdi.ReferenceTypeID {
	return r.TypeID
}

func (r *ReferenceTypeImpl) GetTypeName() string {
	return jdi.TranslateSignatureToClassName(r.GetSignature())
}

func (r *ReferenceTypeImpl) GetSignature() string {
	if r.signatureName == "" {
		r.signatureName = r.referenceTypeSignature(r.TypeID)
	}
	return r.signatureName
}

func (r *ReferenceTypeImpl) GetGenericSignature() string {
	if r.genericSignature == "" {
		r.signatureName, r.genericSignature = r.referenceSignatureWithGeneric(r.TypeID)
	}
	return r.genericSignature
}

func (r *ReferenceTypeImpl) GetClassLoader() jdi.ClassLoaderReference {
	if r.classLoader == nil {
		r.classLoader = r.referenceTypeClassLoader(r.GetUniqueID())
	}
	return r.classLoader
}

// GetModule commandSet 2 command 19 文档找不到
func (r *ReferenceTypeImpl) GetModule() jdi.ModuleReference {
	return nil
}

func (r *ReferenceTypeImpl) getModifiers() {
	if !r.hasGetModifiers {
		r.modifiers = r.referenceTypeModifiers(r.TypeID)
	}
}

func (r *ReferenceTypeImpl) IsStatic() bool {
	if r.modifiers == 0 {
		r.getModifiers()
	}
	return (r.modifiers & PUBLIC) > 0
}

func (r *ReferenceTypeImpl) IsAbstract() bool {
	if r.modifiers == 0 {
		r.getModifiers()
	}
	return (r.modifiers & ABSTRACT) > 0
}

func (r *ReferenceTypeImpl) IsFinal() bool {
	if r.modifiers == 0 {
		r.getModifiers()
	}
	return (r.modifiers & FINAL) > 0
}

func (r *ReferenceTypeImpl) IsPrepared() bool {
	if !r.hasStatus {
		r.status = r.referenceTypeStatus(r.TypeID)
		r.hasStatus = true
	}
	return (r.status & jdi.StatusPrepared) != 0
}

func (r *ReferenceTypeImpl) IsVerified() bool {
	if r.status&jdi.StatusVerified == 0 {
		r.status = r.referenceTypeStatus(r.TypeID)
		r.hasStatus = true
	}
	return (r.status & jdi.StatusVerified) != 0
}

func (r *ReferenceTypeImpl) IsInitialized() bool {
	if (r.status & (jdi.StatusInitialized | jdi.StatusError)) == 0 {
		r.status = r.referenceTypeStatus(r.TypeID)
		r.hasStatus = true
	}
	return (r.status & jdi.StatusInitialized) != 0
}

func (r *ReferenceTypeImpl) FailedToInitialize() bool {
	if (r.status & (jdi.StatusInitialized | jdi.StatusError)) == 0 {
		r.status = r.referenceTypeStatus(r.TypeID)
		r.hasStatus = true
	}
	return (r.status & jdi.StatusError) != 0
}

func (r *ReferenceTypeImpl) GetFields() []jdi.Field {
	return *r.referenceTypeFields(r, r.TypeID)
}

func (r *ReferenceTypeImpl) GetAllVisibleFields() []jdi.Field {
	types := r.inheritedTypes()
	var out []jdi.Field
	for _, typeValue := range types {
		out = append(out, typeValue.GetFields()...)
	}
	return out
}

func (r *ReferenceTypeImpl) GetAllFields() []jdi.Field {
	return r.GetAllVisibleFields()
}
func (r *ReferenceTypeImpl) inheritedTypes() []jdi.ReferenceType {
	var out = []jdi.ReferenceType{r}
	switch r.Kind {
	case jdi.ArrayTypeTag:
		return out
	case jdi.ClassTypeTag:
		if r.getSuperClass() != nil {
			out = append(out, r.getSuperClass())
		}
		for _, value := range r.GetAllInterfaces() {
			out = append(out, value)
		}
	}
	return out
}

func (r *ReferenceTypeImpl) GetFieldByName(name string) jdi.Field {
	for _, value := range r.GetAllFields() {
		if value.GetName() == name {
			return value
		}
	}
	return nil
}

func (r *ReferenceTypeImpl) GetMethods() []jdi.Method {
	return *r.referenceMethodsWithGeneric(r, r.TypeID)
}

func (r *ReferenceTypeImpl) GetVisibleMethods() []jdi.Method {
	types := r.inheritedTypes()
	var out []jdi.Method
	for _, value := range types {
		out = append(out, value.GetMethods()...)
	}
	return out
}

func (r *ReferenceTypeImpl) GetAllMethods() []jdi.Method {
	return r.GetVisibleMethods()
}

func (r *ReferenceTypeImpl) GetMethodsByName(name string) []jdi.Method {
	var out []jdi.Method
	for _, value := range r.GetVisibleMethods() {
		if value.GetName() == name {
			out = append(out, value)
		}
	}
	return out
}

func (r *ReferenceTypeImpl) GetMethodsByNameAndSign(name string, signature string) []jdi.Method {
	var out []jdi.Method
	for _, value := range r.GetVisibleMethods() {
		if value.GetName() == name && value.GetSignature() == signature {
			out = append(out, value)
		}
	}
	return out
}

func (r *ReferenceTypeImpl) GetNestedTypes() []jdi.ReferenceType {
	return *r.referenceTypeNestedTypes(r.TypeID)
}

func (r *ReferenceTypeImpl) GetValue(field jdi.Field) jdi.Value {
	return r.GetValues([]jdi.Field{field})[field]
}

func (r *ReferenceTypeImpl) GetValues(fields []jdi.Field) map[jdi.Field]jdi.Value {
	fieldIdList := make([]jdi.FieldID, len(fields))
	for index, value := range fields {
		if !value.IsStatic() {
			panic("referenceType just can get static field, if not static, please use objectReference")
		}
		fieldIdList[index] = jdi.FieldID(value.GetUniqueID())
	}
	fieldValues := r.referenceTypeGetValues(r.TypeID, fieldIdList)
	out := make(map[jdi.Field]jdi.Value)
	for index, value := range *fieldValues {
		out[fields[index]] = value
	}
	return out
}

func (r *ReferenceTypeImpl) GetClassObject() jdi.ClassObjectReference {
	return (*r.referenceClassObject(r.TypeID))[0]
}

func (r *ReferenceTypeImpl) GetAllLineLocations() []jdi.Location {
	methods := r.GetMethods()
	var out []jdi.Location
	for _, value := range methods {
		out = append(out, value.GetAllLineLocation()...)
	}
	return out
}

func (r *ReferenceTypeImpl) GetLocationsOfLine(lineIndex int) []jdi.Location {
	var out []jdi.Location
	methods := r.GetMethods()
	for _, value := range methods {
		for _, location := range value.GetAllLineLocation() {
			if location.GetLineNumber() == lineIndex {
				out = append(out, location)
			}
		}
	}
	return out
}

func (r *ReferenceTypeImpl) GetInstances(max int64) []jdi.ObjectReference {
	if !r.vm.CanGetInstanceInfo() {
		panic("target does not support getting instances")
	}
	if max < 0 {
		panic("maxInstances is less than zero")
	}
	var intMax int
	if max > MAX_VALUE {
		intMax = MAX_VALUE
	} else {
		intMax = int(max)
	}
	return *r.referenceInstances(r.TypeID, intMax)
}

func (r *ReferenceTypeImpl) GetMinorVersion() int {
	_, minor := r.referenceClassFileVersion(r.TypeID)
	return int(minor)
}

func (r *ReferenceTypeImpl) getSuperClass() jdi.ClassType {
	var out jdi.ClassType
	switch r.Kind {
	case jdi.ClassTypeTag:
		sup := r.classTypeSuperclass(jdi.ClassID(r.TypeID))
		if sup != nil {
			out = sup
		}
		r.initSuperClass = true
	case jdi.InterfaceTypeTag:
		out = nil
		r.initSuperClass = true
	}
	return out
}
func (r *ReferenceTypeImpl) isOneDimensionalPrimitiveArray(sign string) bool {
	i := strings.LastIndex(sign, "[")
	var isPa bool
	if i < 0 || strings.HasPrefix(sign, "[[") {
		isPa = false
	} else {
		char := sign[i : i+1]
		isPa = char != "L"
	}
	return isPa
}
