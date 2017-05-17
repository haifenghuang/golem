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
	//"fmt"
	"reflect"
)

type StructEntry struct {
	Key   string
	Value Value
}

//---------------------------------------------------------------
// struct

// TODO replace with a more efficient data structure
type _struct struct {
	fields map[string]Value
}

func NewStruct(entries []*StructEntry) (Struct, Error) {
	stc := &_struct{make(map[string]Value)}
	for _, e := range entries {
		if _, has := stc.fields[e.Key]; has {
			return nil, DuplicateFieldError(e.Key)
		}

		stc.fields[e.Key] = e.Value
	}
	return stc, nil
}

func BlankStruct(keys []string) (Struct, Error) {
	stc := &_struct{make(map[string]Value)}
	for _, k := range keys {
		if _, has := stc.fields[k]; has {
			return nil, DuplicateFieldError(k)
		}
		stc.fields[k] = NULL
	}
	return stc, nil
}

func (stc *_struct) compositeMarker() {}

func (stc *_struct) TypeOf() Type { return TSTRUCT }

func (stc *_struct) ToStr() Str {

	var buf bytes.Buffer
	buf.WriteString("struct {")
	idx := 0
	for k, v := range stc.fields {
		if idx > 0 {
			buf.WriteString(",")
		}
		idx++
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString(": ")

		buf.WriteString(v.ToStr().String())
	}
	buf.WriteString(" }")
	return MakeStr(buf.String())
}

func (stc *_struct) HashCode() (Int, Error) {
	// TODO $hash()
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (stc *_struct) Eq(v Value) Bool {

	// TODO $eq()
	switch t := v.(type) {
	case *_struct:
		return MakeBool(reflect.DeepEqual(stc, t))
	default:
		return FALSE
	}
}

func (stc *_struct) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (stc *_struct) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(stc, t), nil

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (stc *_struct) Get(index Value) (Value, Error) {
	if s, ok := index.(Str); ok {
		return stc.GetField(s)
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}

func (stc *_struct) Set(index Value, val Value) Error {
	if s, ok := index.(Str); ok {
		return stc.PutField(s, val)
	} else {
		return TypeMismatchError("Expected 'Str'")
	}
}

func (stc *_struct) GetField(key Str) (Value, Error) {
	v, ok := stc.fields[key.String()]
	if ok {
		return v, nil
	} else {
		return nil, NoSuchFieldError(key.String())
	}
}

func (stc *_struct) PutField(key Str, val Value) Error {
	_, ok := stc.fields[key.String()]
	if ok {
		stc.fields[key.String()] = val
		return nil
	} else {
		return NoSuchFieldError(key.String())
	}
}

func (stc *_struct) Has(key Value) (Bool, Error) {
	if s, ok := key.(Str); ok {
		_, has := stc.fields[s.String()]
		return MakeBool(has), nil
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}

func (stc *_struct) keys() []string {

	keys := make([]string, len(stc.fields), len(stc.fields))
	idx := 0
	for k := range stc.fields {
		keys[idx] = k
		idx++
	}
	return keys
}
