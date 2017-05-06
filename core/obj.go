// Copyright 2017 The Golem Project Developers //
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

type ObjEntry struct {
	Key   string
	Value Value
}

//---------------------------------------------------------------
// obj

type obj struct {
	// TODO replace this with a more efficient data structure
	fields map[string]Value
}

func NewObj(entries []*ObjEntry) Obj {
	o := &obj{make(map[string]Value)}
	for _, e := range entries {
		o.fields[e.Key] = e.Value
	}
	return o
}

func BlankObj(keys []string) Obj {
	o := &obj{make(map[string]Value)}
	for _, k := range keys {
		o.fields[k] = NULL
	}
	return o
}

func (o *obj) compositeMarker() {}

func (o *obj) TypeOf() Type { return TOBJ }

func (o *obj) ToStr() Str {
	if len(o.fields) == 0 {
		return MakeStr("obj {}")
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

		buf.WriteString(v.ToStr().String())
	}
	buf.WriteString(" }")
	return MakeStr(buf.String())
}

func (o *obj) HashCode() (Int, Error) {
	// TODO $hash()
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (o *obj) Eq(v Value) Bool {
	// TODO $eq()
	switch t := v.(type) {
	case *obj:
		return MakeBool(reflect.DeepEqual(o.fields, t.fields))
	default:
		return FALSE
	}
}

func (o *obj) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (o *obj) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return Strcat(o, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (o *obj) Get(index Value) (Value, Error) {
	if s, ok := index.(Str); ok {
		return o.GetField(s)
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}

func (o *obj) Set(index Value, val Value) Error {
	if s, ok := index.(Str); ok {
		return o.PutField(s, val)
	} else {
		return TypeMismatchError("Expected 'Str'")
	}
}

func (o *obj) GetField(key Str) (Value, Error) {
	v, ok := o.fields[key.String()]
	if ok {
		return v, nil
	} else {
		return nil, NoSuchFieldError(key.String())
	}
}

func (o *obj) PutField(key Str, val Value) Error {
	_, ok := o.fields[key.String()]
	if ok {
		o.fields[key.String()] = val
		return nil
	} else {
		return NoSuchFieldError(key.String())
	}
}

func (o *obj) Has(key Value) (Bool, Error) {
	if s, ok := key.(Str); ok {
		_, has := o.fields[s.String()]
		return MakeBool(has), nil
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}
