// Copyright (C) 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"fmt"
	jdi "github.com/kyo-w/jdwp"
	"reflect"
)

const (
	// StatusVerified is used to describe a class in the verified state.
	StatusVerified = jdi.ByteID(1)
	// StatusPrepared is used to describe a class in the prepared state.
	StatusPrepared = jdi.ByteID(2)
	// StatusInitialized is used to describe a class in the initialized state.
	StatusInitialized = jdi.ByteID(4)
	// StatusError is used to describe a class in the error state.
	StatusError = jdi.ByteID(8)
)

func unbox(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		return v.Elem()
	}
	return v
}

// encode writes the value v to w, using the JDWP encoding scheme.
func (c *Connection) encode(w Writer, v reflect.Value) error {

	t := v.Type()
	o := v.Interface()

	switch v.Type() {
	case reflect.TypeOf((*jdi.EventModifier)(nil)).Elem():
		w.Uint8(o.(jdi.EventModifier).ModKind())

	case reflect.TypeOf((*jdi.ValueID)(nil)).Elem():
		switch o.(type) {
		case jdi.ArrayID:
			w.Uint8(uint8(jdi.ARRAY))
		case byte:
			w.Uint8(uint8(jdi.BYTE))
		case jdi.Char:
			w.Uint8(uint8(jdi.CHAR))
		case jdi.ObjectID:
			w.Uint8(uint8(jdi.OBJECT))
		case float32:
			w.Uint8(uint8(jdi.FLOAT))
		case float64:
			w.Uint8(uint8(jdi.DOUBLE))
		case int, int32:
			w.Uint8(uint8(jdi.INT))
		case int16:
			w.Uint8(uint8(jdi.SHORT))
		case int64:
			w.Uint8(uint8(jdi.LONG))
		case nil:
			w.Uint8(uint8(jdi.VOID))
		case bool:
			w.Uint8(uint8(jdi.BOOLEAN))
		case jdi.StringID:
			w.Uint8(uint8(jdi.STRING))
		case jdi.ThreadID:
			w.Uint8(uint8(jdi.THREAD))
		case jdi.ThreadGroupID:
			w.Uint8(uint8(jdi.ThreadGroup))
		case jdi.ClassLoaderID:
			w.Uint8(uint8(jdi.ClassLoader))
		case jdi.ClassObjectID:
			w.Uint8(uint8(jdi.ClassObject))
		default:
			panic(fmt.Errorf("Got Value of api %T", o))
		}
	}

	switch o := o.(type) {
	case jdi.ReferenceTypeID, jdi.ClassID, jdi.InterfaceID, jdi.ArrayTypeID:
		WriteUint(w, c.idSizes.ReferenceTypeIDSize*8, unbox(v).Uint())

	case jdi.MethodID:
		WriteUint(w, c.idSizes.MethodIDSize*8, unbox(v).Uint())

	case jdi.FieldID:
		WriteUint(w, c.idSizes.FieldIDSize*8, unbox(v).Uint())

	case jdi.ObjectID, jdi.ThreadID, jdi.ThreadGroupID, jdi.StringID, jdi.ClassLoaderID, jdi.ClassObjectID, jdi.ArrayID:
		WriteUint(w, c.idSizes.ObjectIDSize*8, unbox(v).Uint())

	case []byte: // Optimisation
		w.Uint32(uint32(len(o)))
		w.Data(o)

	default:
		switch t.Kind() {
		case reflect.Ptr, reflect.Interface:
			return c.encode(w, v.Elem())
		case reflect.String:
			w.Uint32(uint32(v.Len()))
			w.Data([]byte(v.String()))
		case reflect.Uint8:
			w.Uint8(uint8(v.Uint()))
		case reflect.Uint64:
			w.Uint64(uint64(v.Uint()))
		case reflect.Int8:
			w.Int8(int8(v.Int()))
		case reflect.Int16:
			w.Int16(int16(v.Int()))
		case reflect.Int32, reflect.Int:
			w.Int32(int32(v.Int()))
		case reflect.Int64:
			w.Int64(v.Int())
		case reflect.Float32:
			w.Float32(float32(v.Float()))
		case reflect.Float64:
			w.Float64(v.Float())
		case reflect.Bool:
			w.Bool(v.Bool())
		case reflect.Struct:
			for i, count := 0, v.NumField(); i < count; i++ {
				c.encode(w, v.Field(i))
			}
		case reflect.Slice:
			count := v.Len()
			w.Uint32(uint32(count))
			for i := 0; i < count; i++ {
				c.encode(w, v.Index(i))
			}
		default:
			panic(fmt.Errorf("Unhandled api %T %v %v", o, t.Name(), t.Kind()))
		}
	}
	return w.Error()
}

// decode reads the value v from r, using the JDWP encoding scheme.
func (c *Connection) decode(r Reader, v reflect.Value) error {
	switch v.Type() {
	// event事件
	case reflect.TypeOf((*jdi.EventResponse)(nil)).Elem():
		var kind jdi.EventKind
		if err := c.decode(r, reflect.ValueOf(&kind)); err != nil {
			return err
		}
		event := kind.Event()
		v.Set(reflect.ValueOf(event))
		v = v.Elem()
	case reflect.TypeOf((*jdi.ArrayRegion)(nil)).Elem():
		tag := jdi.Tag(r.Uint8())
		var ty reflect.Type
		switch tag {
		case jdi.OBJECT, jdi.ARRAY, jdi.STRING, jdi.THREAD, jdi.ThreadGroup, jdi.ClassLoader, jdi.ClassObject:
			ty = reflect.TypeOf([]jdi.TaggedObjectID{})
			count := int(r.Uint32())
			slice := reflect.MakeSlice(ty, count, count)
			for i := 0; i < count; i++ {
				c.decode(r, slice.Index(i))
			}
			v.Field(2).Set(slice)
			return r.Error()
		default:
			ty = reflect.TypeOf([]jdi.ValueID{})
			count := int(r.Uint32())
			slice := reflect.MakeSlice(ty, count, count)
			for i := 0; i < count; i++ {
				c.decode(r, slice.Index(i))
			}
			v.Field(1).Set(slice)
			return r.Error()
		}

	case reflect.TypeOf((*jdi.ValueID)(nil)).Elem():
		tag := jdi.Tag(r.Uint8())
		var ty reflect.Type
		switch tag {
		case jdi.ARRAY:
			ty = reflect.TypeOf(jdi.ArrayID(0))
		case jdi.BYTE:
			ty = reflect.TypeOf(byte(0))
		case jdi.CHAR:
			ty = reflect.TypeOf(jdi.Char(0))
		case jdi.OBJECT:
			ty = reflect.TypeOf(jdi.ObjectID(0))
		case jdi.FLOAT:
			ty = reflect.TypeOf(float32(0))
		case jdi.DOUBLE:
			ty = reflect.TypeOf(float64(0))
		case jdi.INT:
			ty = reflect.TypeOf(0)
		case jdi.SHORT:
			ty = reflect.TypeOf(int16(0))
		case jdi.LONG:
			ty = reflect.TypeOf(int64(0))
		case jdi.BOOLEAN:
			ty = reflect.TypeOf(false)
		case jdi.STRING:
			ty = reflect.TypeOf(jdi.StringID(0))
		case jdi.THREAD:
			ty = reflect.TypeOf(jdi.ThreadID(0))
		case jdi.ThreadGroup:
			ty = reflect.TypeOf(jdi.ThreadGroupID(0))
		case jdi.ClassLoader:
			ty = reflect.TypeOf(jdi.ClassLoaderID(0))
		case jdi.ClassObject:
			ty = reflect.TypeOf(jdi.ClassObjectID(0))
		case jdi.VOID:
			v.Set(reflect.New(v.Type()).Elem())
			return r.Error()
		default:
			panic(fmt.Errorf("Unhandled value api %v", tag))
		}
		data := reflect.New(ty).Elem()
		c.decode(r, data)
		v.Set(data)
		return r.Error()
	}
	t := v.Type()
	o := v.Interface()
	switch o := o.(type) {
	case jdi.ReferenceTypeID, jdi.ClassID, jdi.InterfaceID, jdi.ArrayTypeID:
		v.Set(reflect.ValueOf(ReadUint(r, c.idSizes.ReferenceTypeIDSize*8)).Convert(t))
	case jdi.MethodID:
		v.Set(reflect.ValueOf(ReadUint(r, c.idSizes.MethodIDSize*8)).Convert(t))
	case jdi.FieldID:
		v.Set(reflect.ValueOf(ReadUint(r, c.idSizes.FieldIDSize*8)).Convert(t))
	case jdi.ObjectID, jdi.ThreadID, jdi.ThreadGroupID, jdi.StringID, jdi.ClassLoaderID, jdi.ClassObjectID, jdi.ArrayID:
		byteValue := ReadUint(r, c.idSizes.ObjectIDSize*8)
		valueRef := reflect.ValueOf(byteValue).Convert(t)
		v.Set(valueRef)
	case jdi.EventModifier:
		panic("Cannot decode EventModifiers")

	default:
		switch t.Kind() {
		case reflect.Ptr, reflect.Interface:
			return c.decode(r, v.Elem())
		case reflect.String:
			data := make([]byte, r.Uint32())
			r.Data(data)
			v.Set(reflect.ValueOf(string(data)).Convert(t))
		case reflect.Bool:
			v.Set(reflect.ValueOf(r.Bool()).Convert(t))
		case reflect.Uint8:
			v.Set(reflect.ValueOf(r.Uint8()).Convert(t))
		case reflect.Uint64:
			v.Set(reflect.ValueOf(r.Uint64()).Convert(t))
		case reflect.Int8:
			v.Set(reflect.ValueOf(r.Int8()).Convert(t))
		case reflect.Int16:
			v.Set(reflect.ValueOf(r.Int16()).Convert(t))
		case reflect.Int32, reflect.Int:
			v.Set(reflect.ValueOf(r.Int32()).Convert(t))
		case reflect.Int64:
			v.Set(reflect.ValueOf(r.Int64()).Convert(t))
		case reflect.Struct:
			for i, count := 0, v.NumField(); i < count; i++ {
				c.decode(r, v.Field(i))
			}
		case reflect.Slice:
			count := int(r.Uint32())
			slice := reflect.MakeSlice(t, count, count)
			for i := 0; i < count; i++ {
				c.decode(r, slice.Index(i))
			}
			v.Set(slice)
		default:
			panic(fmt.Errorf("Unhandled api %T %v %v", o, t.Name(), t.Kind()))
		}
	}
	return r.Error()
}
