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

func TestFunc(t *testing.T) {

	a := NewFunc(&Template{})
	b := NewFunc(&Template{})

	okType(t, a, TFUNC)

	z, err := a.Eq(b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(a)
	ok(t, z, err, TRUE)
}

func newObj(fields map[string]Value) Obj {
	o := NewObj()
	def := &ObjDef{[]string{}}
	values := []Value{}
	for k, v := range fields {
		def.Keys = append(def.Keys, k)
		values = append(values, v)
	}
	o.Init(def, values)
	return o
}

func TestObj(t *testing.T) {
	o := newObj(map[string]Value{})
	okType(t, o, TOBJ)

	s, err := o.ToStr()
	ok(t, s, err, MakeStr("obj {}"))

	z, err := o.Eq(newObj(map[string]Value{}))
	ok(t, z, err, TRUE)
	z, err = o.Eq(newObj(map[string]Value{"a": ONE}))
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

	o = newObj(map[string]Value{"a": ONE})
	okType(t, o, TOBJ)

	s, err = o.ToStr()
	ok(t, s, err, MakeStr("obj { a: 1 }"))

	z, err = o.Eq(newObj(map[string]Value{}))
	ok(t, z, err, FALSE)
	z, err = o.Eq(newObj(map[string]Value{"a": ONE}))
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
}

func uninitErr(t *testing.T, err Error) {
	if err.Error() != "UninitializedObj: Obj is not yet initialized" {
		t.Error("bad uninitialized error")
	}
}

func TestUninitialized(t *testing.T) {
	o := NewObj()
	_, e0 := o.TypeOf()
	_, e1 := o.Eq(NULL)
	_, e2 := o.ToStr()
	_, e3 := o.Cmp(NULL)
	_, e4 := o.Add(NULL)

	_, e5 := o.GetField(MakeStr(""))
	e6 := o.PutField(MakeStr(""), NULL)

	_, e7 := o.Get(MakeStr(""))
	e8 := o.Set(MakeStr(""), NULL)

	_, e9 := o.Has(NULL)
	_, e10 := o.HashCode()

	uninitErr(t, e0)
	uninitErr(t, e1)
	uninitErr(t, e2)
	uninitErr(t, e3)
	uninitErr(t, e4)
	uninitErr(t, e5)
	uninitErr(t, e6)
	uninitErr(t, e7)
	uninitErr(t, e8)
	uninitErr(t, e9)
	uninitErr(t, e10)
}

func TestLineNumber(t *testing.T) {

	tp := &Template{0, 0, 0, nil,
		[]OpcLine{
			OpcLine{0, 0},
			OpcLine{1, 2},
			OpcLine{11, 3},
			OpcLine{20, 4},
			OpcLine{29, 0}}}

	assert(t, tp.LineNumber(0) == 0)
	assert(t, tp.LineNumber(1) == 2)
	assert(t, tp.LineNumber(10) == 2)
	assert(t, tp.LineNumber(11) == 3)
	assert(t, tp.LineNumber(19) == 3)
	assert(t, tp.LineNumber(20) == 4)
	assert(t, tp.LineNumber(28) == 4)
	assert(t, tp.LineNumber(29) == 0)
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
	ok(t, v, err, MakeInt(0))

	err = ls.Append(MakeStr("a"))
	assert(t, err == nil)

	v, err = ls.Eq(NewList([]Value{}))
	ok(t, v, err, FALSE)

	v, err = ls.Eq(NewList([]Value{MakeStr("a")}))
	ok(t, v, err, TRUE)

	v, err = ls.Len()
	ok(t, v, err, ONE)

	v, err = ls.Get(MakeInt(0))
	ok(t, v, err, MakeStr("a"))

	err = ls.Set(MakeInt(0), MakeStr("b"))
	assert(t, err == nil)

	v, err = ls.Get(MakeInt(0))
	ok(t, v, err, MakeStr("b"))

	v, err = ls.Get(MakeInt(-1))
	fail(t, v, err, "IndexOutOfBounds")

	v, err = ls.Get(ONE)
	fail(t, v, err, "IndexOutOfBounds")

	err = ls.Set(MakeInt(-1), TRUE)
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
	h, err := NewFunc(&Template{}).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = NewList([]Value{}).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = newObj(map[string]Value{}).HashCode()
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
	ok(t, v, err, MakeInt(0))

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
}
