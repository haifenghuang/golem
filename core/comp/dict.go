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

package comp

import (
	"bytes"
	g "golem/core"
	"golem/core/coll"
	"reflect"
)

type dict struct {
	hashMap *coll.HashMap
}

func NewDict(hashMap *coll.HashMap) Dict {
	return &dict{hashMap}
}

func (d *dict) compositeMarker() {}

func (d *dict) TypeOf() (g.Type, g.Error) {
	return g.TDICT, nil
}

func (d *dict) ToStr() (g.Str, g.Error) {
	if d.hashMap.Len().IntVal() == 0 {
		return g.MakeStr("dict {}"), nil
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
	return g.MakeStr(buf.String()), nil
}

func (d *dict) HashCode() (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Hashable Type")
}

func (d *dict) Eq(v g.Value) (g.Bool, g.Error) {
	switch t := v.(type) {
	case *dict:
		return g.MakeBool(reflect.DeepEqual(d.hashMap, t.hashMap)), nil
	default:
		return g.FALSE, nil
	}
}

func (d *dict) Cmp(v g.Value) (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Comparable Type")
}

func (d *dict) Add(v g.Value) (g.Value, g.Error) {
	switch t := v.(type) {

	case g.Str:
		return g.Strcat(d, t)

	default:
		return nil, g.TypeMismatchError("Expected Number Type")
	}
}

func (d *dict) Get(key g.Value) (g.Value, g.Error) {
	return d.hashMap.Get(key)
}

func (d *dict) Set(key g.Value, val g.Value) g.Error {
	return d.hashMap.Put(key, val)
}

func (d *dict) Len() (g.Int, g.Error) {
	return d.hashMap.Len(), nil
}
