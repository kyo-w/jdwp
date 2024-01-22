package impl

import (
	jdi "github.com/kyo-w/jdwp"
)

type ObjectReferenceImpl struct {
	*MirrorImpl
	// 必须
	ObjectId       jdi.ObjectID
	refType        *ReferenceTypeImpl
	gcDisableCount int
}

func (o *ObjectReferenceImpl) GetUniqueID() jdi.ObjectID {
	return o.ObjectId
}

func (o *ObjectReferenceImpl) GetType() jdi.Type {
	return o.GetReferenceType()
}

func (o *ObjectReferenceImpl) GetReferenceType() jdi.ReferenceType {
	return o.objectReferenceReferenceType(o.ObjectId)
}

func (o *ObjectReferenceImpl) GetValueByField(field jdi.Field) jdi.Value {
	out := o.GetValuesByFields([]jdi.Field{field})
	return out[field]
}

func (o *ObjectReferenceImpl) GetValuesByFields(fields []jdi.Field) map[jdi.Field]jdi.Value {
	fieldId := make([]jdi.FieldID, len(fields))
	for index, value := range fields {
		fieldId[index] = jdi.FieldID(value.GetUniqueID())
	}
	fieldResult := o.objectReferenceGetValues(o.ObjectId, fieldId)
	var out = make(map[jdi.Field]jdi.Value)
	for index, value := range *fieldResult {
		out[fields[index]] = value
	}
	return out
}
func (o *ObjectReferenceImpl) InvokeMethod(thread jdi.ThreadReference, method jdi.Method, args []jdi.Value, options jdi.InvokeOptions) (jdi.Value, jdi.ObjectReference) {
	referType := o.GetReferenceType()
	if method.IsStatic() {
		interfaceType, isInterface := referType.(jdi.InterfaceType)
		if isInterface {
			return interfaceType.InvokeMethod(thread, method, args, options)
		}
		classType, isClassType := referType.(jdi.ClassType)
		if isClassType {
			return classType.InvokeMethod(thread, method, args, options)
		}
		panic("unknown method type")
	}
	return invokeObjectMethod(o.MirrorImpl, o.ObjectId, thread, jdi.ClassID(referType.GetUniqueID()), method, args, options)
}

func (o *ObjectReferenceImpl) DisableCollection() {
	if o.gcDisableCount == 0 {
		o.objectReferenceDisableCollection(o.ObjectId)
	}
	o.gcDisableCount++
}

func (o *ObjectReferenceImpl) EnableCollection() {
	o.gcDisableCount--
	if o.gcDisableCount == 0 {
		o.objectReferenceEnableCollection(o.ObjectId)
	}
}

func (o *ObjectReferenceImpl) IsCollected() bool {
	return o.objectReferenceIsCollected(o.ObjectId)
}

func (o *ObjectReferenceImpl) UniqueID() jdi.ObjectID {
	return o.ObjectId
}

func (o *ObjectReferenceImpl) GetReferringObjects(maxReferrers int) []jdi.ObjectReference {
	return *o.objectReferenceReferringObjects(o.ObjectId, maxReferrers)
}
func (o *ObjectReferenceImpl) GetTagType() jdi.Tag {
	return jdi.OBJECT
}
func (o *ObjectReferenceImpl) GetValuesByFieldNames(fieldNames ...string) jdi.Value {
	var tmpObjectRef jdi.Value
	tmpObjectRef = o
	for _, fieldName := range fieldNames {
		objectRef, isObject := tmpObjectRef.(jdi.ObjectReference)
		if !isObject {
			panic(tmpObjectRef.GetType().GetSignature() + " is not object")
		}
		fieldRef := objectRef.GetType().(jdi.ClassType).GetFieldByName(fieldName)
		tmpObjectRef = objectRef.GetValueByField(fieldRef)
	}
	return tmpObjectRef
}
