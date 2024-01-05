package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type FieldImpl struct {
	*TypeComponentImpl
}

func (f *FieldImpl) GetUniqueID() jdi.ObjectID {
	return f.RefId
}

func (f *FieldImpl) GetDeclaringType() jdi.ReferenceType {
	return f.DeclaringType
}

func (f *FieldImpl) GetTypeName() string {
	if f.GetSignature() == "" {
		f.Name = jdi.TranslateSignatureToClassName(f.GetSignature())
	}
	return f.Name
}

func (f *FieldImpl) GetType() jdi.Type {
	enclosing := f.GetDeclaringType()
	return findType(enclosing, f.GetSignature())
}

func (f *FieldImpl) IsTransient() bool {
	return f.isModifierSet(TRANSIENT)
}

func (f *FieldImpl) IsVolatile() bool {
	return f.isModifierSet(VOLATILE)
}

func (f *FieldImpl) IsEnumConstant() bool {
	return f.isModifierSet(ENUM_CONSTANT)
}
