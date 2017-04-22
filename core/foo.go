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
	"golang.org/x/exp/utf8string"
)

type foo struct {
	ustr *utf8string.String
}

func (s *foo) StrVal() string {
	return s.ustr.String()
}

//func (s *foo) Runes() []rune {
//	return s
//}

func MakeFoo(s string) *foo {
	return &foo{utf8string.NewString(s)}
}

func (s *foo) basicMarker() {}

func (s *foo) TypeOf() (Type, Error) { return TSTR, nil }

func (s *foo) String() (Str, Error) { return s, nil }

func (s *foo) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *foo:
		a := s.ustr.String()
		b := t.ustr.String()
		return MakeBool(a == b), nil

	default:
		return FALSE, nil
	}
}

func (s *foo) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {
	case *foo:
		a := s.ustr.String()
		b := t.ustr.String()
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
	return foocat(s, v)
}

func (s *foo) Get(index Value) (Value, Error) {
	idx, err := parseIndex(index, s.ustr.RuneCount())
	if err != nil {
		return nil, err
	}

	n := int(idx.IntVal())
	result := string([]rune{s.ustr.At(n)})
	return MakeFoo(result), nil
}

func (s *foo) Len() (Int, Error) {
	return MakeInt(int64(s.ustr.RuneCount())), nil
}

func (s *foo) Slice(from Value, to Value) (Value, Error) {

	f, err := parseIndex(from, s.ustr.RuneCount())
	if err != nil {
		return nil, err
	}

	t, err := parseIndex(to, s.ustr.RuneCount()+1)
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
	a := []byte(s.ustr.Slice(fn, tn))
	b := make([]byte, len(a))
	copy(b, a)
	return MakeFoo(string(b)), nil
}

func (s *foo) SliceFrom(from Value) (Value, Error) {
	return s.Slice(from, MakeInt(int64(s.ustr.RuneCount())))
}

func (s *foo) SliceTo(to Value) (Value, Error) {
	return s.Slice(ZERO, to)
}
