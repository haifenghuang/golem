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
)

type _str string

func (s _str) StrVal() string {
	return string(s)
}

func MakeStr(str string) Str {
	return _str(str)
}

func (s _str) TypeOf() (Type, Error) { return TSTR, nil }

func (s _str) String() (Str, Error) { return s, nil }

func (s _str) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {

	case _str:
		return MakeBool(s == t), nil

	default:
		return FALSE, nil
	}
}

func (s _str) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {

	case _str:
		if s < t {
			return NEG_ONE, nil
		} else if s > t {
			return ONE, nil
		} else {
			return ZERO, nil
		}

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (s _str) Add(v Value) (Value, Error) {
	return strcat([]Value{s, v})
}

func strcat(a []Value) (_str, Error) {
	var buf bytes.Buffer
	for _, v := range a {
		s, err := v.String()
		if err != nil {
			return _str(""), err
		}
		buf.WriteString(s.StrVal())
	}
	return _str(buf.String()), nil
}

func (s _str) GetField(key string) (Value, Error)   { return nil, TypeMismatchError("Expected 'Obj'") }
func (s _str) PutField(key string, val Value) Error { return TypeMismatchError("Expected 'Obj'") }
