package impl

import jdi "github.com/kyo-w/jdwp"

type TypeComponentImpl struct {
	*MirrorImpl
	//一定存在参数
	RefId     jdi.ObjectID
	Name      string
	Modifiers int

	DeclaringType    jdi.ReferenceType
	signature        string
	genericSignature string
	hasGetModifiers  bool
}

func (t *TypeComponentImpl) GetUniqueID() jdi.ObjectID {
	return t.RefId
}
func (t *TypeComponentImpl) GetModifiers() int {
	return t.Modifiers
}
func (t *TypeComponentImpl) IsPrivate() bool {
	return t.isModifierSet(PRIVATE)
}
func (t *TypeComponentImpl) IsPackagePrivate() bool {
	return t.isModifierSet(PRIVATE | PROTECTED | PUBLIC)

}
func (t *TypeComponentImpl) IsProtected() bool {
	return t.isModifierSet(PROTECTED)

}
func (t *TypeComponentImpl) IsPublic() bool {
	return t.isModifierSet(PUBLIC)
}
func (t *TypeComponentImpl) GetName() string {
	return t.Name
}
func (t *TypeComponentImpl) GetSignature() string {
	return t.signature
}
func (t *TypeComponentImpl) GetGenericSignature() string {
	return t.genericSignature
}
func (t *TypeComponentImpl) GetDeclaringType() jdi.ReferenceType {
	return t.DeclaringType
}
func (t *TypeComponentImpl) IsStatic() bool {
	return t.isModifierSet(STATIC)
}
func (t *TypeComponentImpl) IsFinal() bool {
	return t.isModifierSet(FINAL)
}
func (t *TypeComponentImpl) IsSynthetic() bool {
	return t.isModifierSet(SYNTHETIC)
}
func (t *TypeComponentImpl) isModifierSet(compareBits int) bool {
	return (t.Modifiers & compareBits) != 0
}
