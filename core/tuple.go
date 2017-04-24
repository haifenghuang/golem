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
// tuple

type tuple struct {
	array []Value
}

func NewTuple(values []Value) Tuple {
	return &tuple{values}
}

func (tp *tuple) compositeMarker() {}

func (tp *tuple) TypeOf() (Type, Error) {
	return TTUPLE, nil
}

func (tp *tuple) ToStr() (Str, Error) {

	if len(tp.array) == 0 {
		return MakeStr("()"), nil
	}

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
	return MakeStr(buf.String()), nil
}

func (tp *tuple) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (tp *tuple) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *tuple:
		return MakeBool(reflect.DeepEqual(tp.array, t.array)), nil
	default:
		return FALSE, nil
	}
}

func (tp *tuple) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (tp *tuple) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(tp, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (tp *tuple) Get(index Value) (Value, Error) {
	idx, err := parseIndex(index, len(tp.array))
	if err != nil {
		return nil, err
	}
	return tp.array[idx.IntVal()], nil
}

func (tp *tuple) Len() (Int, Error) {
	return MakeInt(int64(len(tp.array))), nil
}
