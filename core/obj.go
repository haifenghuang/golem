// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"bytes"
	"reflect"
)

//---------------------------------------------------------------
// ObjDef represents the information needed to instantiate an obj
// instance.  ObjDefs are created at compile time, and
// are immutable at run time.

type ObjDef struct {
	Keys []string
}

//---------------------------------------------------------------
// Obj

type Obj struct {
	Fields map[string]Value
	Inited bool
}

func NewObj() *Obj {
	return &Obj{nil, false}
}

func (o *Obj) Init(def *ObjDef, vals []Value) {
	o.Fields = make(map[string]Value)
	for i, k := range def.Keys {
		o.Fields[k] = vals[i]
	}
	o.Inited = true
}

func (o *Obj) TypeOf() (Type, Error) {
	if !o.Inited {
		return TOBJ, UninitializedObjError()
	}

	return TOBJ, nil
}

func (o *Obj) String() (Str, Error) {
	if !o.Inited {
		return Str(""), UninitializedObjError()
	}

	var buf bytes.Buffer
	buf.WriteString("obj {")
	idx := 0
	for k, v := range o.Fields {
		if idx > 0 {
			buf.WriteString(",")
		}
		idx = idx + 1
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString(": ")

		s, err := v.String()
		if err != nil {
			return Str(""), err
		}
		buf.WriteString(string(s))
	}
	buf.WriteString(" }")
	return Str(buf.String()), nil
}

func (o *Obj) Eq(v Value) (Bool, Error) {
	if !o.Inited {
		return Bool(false), UninitializedObjError()
	}

	switch t := v.(type) {
	case *Obj:
		return Bool(reflect.DeepEqual(o.Fields, t.Fields)), nil
	default:
		return Bool(false), nil
	}
}

func (o *Obj) Cmp(v Value) (Int, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}

	return Int(0), TypeMismatchError("Expected Comparable Type")
}

func (o *Obj) Add(v Value) (Value, Error) {
	if !o.Inited {
		return nil, UninitializedObjError()
	}

	switch t := v.(type) {

	case Str:
		return strcat([]Value{o, t})

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (o *Obj) Sub(v Value) (Number, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return nil, TypeMismatchError("Expected Number Type")
}
func (o *Obj) Mul(v Value) (Number, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return nil, TypeMismatchError("Expected Number Type")
}
func (o *Obj) Div(v Value) (Number, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return nil, TypeMismatchError("Expected Number Type")
}
func (o *Obj) Rem(v Value) (Int, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return Int(0), TypeMismatchError("Expected 'Int'")
}
func (o *Obj) BitAnd(v Value) (Int, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return Int(0), TypeMismatchError("Expected 'Int'")
}
func (o *Obj) BitOr(v Value) (Int, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return Int(0), TypeMismatchError("Expected 'Int'")
}
func (o *Obj) BitXOr(v Value) (Int, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return Int(0), TypeMismatchError("Expected 'Int'")
}
func (o *Obj) LeftShift(v Value) (Int, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return Int(0), TypeMismatchError("Expected 'Int'")
}
func (o *Obj) RightShift(Value) (Int, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return Int(0), TypeMismatchError("Expected 'Int'")
}
func (o *Obj) Negate() (Number, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return Int(0), TypeMismatchError("Expected Number Type")
}
func (o *Obj) Not() (Bool, Error) {
	if !o.Inited {
		return Bool(false), UninitializedObjError()
	}
	return false, TypeMismatchError("Expected 'Bool'")
}
func (o *Obj) Complement() (Int, Error) {
	if !o.Inited {
		return Int(0), UninitializedObjError()
	}
	return Int(0), TypeMismatchError("Expected 'Int'")
}

func (o *Obj) GetField(key string) (Value, Error) {
	if !o.Inited {
		return nil, UninitializedObjError()
	}

	v, ok := o.Fields[key]
	if ok {
		return v, nil
	} else {
		return nil, NoSuchFieldError(key)
	}
}

func (o *Obj) PutField(key string, val Value) Error {
	if !o.Inited {
		return UninitializedObjError()
	}

	_, ok := o.Fields[key]
	if ok {
		o.Fields[key] = val
		return nil
	} else {
		return NoSuchFieldError(key)
	}
}
