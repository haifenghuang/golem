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

func (ls *list) ToStr() (Str, Error) {

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
		s, err := v.ToStr()
		if err != nil {
			return nil, err
		}
		buf.WriteString(s.String())
	}
	buf.WriteString(" ]")
	return MakeStr(buf.String()), nil
}

func (ls *list) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
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
		return Strcat(ls, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (ls *list) Get(index Value) (Value, Error) {
	idx, err := ParseIndex(index, len(ls.array))
	if err != nil {
		return nil, err
	}
	return ls.array[idx.IntVal()], nil
}

func (ls *list) Set(index Value, val Value) Error {
	idx, err := ParseIndex(index, len(ls.array))
	if err != nil {
		return err
	}

	ls.array[idx.IntVal()] = val
	return nil
}

func (ls *list) Append(val Value) Error {
	ls.array = append(ls.array, val)
	return nil
}

func (ls *list) Len() (Int, Error) {
	return MakeInt(int64(len(ls.array))), nil
}

func (ls *list) Slice(from Value, to Value) (Value, Error) {

	f, err := ParseIndex(from, len(ls.array))
	if err != nil {
		return nil, err
	}

	t, err := ParseIndex(to, len(ls.array)+1)
	if err != nil {
		return nil, err
	}

	// TODO do we want a different error here?
	if t.IntVal() < f.IntVal() {
		return nil, IndexOutOfBoundsError()
	}

	a := ls.array[f.IntVal():t.IntVal()]
	b := make([]Value, len(a))
	copy(b, a)
	return NewList(b), nil
}

func (ls *list) SliceFrom(from Value) (Value, Error) {
	return ls.Slice(from, MakeInt(int64(len(ls.array))))
}

func (ls *list) SliceTo(to Value) (Value, Error) {
	return ls.Slice(ZERO, to)
}

//---------------------------------------------------------------
// Iterator

type listIterator struct {
	Obj
	ls *list
	n  int
}

func (ls *list) NewIterator() Iterator {

	// TODO make this immutable
	obj := NewObj()

	iter := &listIterator{obj, ls, -1}

	obj.Init(
		&ObjDef{[]string{
			"nextValue",
			"getValue"}},
		[]Value{
			&nativeIterNext{&nativeFunc{}, iter},
			&nativeIterGet{&nativeFunc{}, iter}})

	return iter
}

func (i *listIterator) IterNext() Bool {

	i.n++
	return MakeBool(i.n < len(i.ls.array))
}

func (i *listIterator) IterGet() (Value, Error) {

	if (i.n >= 0) && (i.n < len(i.ls.array)) {
		return i.ls.array[i.n], nil
	} else {
		return nil, NoSuchElementError()
	}
}
