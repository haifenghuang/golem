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
	"reflect"
)

//---------------------------------------------------------------
// tuple

type tuple struct {
	array []g.Value
}

func NewTuple(values []g.Value) Tuple {
	if len(values) < 2 {
		panic("invalid tuple size")
	}
	return &tuple{values}
}

func (tp *tuple) compositeMarker() {}

func (tp *tuple) TypeOf() (g.Type, g.Error) {
	return g.TTUPLE, nil
}

func (tp *tuple) ToStr() (g.Str, g.Error) {
	var buf bytes.Buffer
	buf.WriteString("(")
	for idx, v := range tp.array {
		if idx > 0 {
			buf.WriteString(", ")
		}
		s, err := v.ToStr()
		if err != nil {
			return nil, err
		}
		buf.WriteString(s.String())
	}
	buf.WriteString(")")
	return g.MakeStr(buf.String()), nil
}

func (tp *tuple) HashCode() (g.Int, g.Error) {

	// https://en.wikipedia.org/wiki/Jenkins_hash_function
	var hash int64 = 0
	for _, v := range tp.array {
		h, err := v.HashCode()
		if err != nil {
			return nil, err
		}
		hash += h.IntVal()
		hash += hash << 10
		hash ^= hash >> 6
	}
	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15
	return g.MakeInt(hash), nil
}

func (tp *tuple) Eq(v g.Value) (g.Bool, g.Error) {
	switch t := v.(type) {
	case *tuple:
		return g.MakeBool(reflect.DeepEqual(tp.array, t.array)), nil
	default:
		return g.FALSE, nil
	}
}

func (tp *tuple) Cmp(v g.Value) (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Comparable Type")
}

func (tp *tuple) Add(v g.Value) (g.Value, g.Error) {
	switch t := v.(type) {

	case g.Str:
		return g.Strcat(tp, t)

	default:
		return nil, g.TypeMismatchError("Expected Number Type")
	}
}

func (tp *tuple) Get(index g.Value) (g.Value, g.Error) {
	idx, err := g.ParseIndex(index, len(tp.array))
	if err != nil {
		return nil, err
	}
	return tp.array[idx.IntVal()], nil
}

func (tp *tuple) Len() (g.Int, g.Error) {
	return g.MakeInt(int64(len(tp.array))), nil
}
