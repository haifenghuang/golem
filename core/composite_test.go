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

func newObj(fields map[string]Value) *Obj {
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
	ok(t, s, err, MakeStr("obj { }"))

	z, err := o.Eq(newObj(map[string]Value{}))
	ok(t, z, err, TRUE)
	z, err = o.Eq(newObj(map[string]Value{"a": Int(1)}))
	ok(t, z, err, FALSE)

	val, err := o.Add(MakeStr("a"))
	ok(t, val, err, MakeStr("obj { }a"))

	val, err = o.GetField("a")
	fail(t, val, err, "NoSuchField: Field 'a' not found.")

	//////////////////

	o = newObj(map[string]Value{"a": Int(1)})
	okType(t, o, TOBJ)

	s, err = o.String()
	ok(t, s, err, MakeStr("obj { a: 1 }"))

	z, err = o.Eq(newObj(map[string]Value{}))
	ok(t, z, err, FALSE)
	z, err = o.Eq(newObj(map[string]Value{"a": Int(1)}))
	ok(t, z, err, TRUE)

	val, err = o.Add(MakeStr("a"))
	ok(t, val, err, MakeStr("obj { a: 1 }a"))

	val, err = o.GetField("a")
	ok(t, val, err, Int(1))

	val, err = o.GetField("b")
	fail(t, val, err, "NoSuchField: Field 'b' not found.")

	err = o.PutField("a", Int(123))
	if err != nil {
		panic("unexpected error")
	}

	val, err = o.GetField("a")
	ok(t, val, err, Int(123))
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
	_, e5 := o.Sub(NULL)
	_, e6 := o.Mul(NULL)
	_, e7 := o.Div(NULL)

	_, e8 := o.Negate()

	_, e10 := o.GetField("")
	e11 := o.PutField("", NULL)

	uninitErr(t, e0)
	uninitErr(t, e1)
	uninitErr(t, e2)
	uninitErr(t, e3)
	uninitErr(t, e4)
	uninitErr(t, e5)
	uninitErr(t, e6)
	uninitErr(t, e7)
	uninitErr(t, e8)
	uninitErr(t, e10)
	uninitErr(t, e11)
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
