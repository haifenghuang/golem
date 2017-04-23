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
// An ObjDef contains the information needed to instantiate an Obj
// instance.  ObjDefs are created at compile time, and
// are immutable at run time.

type ObjDef struct {
	Keys []string
}

//---------------------------------------------------------------
// obj

type obj struct {
	// TODO replace this with a more efficient data structure
	fields map[string]Value
	inited bool
}

func NewObj() Obj {
	return &obj{nil, false}
}

func (o *obj) Init(def *ObjDef, vals []Value) {
	o.fields = make(map[string]Value)
	for i, k := range def.Keys {
		o.fields[k] = vals[i]
	}
	o.inited = true
}

func (o *obj) compositeMarker() {}

func (o *obj) TypeOf() (Type, Error) {
	if !o.inited {
		return TOBJ, UninitializedObjError()
	}

	return TOBJ, nil
}

func (o *obj) ToStr() (Str, Error) {
	if !o.inited {
		return nil, UninitializedObjError()
	}

	if len(o.fields) == 0 {
		return MakeStr("obj {}"), nil
	}

	var buf bytes.Buffer
	buf.WriteString("obj {")
	idx := 0
	for k, v := range o.fields {
		if idx > 0 {
			buf.WriteString(",")
		}
		idx = idx + 1
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString(": ")

		s, err := v.ToStr()
		if err != nil {
			return nil, err
		}
		buf.WriteString(s.String())
	}
	buf.WriteString(" }")
	return MakeStr(buf.String()), nil
}

func (o *obj) HashCode() (Int, Error) {
	if !o.inited {
		return nil, UninitializedObjError()
	}

	// TODO $hash()
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (o *obj) Eq(v Value) (Bool, Error) {
	if !o.inited {
		return FALSE, UninitializedObjError()
	}

	// TODO $eq()
	switch t := v.(type) {
	case *obj:
		return MakeBool(reflect.DeepEqual(o.fields, t.fields)), nil
	default:
		return FALSE, nil
	}
}

func (o *obj) Cmp(v Value) (Int, Error) {
	if !o.inited {
		return nil, UninitializedObjError()
	}

	return nil, TypeMismatchError("Expected Comparable Type")
}

func (o *obj) Add(v Value) (Value, Error) {
	if !o.inited {
		return nil, UninitializedObjError()
	}

	switch t := v.(type) {

	case Str:
		return strcat(o, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (o *obj) Get(index Value) (Value, Error) {
	if !o.inited {
		return nil, UninitializedObjError()
	}

	if s, ok := index.(Str); ok {
		return o.GetField(s)
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}

func (o *obj) Set(index Value, val Value) Error {
	if !o.inited {
		return UninitializedObjError()
	}

	if s, ok := index.(Str); ok {
		return o.PutField(s, val)
	} else {
		return TypeMismatchError("Expected 'Str'")
	}
}

func (o *obj) GetField(key Str) (Value, Error) {
	if !o.inited {
		return nil, UninitializedObjError()
	}

	v, ok := o.fields[key.String()]
	if ok {
		return v, nil
	} else {
		return nil, NoSuchFieldError(key.String())
	}
}

func (o *obj) PutField(key Str, val Value) Error {
	if !o.inited {
		return UninitializedObjError()
	}

	_, ok := o.fields[key.String()]
	if ok {
		o.fields[key.String()] = val
		return nil
	} else {
		return NoSuchFieldError(key.String())
	}
}

func (o *obj) Has(key Value) (Bool, Error) {
	if !o.inited {
		return nil, UninitializedObjError()
	}

	if s, ok := key.(Str); ok {
		_, has := o.fields[s.String()]
		return MakeBool(has), nil
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}
