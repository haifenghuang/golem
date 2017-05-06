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

func TestObj(t *testing.T) {
	o := NewObj([]*ObjEntry{})
	okType(t, o, TOBJ)

	s, err := o.ToStr()
	ok(t, s, err, MakeStr("obj {}"))

	z, err := o.Eq(NewObj([]*ObjEntry{}))
	ok(t, z, err, TRUE)
	z, err = o.Eq(NewObj([]*ObjEntry{&ObjEntry{"a", ONE}}))
	ok(t, z, err, FALSE)

	val, err := o.Add(MakeStr("a"))
	ok(t, val, err, MakeStr("obj {}a"))

	val, err = o.GetField(MakeStr("a"))
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	val, err = o.Get(MakeStr("a"))
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	val, err = o.Get(ZERO)
	fail(t, val, err, "TypeMismatch: Expected 'Str'")

	//////////////////

	o = NewObj([]*ObjEntry{&ObjEntry{"a", ONE}})
	okType(t, o, TOBJ)

	s, err = o.ToStr()
	ok(t, s, err, MakeStr("obj { a: 1 }"))

	z, err = o.Eq(NewObj([]*ObjEntry{}))
	ok(t, z, err, FALSE)
	z, err = o.Eq(NewObj([]*ObjEntry{&ObjEntry{"a", ONE}}))
	ok(t, z, err, TRUE)

	val, err = o.Add(MakeStr("a"))
	ok(t, val, err, MakeStr("obj { a: 1 }a"))

	val, err = o.GetField(MakeStr("a"))
	ok(t, val, err, ONE)

	val, err = o.GetField(MakeStr("b"))
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	val, err = o.Get(MakeStr("a"))
	ok(t, val, err, ONE)

	val, err = o.Get(MakeStr("b"))
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	err = o.PutField(MakeStr("a"), MakeInt(123))
	if err != nil {
		panic("unexpected error")
	}

	val, err = o.GetField(MakeStr("a"))
	ok(t, val, err, MakeInt(123))

	val, err = o.Get(MakeStr("a"))
	ok(t, val, err, MakeInt(123))

	err = o.Set(MakeStr("a"), MakeInt(456))
	if err != nil {
		panic("unexpected error")
	}

	val, err = o.GetField(MakeStr("a"))
	ok(t, val, err, MakeInt(456))

	val, err = o.Get(MakeStr("a"))
	ok(t, val, err, MakeInt(456))

	val, err = o.Has(MakeStr("a"))
	ok(t, val, err, TRUE)

	val, err = o.Has(MakeStr("abc"))
	ok(t, val, err, FALSE)

	val, err = o.Has(ZERO)
	fail(t, val, err, "TypeMismatch: Expected 'Str'")

	o = BlankObj([]string{"a"})
	val, err = o.GetField(MakeStr("a"))
	ok(t, val, err, NULL)
}

func TestList(t *testing.T) {
	ls := NewList([]Value{})
	okType(t, ls, TLIST)

	var v Value
	v, err := ls.ToStr()
	ok(t, v, err, MakeStr("[]"))

	v, err = ls.Eq(NewList([]Value{}))
	ok(t, v, err, TRUE)

	v, err = ls.Eq(NewList([]Value{MakeStr("a")}))
	ok(t, v, err, FALSE)

	v, err = ls.Eq(NULL)
	ok(t, v, err, FALSE)

	v, err = ls.Len()
	ok(t, v, err, ZERO)

	err = ls.Append(MakeStr("a"))
	assert(t, err == nil)

	v, err = ls.Eq(NewList([]Value{}))
	ok(t, v, err, FALSE)

	v, err = ls.Eq(NewList([]Value{MakeStr("a")}))
	ok(t, v, err, TRUE)

	v, err = ls.Len()
	ok(t, v, err, ONE)

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

	v, err = ls.ToStr()
	ok(t, v, err, MakeStr("[ b ]"))

	err = ls.Append(MakeStr("z"))
	assert(t, err == nil)

	v, err = ls.ToStr()
	ok(t, v, err, MakeStr("[ b, z ]"))

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
	h, err := NewDict(NewHashMap([]*HEntry{})).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = NewList([]Value{}).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = NewObj([]*ObjEntry{}).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")
}

func TestDict(t *testing.T) {
	d := NewDict(NewHashMap([]*HEntry{}))
	okType(t, d, TDICT)

	var v Value
	v, err := d.ToStr()
	ok(t, v, err, MakeStr("dict {}"))

	v, err = d.Eq(NewDict(NewHashMap([]*HEntry{})))
	ok(t, v, err, TRUE)

	v, err = d.Eq(NULL)
	ok(t, v, err, FALSE)

	v, err = d.Len()
	ok(t, v, err, ZERO)

	v, err = d.Get(MakeStr("a"))
	ok(t, v, err, NULL)

	err = d.Set(MakeStr("a"), ONE)
	assert(t, err == nil)

	v, err = d.Get(MakeStr("a"))
	ok(t, v, err, ONE)

	v, err = d.Eq(NewDict(NewHashMap([]*HEntry{})))
	ok(t, v, err, FALSE)

	v, err = d.Eq(NewDict(NewHashMap([]*HEntry{
		&HEntry{MakeStr("a"), ONE}})))
	ok(t, v, err, TRUE)

	v, err = d.Len()
	ok(t, v, err, ONE)

	v, err = d.ToStr()
	ok(t, v, err, MakeStr("dict { a: 1 }"))

	err = d.Set(MakeStr("b"), MakeInt(2))
	assert(t, err == nil)

	v, err = d.Get(MakeStr("b"))
	ok(t, v, err, MakeInt(2))

	v, err = d.ToStr()
	ok(t, v, err, MakeStr("dict { b: 2, a: 1 }"))

	tp := NewTuple([]Value{ONE, ZERO})
	d = NewDict(NewHashMap([]*HEntry{
		&HEntry{tp, TRUE}}))

	v, err = d.ToStr()
	ok(t, v, err, MakeStr("dict { (1, 0): true }"))

	v, err = d.Get(tp)
	ok(t, v, err, TRUE)
}

func TestTuple(t *testing.T) {
	var v Value

	tp := NewTuple([]Value{ONE, ZERO})
	okType(t, tp, TTUPLE)

	v, err := tp.Eq(NewTuple([]Value{ZERO, ZERO}))
	ok(t, v, err, FALSE)

	v, err = tp.Eq(NewTuple([]Value{ONE, ZERO}))
	ok(t, v, err, TRUE)

	v, err = tp.Eq(NULL)
	ok(t, v, err, FALSE)

	v, err = tp.Get(ZERO)
	ok(t, v, err, ONE)

	v, err = tp.Get(ONE)
	ok(t, v, err, ZERO)

	v, err = tp.Get(NEG_ONE)
	fail(t, v, err, "IndexOutOfBounds")

	v, err = tp.Get(MakeInt(2))
	fail(t, v, err, "IndexOutOfBounds")

	v, err = tp.ToStr()
	ok(t, v, err, MakeStr("(1, 0)"))

	v, err = tp.Len()
	ok(t, v, err, MakeInt(2))
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

	r := newRange(0, 5, 1)
	okType(t, r, TRANGE)

	v, err := r.Eq(newRange(0, 5, 2))
	ok(t, v, err, FALSE)

	v, err = r.Eq(newRange(0, 5, 1))
	ok(t, v, err, TRUE)

	v, err = r.Eq(NULL)
	ok(t, v, err, FALSE)

	v, err = r.Len()
	ok(t, v, err, MakeInt(5))

	v, err = newRange(0, 6, 3).Len()
	ok(t, v, err, MakeInt(2))
	v, err = newRange(0, 7, 3).Len()
	ok(t, v, err, MakeInt(2))
	v, err = newRange(0, 8, 3).Len()
	ok(t, v, err, MakeInt(2))
	v, err = newRange(0, 9, 3).Len()
	ok(t, v, err, MakeInt(3))

	v, err = newRange(0, 0, 3).Len()
	ok(t, v, err, MakeInt(0))
	v, err = newRange(1, 0, 1).Len()
	ok(t, v, err, MakeInt(0))

	v, err = NewRange(1, 0, 0)
	fail(t, v, err, "InvalidArgument: step cannot be 0")

	v, err = newRange(0, -5, -1).Len()
	ok(t, v, err, MakeInt(5))
	v, err = newRange(-1, -8, -3).Len()
	ok(t, v, err, MakeInt(2))

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
	for objInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := objInvokeFunc(t, itr, MakeStr("getValue"))

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
	for objInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := objInvokeFunc(t, itr, MakeStr("getValue"))

		i, ok := v.(Int)
		assert(t, ok)
		n *= i.IntVal()
	}
	assert(t, n == 24)
}

func TestDictIterator(t *testing.T) {

	var ibl Iterable = NewDict(
		NewHashMap([]*HEntry{
			&HEntry{MakeStr("a"), ONE},
			&HEntry{MakeStr("b"), MakeInt(2)},
			&HEntry{MakeStr("c"), MakeInt(3)}}))

	var itr Iterator = ibl.NewIterator()
	v, err := itr.IterGet()
	fail(t, v, err, "NoSuchElement")
	s := MakeStr("")
	for itr.IterNext().BoolVal() {
		v, err = itr.IterGet()
		assert(t, err == nil)

		tp, ok := v.(Tuple)
		assert(t, ok)
		s, err = Strcat(s, tp)
		assert(t, err == nil)
	}
	ok(t, s, nil, MakeStr("(b, 2)(a, 1)(c, 3)"))
	v, err = itr.IterGet()
	fail(t, v, err, "NoSuchElement")

	itr = ibl.NewIterator()
	s = MakeStr("")
	for objInvokeBoolFunc(t, itr, MakeStr("nextValue")).BoolVal() {
		v := objInvokeFunc(t, itr, MakeStr("getValue"))

		tp, ok := v.(Tuple)
		assert(t, ok)
		s, err = Strcat(s, tp)
		assert(t, err == nil)
	}
	ok(t, s, nil, MakeStr("(b, 2)(a, 1)(c, 3)"))
}
