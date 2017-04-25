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

type dict struct {
	hashMap *HashMap
}

func NewDict(hashMap *HashMap) Dict {
	return &dict{hashMap}
}

func (d *dict) compositeMarker() {}

func (d *dict) TypeOf() (Type, Error) {
	return TDICT, nil
}

func (d *dict) ToStr() (Str, Error) {
	if d.hashMap.Len().IntVal() == 0 {
		return MakeStr("dict {}"), nil
	}

	var buf bytes.Buffer
	buf.WriteString("dict {")
	idx := 0
	itr := d.hashMap.Iterator()

	for itr.Next() {
		entry := itr.Get()
		if idx > 0 {
			buf.WriteString(",")
		}
		idx++

		buf.WriteString(" ")
		s, err := entry.Key.ToStr()
		if err != nil {
			return s, err
		}
		buf.WriteString(s.String())

		buf.WriteString(": ")
		s, err = entry.Value.ToStr()
		if err != nil {
			return s, err
		}
		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return MakeStr(buf.String()), nil
}

func (d *dict) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (d *dict) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *dict:
		return MakeBool(reflect.DeepEqual(d.hashMap, t.hashMap)), nil
	default:
		return FALSE, nil
	}
}

func (d *dict) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (d *dict) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(d, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (d *dict) Get(key Value) (Value, Error) {
	return d.hashMap.Get(key)
}

func (d *dict) Set(key Value, val Value) Error {
	return d.hashMap.Put(key, val)
}

func (d *dict) Len() (Int, Error) {
	return d.hashMap.Len(), nil
}
