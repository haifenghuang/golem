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
	//"golang.org/x/exp/utf8string"
	"unicode/utf8"
)

type foo struct {
	contents string
	isAscii  bool
}

func (s *foo) StrVal() string {
	return s.contents
}

func MakeFoo(s string) *foo {

	isAscii := true
	for i := 0; i < len(s); i++ {
		if s[i] >= utf8.RuneSelf {
			isAscii = false
			break
		}
	}

	return &foo{s, isAscii}
}

func (s *foo) basicMarker() {}

func (s *foo) TypeOf() (Type, Error) { return TSTR, nil }

func (s *foo) String() (Str, Error) { return s, nil }

func (s *foo) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *foo:
		a := s.contents
		b := t.contents
		return MakeBool(a == b), nil

	default:
		return FALSE, nil
	}
}

func (s *foo) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {
	case *foo:
		a := s.contents
		b := t.contents
		if a < b {
			return NEG_ONE, nil
		} else if a > b {
			return ONE, nil
		} else {
			return ZERO, nil
		}

	default:
		return nil, TypeMismatchError("Expected Comparable Type")
	}
}

func (s *foo) Add(v Value) (Value, Error) {
	return strcat(s, v)
}

func (s *foo) Len() (Int, Error) {
	return MakeInt(int64(len(s.contents))), nil
}

func (s *foo) Get(index Value) (Value, Error) {
	idx, err := parseIndex(index, len(s.contents))
	if err != nil {
		return nil, err
	}
	n := int(idx.IntVal())

	// NOTE: We are copying the slice.
	// This will allow garbage collection of strings
	// that are no longer referenced directly.
	if s.isAscii {
		z := s.contents[n]
		return MakeFoo(string([]byte{z})), nil
	} else {
		// TODO we need to figure out how to do this efficiently...
		runes := toRunes(s.contents)
		z := runes[n]
		return MakeFoo(string([]rune{z})), nil
	}
}

//a := []byte(s.contents.Slice(fn, tn))
//b := make([]byte, len(a))
//copy(b, a)
//return MakeFoo(string(b)), nil

func (s *foo) Slice(from Value, to Value) (Value, Error) {

	f, err := parseIndex(from, len(s.contents))
	if err != nil {
		return nil, err
	}

	t, err := parseIndex(to, len(s.contents)+1)
	if err != nil {
		return nil, err
	}

	fn := int(f.IntVal())
	tn := int(t.IntVal())

	// TODO do we want a different error here?
	if tn < fn {
		return nil, IndexOutOfBoundsError()
	}

	// NOTE: We are copying the slice.
	// This will allow garbage collection of strings
	// that are no longer referenced directly.
	if s.isAscii {
		a := []byte(s.contents[fn:tn])
		b := make([]byte, len(a))
		copy(b, a)
		return MakeFoo(string(b)), nil
	} else {
		// TODO we need to figure out how to do this efficiently...
		runes := toRunes(s.contents)
		a := runes[fn:tn]
		b := make([]rune, len(a))
		copy(b, a)
		return MakeFoo(string(b)), nil
	}
}

func (s *foo) SliceFrom(from Value) (Value, Error) {
	return s.Slice(from, MakeInt(int64(len(s.contents))))
}

func (s *foo) SliceTo(to Value) (Value, Error) {
	return s.Slice(ZERO, to)
}

func toRunes(s string) []rune {
	runes := []rune{}
	for _, r := range s {
		runes = append(runes, r)
	}
	return runes
}
