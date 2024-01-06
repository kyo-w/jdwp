package impl

import (
	jdi "github.com/kyo-w/jdwp"
	connect "github.com/kyo-w/jdwp/impl/internal"
)

type BooleanValueImpl struct {
	*MirrorImpl
	value bool
}
type ByteValueImpl struct {
	*MirrorImpl
	value byte
}
type CharValueImpl struct {
	*MirrorImpl
	value jdi.Char
}
type DoubleValueImpl struct {
	*MirrorImpl
	value jdi.Double
}

type FloatValueImpl struct {
	*MirrorImpl
	value jdi.Float
}
type IntegerValueImpl struct {
	*MirrorImpl
	value jdi.Int
}

type LongValueImpl struct {
	*MirrorImpl
	value jdi.Long
}
type ShortValueImpl struct {
	*MirrorImpl
	value jdi.Short
}

type VoidValueImpl struct {
	*MirrorImpl
}

type StringReferenceImpl struct {
	*ObjectReferenceImpl
	value string
}

func (s *StringReferenceImpl) GetStringValue() string {
	if s.value == "" {
		err := s.GetConnect().SendCommand(connect.CmdStringReferenceValue, &s.ObjectId, &s.value)
		if err != nil {
			panic(err)
		}
	}
	return s.value
}
func (s *StringReferenceImpl) GetTagType() jdi.Tag {
	return jdi.STRING
}

func (b *BooleanValueImpl) GetType() jdi.Type {
	return &jdi.BooleanType{Vm: b.vm}
}
func (b *BooleanValueImpl) GetValue() bool {
	return b.value
}
func (b *BooleanValueImpl) GetTagType() jdi.Tag {
	return jdi.BOOLEAN
}

func (b *ByteValueImpl) GetType() jdi.Type {
	return &jdi.BooleanType{Vm: b.vm}
}
func (b *ByteValueImpl) GetValue() byte {
	return b.value
}
func (b *ByteValueImpl) GetTagType() jdi.Tag {
	return jdi.BYTE
}

func (c *CharValueImpl) GetType() jdi.Type {
	return &jdi.CharType{Vm: c.vm}
}
func (c *CharValueImpl) GetValue() jdi.Char {
	return c.value
}
func (c *CharValueImpl) GetTagType() jdi.Tag {
	return jdi.CHAR
}

func (d *DoubleValueImpl) GetType() jdi.Type {
	return &jdi.DoubleType{Vm: d.vm}
}
func (d *DoubleValueImpl) GetValue() jdi.Double {
	return d.value
}
func (d *DoubleValueImpl) GetTagType() jdi.Tag {
	return jdi.DOUBLE
}

func (f *FloatValueImpl) GetType() jdi.Type {
	return &jdi.FloatType{Vm: f.vm}
}
func (f *FloatValueImpl) GetValue() jdi.Float {
	return f.value
}
func (f *FloatValueImpl) GetTagType() jdi.Tag {
	return jdi.FLOAT
}

func (i *IntegerValueImpl) GetType() jdi.Type {
	return &jdi.IntegerType{Vm: i.vm}
}
func (i *IntegerValueImpl) GetValue() jdi.Int {
	return i.value
}
func (i *IntegerValueImpl) GetTagType() jdi.Tag {
	return jdi.INT
}

func (l *LongValueImpl) GetType() jdi.Type {
	return &jdi.LongType{Vm: l.vm}
}
func (l *LongValueImpl) GetValue() jdi.Long {
	return l.value
}
func (l *LongValueImpl) GetTagType() jdi.Tag {
	return jdi.LONG
}

func (s *ShortValueImpl) GetType() jdi.Type {
	return &jdi.ShortType{Vm: s.vm}
}
func (s *ShortValueImpl) GetValue() jdi.Short {
	return s.value
}
func (s *ShortValueImpl) GetTagType() jdi.Tag {
	return jdi.SHORT
}

// Void 占位符Empty
func (v *VoidValueImpl) Void() {
}
func (v *VoidValueImpl) GetType() jdi.Type {
	return &jdi.VoidType{Vm: v.vm}
}
func (v *VoidValueImpl) GetTagType() jdi.Tag {
	return jdi.VOID
}
