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

	s, err := o.String()
	ok(t, s, err, MakeStr("obj {}"))

	z, err := o.Eq(newObj(map[string]Value{}))
	ok(t, z, err, TRUE)
	z, err = o.Eq(newObj(map[string]Value{"a": MakeInt(1)}))
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

	o = newObj(map[string]Value{"a": MakeInt(1)})
	okType(t, o, TOBJ)

	s, err = o.String()
	ok(t, s, err, MakeStr("obj { a: 1 }"))

	z, err = o.Eq(newObj(map[string]Value{}))
	ok(t, z, err, FALSE)
	z, err = o.Eq(newObj(map[string]Value{"a": MakeInt(1)}))
	ok(t, z, err, TRUE)

	val, err = o.Add(MakeStr("a"))
	ok(t, val, err, MakeStr("obj { a: 1 }a"))

	val, err = o.GetField(MakeStr("a"))
	ok(t, val, err, MakeInt(1))

	val, err = o.GetField(MakeStr("b"))
	fail(t, val, err, "NoSuchField: Field 'b' not found")

	val, err = o.Get(MakeStr("a"))
	ok(t, val, err, MakeInt(1))

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
	_, e2 := o.String()
	_, e3 := o.Cmp(NULL)
	_, e4 := o.Add(NULL)

	_, e5 := o.GetField(MakeStr(""))
	e6 := o.PutField(MakeStr(""), NULL)

	uninitErr(t, e0)
	uninitErr(t, e1)
	uninitErr(t, e2)
	uninitErr(t, e3)
	uninitErr(t, e4)
	uninitErr(t, e5)
	uninitErr(t, e6)
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
	v, err := ls.String()
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
	ok(t, v, err, MakeInt(1))

	v, err = ls.Get(MakeInt(0))
	ok(t, v, err, MakeStr("a"))

	err = ls.Set(MakeInt(0), MakeStr("b"))
	assert(t, err == nil)

	v, err = ls.Get(MakeInt(0))
	ok(t, v, err, MakeStr("b"))

	v, err = ls.Get(MakeInt(-1))
	fail(t, v, err, "IndexOutOfBounds")

	v, err = ls.Get(MakeInt(1))
	fail(t, v, err, "IndexOutOfBounds")

	err = ls.Set(MakeInt(-1), TRUE)
	fail(t, nil, err, "IndexOutOfBounds")

	err = ls.Set(MakeInt(1), TRUE)
	fail(t, nil, err, "IndexOutOfBounds")

	v, err = ls.String()
	ok(t, v, err, MakeStr("[ b ]"))

	err = ls.Append(MakeStr("z"))
	assert(t, err == nil)

	v, err = ls.String()
	ok(t, v, err, MakeStr("[ b, z ]"))
}
