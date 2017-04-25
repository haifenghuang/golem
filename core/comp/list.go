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

package comp

import (
	"bytes"
	g "golem/core"
	"reflect"
)

//---------------------------------------------------------------
// list

type list struct {
	array []g.Value
}

func NewList(values []g.Value) List {
	return &list{values}
}

func (ls *list) compositeMarker() {}

func (ls *list) TypeOf() (g.Type, g.Error) {
	return g.TLIST, nil
}

func (ls *list) ToStr() (g.Str, g.Error) {

	if len(ls.array) == 0 {
		return g.MakeStr("[]"), nil
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
	return g.MakeStr(buf.String()), nil
}

func (ls *list) HashCode() (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Hashable Type")
}

func (ls *list) Eq(v g.Value) (g.Bool, g.Error) {
	switch t := v.(type) {
	case *list:
		return g.MakeBool(reflect.DeepEqual(ls.array, t.array)), nil
	default:
		return g.FALSE, nil
	}
}

func (ls *list) Cmp(v g.Value) (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Comparable Type")
}

func (ls *list) Add(v g.Value) (g.Value, g.Error) {
	switch t := v.(type) {

	case g.Str:
		return g.Strcat(ls, t)

	default:
		return nil, g.TypeMismatchError("Expected Number Type")
	}
}

func (ls *list) Get(index g.Value) (g.Value, g.Error) {
	idx, err := g.ParseIndex(index, len(ls.array))
	if err != nil {
		return nil, err
	}
	return ls.array[idx.IntVal()], nil
}

func (ls *list) Set(index g.Value, val g.Value) g.Error {
	idx, err := g.ParseIndex(index, len(ls.array))
	if err != nil {
		return err
	}

	ls.array[idx.IntVal()] = val
	return nil
}

func (ls *list) Append(val g.Value) g.Error {
	ls.array = append(ls.array, val)
	return nil
}

func (ls *list) Len() (g.Int, g.Error) {
	return g.MakeInt(int64(len(ls.array))), nil
}

func (ls *list) Slice(from g.Value, to g.Value) (g.Value, g.Error) {

	f, err := g.ParseIndex(from, len(ls.array))
	if err != nil {
		return nil, err
	}

	t, err := g.ParseIndex(to, len(ls.array)+1)
	if err != nil {
		return nil, err
	}

	// TODO do we want a different error here?
	if t.IntVal() < f.IntVal() {
		return nil, g.IndexOutOfBoundsError()
	}

	a := ls.array[f.IntVal():t.IntVal()]
	b := make([]g.Value, len(a))
	copy(b, a)
	return NewList(b), nil
}

func (ls *list) SliceFrom(from g.Value) (g.Value, g.Error) {
	return ls.Slice(from, g.MakeInt(int64(len(ls.array))))
}

func (ls *list) SliceTo(to g.Value) (g.Value, g.Error) {
	return ls.Slice(g.ZERO, to)
}
