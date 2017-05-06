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
	return toRunes(str)
}

func (s str) basicMarker() {}

func (s str) TypeOf() Type { return TSTR }

func (s str) ToStr() Str { return s }

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

func (s str) Eq(v Value) Bool {
	switch t := v.(type) {

	case str:
		return MakeBool(runesEq(s, t))

	default:
		return FALSE
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

func (s str) Len() Int {
	return MakeInt(int64(len(s)))
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

func toRunes(s string) str {
	z := str{}
	for _, r := range s {
		z = append(z, r)
	}
	return z
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

//---------------------------------------------------------------
// Iterator

type strIterator struct {
	Obj
	s str
	n int
}

func (s str) NewIterator() Iterator {

	next := &nativeIterNext{&nativeFunc{}, nil}
	get := &nativeIterGet{&nativeFunc{}, nil}
	// TODO make this immutable
	obj := NewObj([]*ObjEntry{
		&ObjEntry{"nextValue", next},
		&ObjEntry{"getValue", get}})

	itr := &strIterator{obj, s, -1}

	next.itr = itr
	get.itr = itr
	return itr
}

func (i *strIterator) IterNext() Bool {
	i.n++
	return MakeBool(i.n < len(i.s))
}

func (i *strIterator) IterGet() (Value, Error) {

	if (i.n >= 0) && (i.n < len(i.s)) {
		return str([]rune{i.s[i.n]}), nil
	} else {
		return nil, NoSuchElementError()
	}
}
