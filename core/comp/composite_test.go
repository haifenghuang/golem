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
	//"fmt"
	g "golem/core"
	"reflect"
	"testing"
)

func assert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}

func ok(t *testing.T, val g.Value, err g.Error, expect g.Value) {

	if err != nil {
		t.Error(err, " != ", nil)
	}

	if !reflect.DeepEqual(val, expect) {
		panic("asdfasfad")
		t.Error(val, " != ", expect)
	}
}

func fail(t *testing.T, val g.Value, err g.Error, expect string) {

	if val != nil {
		t.Error(val, " != ", nil)
	}

	if err.Error() != expect {
		t.Error(err.Error(), " != ", expect)
	}
}

func okType(t *testing.T, val g.Value, expected g.Type) {
	tp, err := val.TypeOf()
	assert(t, tp == expected)
	assert(t, err == nil)
}

func newObj(fields map[string]g.Value) Obj {
	o := NewObj()
	def := &ObjDef{[]string{}}
	values := []g.Value{}
	for k, v := range fields {
		def.Keys = append(def.Keys, k)
		values = append(values, v)
	}
	o.Init(def, values)
	return o
}

func TestObj(t *testing.T) {
	o := newObj(map[string]g.Value{})
	okType(t, o, g.TOBJ)

	s, err := o.ToStr()
	ok(t, s, err, g.MakeStr("obj {}"))

	z, err := o.Eq(newObj(map[string]g.Value{}))
	ok(t, z, err, g.TRUE)
	z, err = o.Eq(newObj(map[string]g.Value{"a": g.ONE}))
	ok(t, z, err, g.FALSE)

	val, err := o.Add(g.MakeStr("a"))
	ok(t, val, err, g.MakeStr("obj {}a"))

	val, err = o.GetField(g.MakeStr("a"))
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	val, err = o.Get(g.MakeStr("a"))
	fail(t, val, err, "NoSuchField: Field 'a' not found")

	val, err = o.Get(g.ZERO)
	fail(t, val, err, "TypeMismatch: Expected 'Str'")

	//////////////////

	o = newObj(map[string]g.Value{"a": g.ONE})
	okType(t, o, g.TOBJ)

	s, err = o.ToStr()
	ok(t, s, err, g.MakeStr("obj { a: 1 }"))

	z, err = o.Eq(newObj(map[string]g.Value{}))
	ok(t, z, err, g.FALSE)
	z, err = o.Eq(newObj(map[string]g.Value{"a": g.ONE}))
	ok(t, z, err, g.TRUE)

	val, err = o.Add(g.MakeStr("a"))
	ok(t, val, err, g.MakeStr("obj { a: 1 }a"))

	val, err = o.GetField(g.MakeStr("a"))
	ok(t, val, err, g.ONE)

	val, err = o.GetField(g.MakeStr("b"))
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	val, err = o.Get(g.MakeStr("a"))
	ok(t, val, err, g.ONE)

	val, err = o.Get(g.MakeStr("b"))
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	err = o.PutField(g.MakeStr("a"), g.MakeInt(123))
	if err != nil {
		panic("unexpected error")
	}

	val, err = o.GetField(g.MakeStr("a"))
	ok(t, val, err, g.MakeInt(123))

	val, err = o.Get(g.MakeStr("a"))
	ok(t, val, err, g.MakeInt(123))

	err = o.Set(g.MakeStr("a"), g.MakeInt(456))
	if err != nil {
		panic("unexpected error")
	}

	val, err = o.GetField(g.MakeStr("a"))
	ok(t, val, err, g.MakeInt(456))

	val, err = o.Get(g.MakeStr("a"))
	ok(t, val, err, g.MakeInt(456))

	val, err = o.Has(g.MakeStr("a"))
	ok(t, val, err, g.TRUE)

	val, err = o.Has(g.MakeStr("abc"))
	ok(t, val, err, g.FALSE)

	val, err = o.Has(g.ZERO)
	fail(t, val, err, "TypeMismatch: Expected 'Str'")
}

func uninitErr(t *testing.T, err g.Error) {
	if err.Error() != "UninitializedObj: Obj is not yet initialized" {
		t.Error("bad uninitialized error")
	}
}

func TestUninitialized(t *testing.T) {
	o := NewObj()
	_, e0 := o.TypeOf()
	_, e1 := o.Eq(g.NULL)
	_, e2 := o.ToStr()
	_, e3 := o.Cmp(g.NULL)
	_, e4 := o.Add(g.NULL)

	_, e5 := o.GetField(g.MakeStr(""))
	e6 := o.PutField(g.MakeStr(""), g.NULL)

	_, e7 := o.Get(g.MakeStr(""))
	e8 := o.Set(g.MakeStr(""), g.NULL)

	_, e9 := o.Has(g.NULL)
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

func TestList(t *testing.T) {
	ls := NewList([]g.Value{})
	okType(t, ls, g.TLIST)

	var v g.Value
	v, err := ls.ToStr()
	ok(t, v, err, g.MakeStr("[]"))

	v, err = ls.Eq(NewList([]g.Value{}))
	ok(t, v, err, g.TRUE)

	v, err = ls.Eq(NewList([]g.Value{g.MakeStr("a")}))
	ok(t, v, err, g.FALSE)

	v, err = ls.Eq(g.NULL)
	ok(t, v, err, g.FALSE)

	v, err = ls.Len()
	ok(t, v, err, g.ZERO)

	err = ls.Append(g.MakeStr("a"))
	assert(t, err == nil)

	v, err = ls.Eq(NewList([]g.Value{}))
	ok(t, v, err, g.FALSE)

	v, err = ls.Eq(NewList([]g.Value{g.MakeStr("a")}))
	ok(t, v, err, g.TRUE)

	v, err = ls.Len()
	ok(t, v, err, g.ONE)

	v, err = ls.Get(g.ZERO)
	ok(t, v, err, g.MakeStr("a"))

	err = ls.Set(g.ZERO, g.MakeStr("b"))
	assert(t, err == nil)

	v, err = ls.Get(g.ZERO)
	ok(t, v, err, g.MakeStr("b"))

	v, err = ls.Get(g.NEG_ONE)
	fail(t, v, err, "IndexOutOfBounds")

	v, err = ls.Get(g.ONE)
	fail(t, v, err, "IndexOutOfBounds")

	err = ls.Set(g.NEG_ONE, g.TRUE)
	fail(t, nil, err, "IndexOutOfBounds")

	err = ls.Set(g.ONE, g.TRUE)
	fail(t, nil, err, "IndexOutOfBounds")

	v, err = ls.ToStr()
	ok(t, v, err, g.MakeStr("[ b ]"))

	err = ls.Append(g.MakeStr("z"))
	assert(t, err == nil)

	v, err = ls.ToStr()
	ok(t, v, err, g.MakeStr("[ b, z ]"))

	//////////////////////////////
	// sliceable

	ls = NewList([]g.Value{})
	v, err = ls.SliceFrom(g.ZERO)
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.SliceTo(g.ZERO)
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.SliceTo(g.ONE)
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.Slice(g.ZERO, g.ONE)
	fail(t, nil, err, "IndexOutOfBounds")

	ls = NewList([]g.Value{g.TRUE, g.FALSE, g.NULL})
	v, err = ls.SliceFrom(g.ONE)
	ok(t, v, err, NewList([]g.Value{g.FALSE, g.NULL}))
	v, err = ls.SliceTo(g.ONE)
	ok(t, v, err, NewList([]g.Value{g.TRUE}))
	v, err = ls.Slice(g.ZERO, g.ONE)
	ok(t, v, err, NewList([]g.Value{g.TRUE}))
	v, err = ls.Slice(g.ZERO, g.MakeInt(3))
	ok(t, v, err, NewList([]g.Value{g.TRUE, g.FALSE, g.NULL}))

	v, err = ls.Slice(g.ZERO, g.ZERO)
	ok(t, v, err, NewList([]g.Value{}))

	v, err = ls.Slice(g.MakeInt(2), g.ZERO)
	fail(t, nil, err, "IndexOutOfBounds")

	v, err = ls.SliceFrom(g.MakeInt(7))
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.SliceTo(g.MakeInt(7))
	fail(t, nil, err, "IndexOutOfBounds")
	v, err = ls.Slice(g.MakeInt(7), g.MakeInt(7))
	fail(t, nil, err, "IndexOutOfBounds")
}

func TestCompositeHashCode(t *testing.T) {
	h, err := NewDict(g.NewHashMap([]*g.HEntry{})).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = NewList([]g.Value{}).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")

	h, err = newObj(map[string]g.Value{}).HashCode()
	fail(t, h, err, "TypeMismatch: Expected Hashable Type")
}

func TestDict(t *testing.T) {
	d := NewDict(g.NewHashMap([]*g.HEntry{}))
	okType(t, d, g.TDICT)

	var v g.Value
	v, err := d.ToStr()
	ok(t, v, err, g.MakeStr("dict {}"))

	v, err = d.Eq(NewDict(g.NewHashMap([]*g.HEntry{})))
	ok(t, v, err, g.TRUE)

	v, err = d.Eq(g.NULL)
	ok(t, v, err, g.FALSE)

	v, err = d.Len()
	ok(t, v, err, g.ZERO)

	v, err = d.Get(g.MakeStr("a"))
	ok(t, v, err, g.NULL)

	err = d.Set(g.MakeStr("a"), g.ONE)
	assert(t, err == nil)

	v, err = d.Get(g.MakeStr("a"))
	ok(t, v, err, g.ONE)

	v, err = d.Eq(NewDict(g.NewHashMap([]*g.HEntry{})))
	ok(t, v, err, g.FALSE)

	v, err = d.Eq(NewDict(g.NewHashMap([]*g.HEntry{
		&g.HEntry{g.MakeStr("a"), g.ONE}})))
	ok(t, v, err, g.TRUE)

	v, err = d.Len()
	ok(t, v, err, g.ONE)

	v, err = d.ToStr()
	ok(t, v, err, g.MakeStr("dict { a: 1 }"))

	err = d.Set(g.MakeStr("b"), g.MakeInt(2))
	assert(t, err == nil)

	v, err = d.Get(g.MakeStr("b"))
	ok(t, v, err, g.MakeInt(2))

	v, err = d.ToStr()
	ok(t, v, err, g.MakeStr("dict { b: 2, a: 1 }"))

	tp := NewTuple([]g.Value{g.ONE, g.ZERO})
	d = NewDict(g.NewHashMap([]*g.HEntry{
		&g.HEntry{tp, g.TRUE}}))

	v, err = d.ToStr()
	ok(t, v, err, g.MakeStr("dict { (1, 0): true }"))

	v, err = d.Get(tp)
	ok(t, v, err, g.TRUE)
}

func TestTuple(t *testing.T) {
	var v g.Value

	tp := NewTuple([]g.Value{g.ONE, g.ZERO})
	okType(t, tp, g.TTUPLE)

	v, err := tp.Eq(NewTuple([]g.Value{g.ZERO, g.ZERO}))
	ok(t, v, err, g.FALSE)

	v, err = tp.Eq(NewTuple([]g.Value{g.ONE, g.ZERO}))
	ok(t, v, err, g.TRUE)

	v, err = tp.Eq(g.NULL)
	ok(t, v, err, g.FALSE)

	v, err = tp.Get(g.ZERO)
	ok(t, v, err, g.ONE)

	v, err = tp.Get(g.ONE)
	ok(t, v, err, g.ZERO)

	v, err = tp.Get(g.NEG_ONE)
	fail(t, v, err, "IndexOutOfBounds")

	v, err = tp.Get(g.MakeInt(2))
	fail(t, v, err, "IndexOutOfBounds")

	v, err = tp.ToStr()
	ok(t, v, err, g.MakeStr("(1, 0)"))

	v, err = tp.Len()
	ok(t, v, err, g.MakeInt(2))
}
