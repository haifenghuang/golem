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
//"bytes"
)

type str []rune

func (s str) StrVal() string {
	return string(s)
}

func (s str) Runes() []rune {
	return s
}

func MakeStr(str string) Str {
	return fromString(str)
}

func (s str) TypeOf() (Type, Error) { return TSTR, nil }

func (s str) String() (Str, Error) { return s, nil }

func (s str) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {

	case str:
		return MakeBool(runesEq(s, t)), nil

	default:
		return FALSE, nil
	}
}

func (s str) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {

	case str:
		return MakeInt(int64(runesCmp(s, t))), nil

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (s str) Add(v Value) (Value, Error) {
	return strcat(s, v)
}

func (s str) Get(index Value) (Value, Error) {
	if i, ok := index.(Int); ok {
		n := int(i.IntVal())
		if (n < 0) || (n >= len(s)) {
			return nil, IndexOutOfBoundsError()
		} else {
			result := make([]rune, 1)
			result[0] = s[n]
			return str(result), nil
		}
	} else {
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

//--------------------------------------------------------------

func fromString(s string) str {
	z := str{}
	for _, r := range s {
		z = append(z, r)
	}
	return z
}

func fromValue(v Value) (str, Error) {
	if sv, ok := v.(str); ok {
		return sv, nil
	} else {
		s, err := v.String()
		if err != nil {
			return nil, err
		}
		return fromString(s.StrVal()), nil
	}
}

func strcat(a Value, b Value) (str, Error) {
	result := str{}

	s, err := fromValue(a)
	if err != nil {
		return nil, err
	}
	result = append(result, s...)

	s, err = fromValue(b)
	if err != nil {
		return nil, err
	}
	result = append(result, s...)

	return result, nil
}

func runesEq(a str, b str) bool {
	if len(a) != len(b) {
		return false
	}
	for i, r := range a {
		if r != b[i] {
			return false
		}
	}
	return true
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func runesCmp(a str, b str) int {
	n := min(len(a), len(b))
	for i := 0; i < n; i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return 1
		}
	}
	return len(a) - len(b)
}
