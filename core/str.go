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
//"fmt"
)

type str []rune

func (s str) String() string {
	return string(s)
}

func MakeStr(str string) Str {
	return fromString(str)
}

func (s str) basicMarker() {}

func (s str) TypeOf() (Type, Error) { return TSTR, nil }

func (s str) ToStr() (Str, Error) { return s, nil }

func (s str) HashCode() (Int, Error) {

	// https://en.wikipedia.org/wiki/Jenkins_hash_function
	var hash rune = 0
	for _, r := range s {
		hash += r
		hash += hash << 10
		hash ^= hash >> 6
	}
	hash += hash << 3
	hash ^= hash >> 11
	hash += hash << 15
	return MakeInt(int64(hash)), nil
}

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
	return Strcat(s, v)
}

func (s str) Get(index Value) (Value, Error) {
	idx, err := ParseIndex(index, len(s))
	if err != nil {
		return nil, err
	}

	return str([]rune{s[idx.IntVal()]}), nil
}

func (s str) Len() (Int, Error) {
	return MakeInt(int64(len(s))), nil
}

func (s str) Slice(from Value, to Value) (Value, Error) {

	f, err := ParseIndex(from, len(s))
	if err != nil {
		return nil, err
	}

	t, err := ParseIndex(to, len(s)+1)
	if err != nil {
		return nil, err
	}

	// TODO do we want a different error here?
	if t.IntVal() < f.IntVal() {
		return nil, IndexOutOfBoundsError()
	}

	// copy to avoid memory leaks
	a := s[f.IntVal():t.IntVal()]
	b := make([]rune, len(a))
	copy(b, a)
	return str(b), nil
}

func (s str) SliceFrom(from Value) (Value, Error) {
	return s.Slice(from, MakeInt(int64(len(s))))
}

func (s str) SliceTo(to Value) (Value, Error) {
	return s.Slice(ZERO, to)
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
		s, err := v.ToStr()
		if err != nil {
			return nil, err
		}
		return fromString(s.String()), nil
	}
}

func Strcat(a Value, b Value) (str, Error) {

	sa, err := fromValue(a)
	if err != nil {
		return nil, err
	}

	sb, err := fromValue(b)
	if err != nil {
		return nil, err
	}

	// copy to avoid memory leaks
	ca := make([]rune, len(sa))
	copy(ca, sa)

	cb := make([]rune, len(sb))
	copy(cb, sb)

	result := make(str, 0, len(ca)+len(cb))
	result = append(result, ca...)
	result = append(result, cb...)
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
