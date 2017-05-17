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
	//"reflect"
)

//---------------------------------------------------------------
// StructEntry

type StructEntry struct {
	Key   string
	Value Value
}

//---------------------------------------------------------------
// structMap

type structMap interface {
	get(Str) (Value, Error)
	put(Str, Value) Error
	has(Value) (Bool, Error)
	keys() []string
}

//---------------------------------------------------------------
// struct

type _struct struct {
	smap structMap
}

func (stc *_struct) compositeMarker() {}

func (stc *_struct) TypeOf() Type { return TSTRUCT }

func (stc *_struct) ToStr() Str {

	var buf bytes.Buffer
	buf.WriteString("struct {")
	for i, k := range stc.keys() {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString(": ")

		v, err := stc.GetField(str(k))
		if err != nil {
			panic("invalid struct")
		}
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

	// same type
	that, ok := v.(Struct)
	if !ok {
		return FALSE
	}

	// same number of keys
	keys := stc.keys()
	if len(keys) != len(that.keys()) {
		return FALSE
	}

	// all keys have same value
	for _, k := range keys {
		a, err := stc.GetField(str(k))
		if err != nil {
			panic("invalid chain")
		}

		b, err := that.GetField(str(k))
		if err != nil {
			if err.Kind() == "NoSuchField" {
				return FALSE
			} else {
				panic("invalid chain")
			}
		}

		if a.Eq(b) != TRUE {
			return FALSE
		}
	}

	// done
	return TRUE
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
	return stc.smap.get(key)
}

func (stc *_struct) PutField(key Str, val Value) Error {
	return stc.smap.put(key, val)
}

func (stc *_struct) Has(key Value) (Bool, Error) {
	return stc.smap.has(key)
}

func (stc *_struct) keys() []string {
	return stc.smap.keys()
}

//---------------------------------------------------------------
// simple struct

// TODO replace with a more efficient data structure
type simpleStruct struct {
	fields map[string]Value
}

func NewStruct(entries []*StructEntry) (Struct, Error) {
	ss := &simpleStruct{make(map[string]Value)}
	for _, e := range entries {
		if _, has := ss.fields[e.Key]; has {
			return nil, DuplicateFieldError(e.Key)
		}

		ss.fields[e.Key] = e.Value
	}
	return &_struct{ss}, nil
}

func BlankStruct(keys []string) (Struct, Error) {
	ss := &simpleStruct{make(map[string]Value)}
	for _, k := range keys {
		if _, has := ss.fields[k]; has {
			return nil, DuplicateFieldError(k)
		}
		ss.fields[k] = NULL
	}
	return &_struct{ss}, nil
}

func (ss *simpleStruct) get(key Str) (Value, Error) {
	v, ok := ss.fields[key.String()]
	if ok {
		return v, nil
	} else {
		return nil, NoSuchFieldError(key.String())
	}
}

func (ss *simpleStruct) put(key Str, val Value) Error {
	_, ok := ss.fields[key.String()]
	if ok {
		ss.fields[key.String()] = val
		return nil
	} else {
		return NoSuchFieldError(key.String())
	}
}

func (ss *simpleStruct) has(key Value) (Bool, Error) {
	if s, ok := key.(Str); ok {
		_, has := ss.fields[s.String()]
		return MakeBool(has), nil
	} else {
		return nil, TypeMismatchError("Expected 'Str'")
	}
}

func (ss *simpleStruct) keys() []string {

	keys := make([]string, len(ss.fields), len(ss.fields))
	idx := 0
	for k := range ss.fields {
		keys[idx] = k
		idx++
	}
	return keys
}
