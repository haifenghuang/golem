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

type str struct {
	ustr *utf8string.String
}

func (s *str) StrVal() string {
	return s.ustr.String()
}

func MakeStr(s string) *str {
	return &str{utf8string.NewString(s)}
}

func (s *str) basicMarker() {}

func (s *str) TypeOf() (Type, Error) { return TSTR, nil }

func (s *str) String() (Str, Error) { return s, nil }

func (s *str) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *str:
		a := s.ustr.String()
		b := t.ustr.String()
		return MakeBool(a == b), nil

	default:
		return FALSE, nil
	}
}

func (s *str) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {
	case *str:
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

func (s *str) Add(v Value) (Value, Error) {
	return strcat(s, v)
}

func (s *str) Get(index Value) (Value, Error) {
	idx, err := parseIndex(index, s.ustr.RuneCount())
	if err != nil {
		return nil, err
	}

	// NOTE: We are copying the slice.
	// This will allow garbage collection of strings
	// that are no longer referenced directly.
	n := int(idx.IntVal())
	result := string([]rune{s.ustr.At(n)})
	return MakeStr(result), nil
}

func (s *str) Len() (Int, Error) {
	return MakeInt(int64(s.ustr.RuneCount())), nil
}

func (s *str) Slice(from Value, to Value) (Value, Error) {

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
	return MakeStr(string(b)), nil
}

func (s *str) SliceFrom(from Value) (Value, Error) {
	return s.Slice(from, MakeInt(int64(s.ustr.RuneCount())))
}

func (s *str) SliceTo(to Value) (Value, Error) {
	return s.Slice(ZERO, to)
}
