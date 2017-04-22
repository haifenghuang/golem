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
// list

type list struct {
	array []Value
}

func NewList(values []Value) List {
	return &list{values}
}

func (ls *list) compositeMarker() {}

func (ls *list) TypeOf() (Type, Error) {
	return TLIST, nil
}

func (ls *list) String() (Str, Error) {

	if len(ls.array) == 0 {
		return MakeStr("[]"), nil
	}

	var buf bytes.Buffer
	buf.WriteString("[")
	for idx, v := range ls.array {
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(" ")
		s, err := v.String()
		if err != nil {
			return nil, err
		}
		buf.WriteString(s.StrVal())
	}
	buf.WriteString(" ]")
	return MakeStr(buf.String()), nil
}

func (ls *list) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *list:
		return MakeBool(reflect.DeepEqual(ls.array, t.array)), nil
	default:
		return FALSE, nil
	}
}

func (ls *list) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (ls *list) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(ls, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (ls *list) Get(index Value) (Value, Error) {
	if i, ok := index.(Int); ok {
		n := int(i.IntVal())
		if (n < 0) || (n >= len(ls.array)) {
			return nil, IndexOutOfBoundsError()
		} else {
			return ls.array[n], nil
		}
	} else {
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func (ls *list) Set(index Value, val Value) Error {
	if i, ok := index.(Int); ok {
		n := int(i.IntVal())
		if (n < 0) || (n >= len(ls.array)) {
			return IndexOutOfBoundsError()
		} else {
			ls.array[n] = val
			return nil
		}
	} else {
		return TypeMismatchError("Expected 'Int'")
	}
}

func (ls *list) Append(val Value) Error {
	ls.array = append(ls.array, val)
	return nil
}

func (ls *list) Len() (Int, Error) {
	return MakeInt(int64(len(ls.array))), nil
}

func (ls *list) Slice(from Value, to Value) (Value, Error) {

	// from
	if f, ok := from.(Int); ok {
		fn := int(f.IntVal())
		if (fn < 0) || (fn >= len(ls.array)) {
			return nil, IndexOutOfBoundsError()
		} else {

			// to
			if t, ok := to.(Int); ok {
				tn := int(t.IntVal())

				if (tn < 0) || (tn > len(ls.array)) {
					return nil, IndexOutOfBoundsError()
				} else if tn < fn {
					// TODO do we want a different error here?
					return nil, IndexOutOfBoundsError()
				} else {

					a := ls.array[fn:tn]
					b := make([]Value, len(a))
					copy(b, a)
					return NewList(b), nil

				}
			} else {
				return nil, TypeMismatchError("Expected 'Int'")
			}

		}
	} else {
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func (ls *list) SliceFrom(from Value) (Value, Error) {
	if f, ok := from.(Int); ok {
		fn := int(f.IntVal())
		if (fn < 0) || (fn >= len(ls.array)) {
			return nil, IndexOutOfBoundsError()
		} else {
			a := ls.array[fn:]
			b := make([]Value, len(a))
			copy(b, a)
			return NewList(b), nil
		}
	} else {
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func (ls *list) SliceTo(to Value) (Value, Error) {
	if t, ok := to.(Int); ok {
		tn := int(t.IntVal())
		if (tn < 0) || (tn > len(ls.array)) {
			return nil, IndexOutOfBoundsError()
		} else {
			a := ls.array[:tn]
			b := make([]Value, len(a))
			copy(b, a)
			return NewList(b), nil
		}
	} else {
		return nil, TypeMismatchError("Expected 'Int'")
	}
}
