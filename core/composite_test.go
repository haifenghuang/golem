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
	"testing"
)

func TestStruct(t *testing.T) {
	stc := NewStruct([]*StructEntry{})
	okType(t, stc, TSTRUCT)

	s := stc.ToStr()
	ok(t, s, nil, MakeStr("struct { }"))

	z := stc.Eq(NewStruct([]*StructEntry{}))
	ok(t, z, nil, TRUE)
	z = stc.Eq(NewStruct([]*StructEntry{{"a", ONE}}))
	ok(t, z, nil, FALSE)

	val, err := stc.Plus(MakeStr("a"))
	ok(t, val, err, MakeStr("struct { }a"))

	val, err = stc.GetField(MakeStr("a"))
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	val, err = stc.Get(MakeStr("a"))
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	val, err = stc.Get(ZERO)
	fail(t, val, err, "TypeMismatch: Expected 'Str'")

	//////////////////

	stc = NewStruct([]*StructEntry{{"a", ONE}})
	okType(t, stc, TSTRUCT)

	s = stc.ToStr()
	ok(t, s, nil, MakeStr("struct { a: 1 }"))

	z = stc.Eq(NewStruct([]*StructEntry{}))
	ok(t, z, nil, FALSE)
	z = stc.Eq(NewStruct([]*StructEntry{{"a", ONE}}))
	ok(t, z, nil, TRUE)

	val, err = stc.Plus(MakeStr("a"))
	ok(t, val, err, MakeStr("struct { a: 1 }a"))

	val, err = stc.GetField(MakeStr("a"))
	ok(t, val, err, ONE)

	val, err = stc.GetField(MakeStr("b"))
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	val, err = stc.Get(MakeStr("a"))
	ok(t, val, err, ONE)

	val, err = stc.Get(MakeStr("b"))
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	err = stc.PutField(MakeStr("a"), MakeInt(123))
	if err != nil {
		panic("unexpected error")
	}

	val, err = stc.GetField(MakeStr("a"))
	ok(t, val, err, MakeInt(123))

	val, err = stc.Get(MakeStr("a"))
	ok(t, val, err, MakeInt(123))

	err = stc.Set(MakeStr("a"), MakeInt(456))
	if err != nil {
		panic("unexpected error")
	}

	val, err = stc.GetField(MakeStr("a"))
	ok(t, val, err, MakeInt(456))

	val, err = stc.Get(MakeStr("a"))
	ok(t, val, err, MakeInt(456))

	val, err = stc.Has(MakeStr("a"))
	ok(t, val, err, TRUE)

	val, err = stc.Has(MakeStr("abc"))
	ok(t, val, err, FALSE)

	val, err = stc.Has(ZERO)
	fail(t, val, err, "TypeMismatch: Expected 'Str'")

	stc = BlankStruct([]string{"a"})
	val, err = stc.GetField(MakeStr("a"))
	ok(t, val, err, NULL)
}

func TestList(t *testing.T) {
	ls := NewList([]Value{})
	okType(t, ls, TLIST)

	var v Value
	var err Error

	v = ls.ToStr()
	ok(t, v, nil, MakeStr("[ ]"))

	v = ls.Eq(NewList([]Value{}))
	ok(t, v, nil, TRUE)

	v = ls.Eq(NewList([]Value{MakeStr("a")}))
	ok(t, v, nil, FALSE)

	v = ls.Eq(NULL)
	ok(t, v, nil, FALSE)

	v = ls.Len()
	ok(t, v, nil, ZERO)

	ls.Add(MakeStr("a"))

	v = ls.Eq(NewList([]Value{}))
	ok(t, v, nil, FALSE)

	v = ls.Eq(NewList([]Value{MakeStr("a")}))
	ok(t, v, nil, TRUE)

	v = ls.Len()
	ok(t, v, nil, ONE)

	v, err = ls.Get(ZERO)
	ok(t, v, err, MakeStr("a"))

	err = ls.Set(ZERO, MakeStr("b"))
	assert(t, err == nil)

	v, err = ls.Get(ZERO)
	ok(t, v, err, MakeStr("b"))

	v, err = ls.Get(NEG_ONE)
	fail(t, v, err, "IndexOutOfBounds")

	v, err = ls.Get(ONE)
	fail(t, v, err, "IndexOutOfBounds")

	err = ls.Set(NEG_ONE, TRUE)
	fail(t, nil, err, "IndexOutOfBounds")

	err = ls.Set(ONE, TRUE)
	fail(t, nil, err, "IndexOutOfBounds")

	v = ls.ToStr()
	ok(t, v, nil, MakeStr("[ b ]"))

	ls.Add(MakeStr("z"))

	v = ls.ToStr()
	ok(t, v, nil, MakeStr("[ b, z ]"))

	//////////////////////////////
	// sliceable

	ls = NewList([]Value{})
	v, err = ls.SliceFrom(ZERO)
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.SliceTo(ZERO)
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.SliceTo(ONE)
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.Slice(ZERO, ONE)
	fail(t, nil, err, "IndexOutOfBounds")

	ls = NewList([]Value{TRUE, FALSE, NULL})
	v, err = ls.SliceFrom(ONE)
	ok(t, v, err, NewList([]Value{FALSE, NULL}))
	v, err = ls.SliceTo(ONE)
	ok(t, v, err, NewList([]Value{TRUE}))
	v, err = ls.Slice(ZERO, ONE)
	ok(t, v, err, NewList([]Value{TRUE}))
	v, err = ls.Slice(ZERO, MakeInt(3))
	ok(t, v, err, NewList([]Value{TRUE, FALSE, NULL}))

	v, err = ls.Slice(ZERO, ZERO)
	ok(t, v, err, NewList([]Value{}))

	v, err = ls.Slice(MakeInt(2), ZERO)
	fail(t, nil, err, "IndexOutOfBounds")

	v, err = ls.SliceFrom(MakeInt(7))
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.SliceTo(MakeInt(7))
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.Slice(MakeInt(7), MakeInt(7))
	fail(t, nil, err, "IndexOutOfBounds")
}

func TestCompositeHashCode(t *testing.T) {
	h, err := NewDict([]*HEntry{}).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = NewList([]Value{}).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = NewStruct([]*StructEntry{}).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")
}

func TestDict(t *testing.T) {
	d := NewDict([]*HEntry{})
	okType(t, d, TDICT)

	var v Value
	var err Error

	v = d.ToStr()
	ok(t, v, err, MakeStr("dict { }"))

	v = d.Eq(NewDict([]*HEntry{}))
	ok(t, v, nil, TRUE)

	v = d.Eq(NULL)
	ok(t, v, nil, FALSE)

	v = d.Len()
	ok(t, v, nil, ZERO)

	v, err = d.Get(MakeStr("a"))
	ok(t, v, err, NULL)

	err = d.Set(MakeStr("a"), ONE)
	assert(t, err == nil)

	v, err = d.Get(MakeStr("a"))
	ok(t, v, err, ONE)

	v = d.Eq(NewDict([]*HEntry{}))
	ok(t, v, nil, FALSE)

	v = d.Eq(NewDict([]*HEntry{{MakeStr("a"), ONE}}))
	ok(t, v, nil, TRUE)

	v = d.Len()
	ok(t, v, nil, ONE)

	v = d.ToStr()
	ok(t, v, nil, MakeStr("dict { a: 1 }"))

	err = d.Set(MakeStr("b"), MakeInt(2))
	assert(t, err == nil)

	v, err = d.Get(MakeStr("b"))
	ok(t, v, err, MakeInt(2))

	v = d.ToStr()
	ok(t, v, nil, MakeStr("dict { b: 2, a: 1 }"))

	tp := NewTuple([]Value{ONE, ZERO})
	d = NewDict([]*HEntry{{tp, TRUE}})

	v = d.ToStr()
	ok(t, v, nil, MakeStr("dict { (1, 0): true }"))

	v, err = d.Get(tp)
	ok(t, v, err, TRUE)
}

func TestSet(t *testing.T) {
	s := NewSet([]Value{})
	okType(t, s, TDICT)

	var v Value
	var err Error

	v = s.ToStr()
	ok(t, v, err, MakeStr("set { }"))

	v = s.Eq(NewSet([]Value{}))
	ok(t, v, nil, TRUE)

	v = s.Eq(NewSet([]Value{ONE}))
	ok(t, v, nil, FALSE)

	v = s.Eq(NULL)
	ok(t, v, nil, FALSE)

	v = s.Len()
	ok(t, v, nil, ZERO)

	s = NewSet([]Value{ONE})

	v = s.ToStr()
	ok(t, v, err, MakeStr("set { 1 }"))

	v = s.Eq(NewSet([]Value{}))
	ok(t, v, nil, FALSE)

	v = s.Eq(NewSet([]Value{ONE, ONE, ONE}))
	ok(t, v, nil, TRUE)

	v = s.Eq(NULL)
	ok(t, v, nil, FALSE)

	v = s.Len()
	ok(t, v, nil, ONE)

	s = NewSet([]Value{ONE, ZERO, ZERO, ONE})

	v = s.ToStr()
	ok(t, v, err, MakeStr("set { 0, 1 }"))

	v = s.Len()
	ok(t, v, nil, MakeInt(2))
}

func TestTuple(t *testing.T) {
	var v Value
	var err Error

	tp := NewTuple([]Value{ONE, ZERO})
	okType(t, tp, TTUPLE)

	v = tp.Eq(NewTuple([]Value{ZERO, ZERO}))
	ok(t, v, nil, FALSE)

	v = tp.Eq(NewTuple([]Value{ONE, ZERO}))
	ok(t, v, nil, TRUE)

	v = tp.Eq(NULL)
	ok(t, v, nil, FALSE)

	v, err = tp.Get(ZERO)
	ok(t, v, err, ONE)

	v, err = tp.Get(ONE)
	ok(t, v, err, ZERO)

	v, err = tp.Get(NEG_ONE)
	fail(t, v, err, "IndexOutOfBounds")

	v, err = tp.Get(MakeInt(2))
	fail(t, v, err, "IndexOutOfBounds")

	v = tp.ToStr()
	ok(t, v, nil, MakeStr("(1, 0)"))

	v = tp.Len()
	ok(t, v, nil, MakeInt(2))
}

func newRange(from int64, to int64, step int64) Range {
	r, err := NewRange(from, to, step)
	if err != nil {
		panic("invalid range")
	}
	return r
}

func TestRange(t *testing.T) {
	var v Value
	var err Error

	r := newRange(0, 5, 1)
	okType(t, r, TRANGE)

	v = r.Eq(newRange(0, 5, 2))
	ok(t, v, err, FALSE)

	v = r.Eq(newRange(0, 5, 1))
	ok(t, v, err, TRUE)

	v = r.Eq(NULL)
	ok(t, v, err, FALSE)

	v = r.Len()
	ok(t, v, nil, MakeInt(5))

	v = newRange(0, 6, 3).Len()
	ok(t, v, nil, MakeInt(2))
	v = newRange(0, 7, 3).Len()
	ok(t, v, nil, MakeInt(2))
	v = newRange(0, 8, 3).Len()
	ok(t, v, nil, MakeInt(2))
	v = newRange(0, 9, 3).Len()
	ok(t, v, nil, MakeInt(3))

	v = newRange(0, 0, 3).Len()
	ok(t, v, nil, MakeInt(0))
	v = newRange(1, 0, 1).Len()
	ok(t, v, nil, MakeInt(0))

	v, err = NewRange(1, 0, 0)
	fail(t, v, err, "InvalidArgument: step cannot be 0")

	v = newRange(0, -5, -1).Len()
	ok(t, v, nil, MakeInt(5))
	v = newRange(-1, -8, -3).Len()
	ok(t, v, nil, MakeInt(2))

	r = newRange(0, 5, 1)
	v, err = r.Get(ONE)
	ok(t, v, err, MakeInt(1))

	r = newRange(3, 9, 2)
	v, err = r.Get(MakeInt(2))
	ok(t, v, err, MakeInt(7))

	r = newRange(-9, -13, -1)
	v, err = r.Get(ONE)
	ok(t, v, err, MakeInt(-10))

	r = newRange(0, 5, 1)
	v, err = r.Slice(ONE, MakeInt(3))
	ok(t, v, err, newRange(1, 3, 1))
	v, err = r.SliceFrom(ONE)
	ok(t, v, err, newRange(1, 5, 1))
	v, err = r.SliceTo(MakeInt(3))
	ok(t, v, err, newRange(0, 3, 1))

	ok(t, r.From(), nil, ZERO)
	ok(t, r.To(), nil, MakeInt(5))
	ok(t, r.Step(), nil, ONE)
}

func TestRangeIterator(t *testing.T) {

	var ibl Iterable = newRange(1, 5, 1)

	var itr Iterator = ibl.NewIterator()
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	var n int64 = 1
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		assert(t, err == nil)

		i, ok := v.(Int)
		assert(t, ok)
		n *= i.IntVal()
	}
	assert(t, n == 24)
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator()
	n = 1
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))

		i, ok := v.(Int)
		assert(t, ok)
		n *= i.IntVal()
	}
	assert(t, n == 24)
}

func TestListIterator(t *testing.T) {

	var ibl Iterable = NewList(
		[]Value{MakeInt(1), MakeInt(2), MakeInt(3), MakeInt(4)})

	var itr Iterator = ibl.NewIterator()
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	var n int64 = 1
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		assert(t, err == nil)

		i, ok := v.(Int)
		assert(t, ok)
		n *= i.IntVal()
	}
	assert(t, n == 24)
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator()
	n = 1
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))

		i, ok := v.(Int)
		assert(t, ok)
		n *= i.IntVal()
	}
	assert(t, n == 24)
}

func TestDictIterator(t *testing.T) {

	var ibl Iterable = NewDict(
		[]*HEntry{
			{MakeStr("a"), ONE},
			{MakeStr("b"), MakeInt(2)},
			{MakeStr("c"), MakeInt(3)}})

	var itr Iterator = ibl.NewIterator()
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := MakeStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		assert(t, err == nil)

		tp, ok := v.(Tuple)
		assert(t, ok)
		s = strcat(s, tp)
	}
	ok(t, s, nil, MakeStr("(b, 2)(a, 1)(c, 3)"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator()
	s = MakeStr("")
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))

		tp, ok := v.(Tuple)
		assert(t, ok)
		s = strcat(s, tp)
	}
	ok(t, s, nil, MakeStr("(b, 2)(a, 1)(c, 3)"))
}

func TestSetIterator(t *testing.T) {

	var ibl Iterable = NewSet(
		[]Value{MakeStr("a"), MakeStr("b"), MakeStr("c")})

	var itr Iterator = ibl.NewIterator()
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := MakeStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		assert(t, err == nil)

		s = strcat(s, v)
	}
	ok(t, s, nil, MakeStr("bac"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator()
	s = MakeStr("")
	for structInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := structInvokeFunc(t, itr, MakeStr("getValue"))

		s = strcat(s, v)
	}
	ok(t, s, nil, MakeStr("bac"))
}
