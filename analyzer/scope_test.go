// Copyright 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain s copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analyzer

import (
	"bytes"
	"fmt"
	"golem/ast"
	"reflect"
	"testing"
)

func testGetOk(test *testing.T, s *scope, symbol string, v *ast.Variable) {

	entry, ok := s.get(symbol)

	if !ok {
		test.Error("not ok")
	}

	if !reflect.DeepEqual(entry, v) {
		test.Error(entry, " != ", v)
	}
}

func testGetMissing(test *testing.T, s *scope, symbol string) {
	_, ok := s.get(symbol)
	if ok {
		test.Error("not missing")
	}
}

func testScopeOk(test *testing.T, s *scope, expect string) {
	ds := dumpScope(s)
	if ("\n" + ds) != expect {
		fmt.Println("--------------------------------------------------------------")
		fmt.Println(ds)
		fmt.Println("--------------------------------------------------------------")
		fmt.Println(expect)
		test.Error("scope not ok")
	}
}

func dumpScope(s *scope) string {
	var buf bytes.Buffer

	for s != nil {
		buf.WriteString(s.String())
		buf.WriteString("\n")
		s = s.parent
	}

	return buf.String()
}

func TestGetPut(test *testing.T) {

	s := newFuncScope(nil)

	testGetMissing(test, s, "a")
	s.put("a", true)
	testGetOk(test, s, "a", &ast.Variable{0, true, false})

	t := newBlockScope(s)
	testGetOk(test, t, "a", &ast.Variable{0, true, false})

	testGetMissing(test, t, "b")
	t.put("b", false)
	testGetOk(test, t, "b", &ast.Variable{1, false, false})

	testGetMissing(test, s, "b")
}

func TestCaptureScope(test *testing.T) {

	s0 := newFuncScope(nil)
	s1 := newBlockScope(s0)
	s2 := newFuncScope(s1)
	s3 := newBlockScope(s2)
	s4 := newFuncScope(s3)
	s5 := newBlockScope(s4)

	s0.put("a", false)
	s1.put("b", false)
	s2.put("c", false)
	s3.put("d", false)
	s4.put("e", false)
	s5.put("f", false)

	s5.get("a")
	s5.get("c")

	testScopeOk(test, s5, `
Block defs:{f: (1,false,false)}
Func  defs:{e: (0,false,false)} captures:{a: (0,false,true), c: (1,false,true)} parentCaptures:{a: (0,false,true), c: (0,false,false)} numLocals:2
Block defs:{d: (1,false,false)}
Func  defs:{c: (0,false,false)} captures:{a: (0,false,true)} parentCaptures:{a: (0,false,false)} numLocals:2
Block defs:{b: (1,false,false)}
Func  defs:{a: (0,false,false)} captures:{} parentCaptures:{} numLocals:2
`)
}

func TestPlainObjScope(test *testing.T) {

	obj := &ast.ObjExpr{nil, nil, nil, -1, nil}

	s0 := newFuncScope(nil)
	s1 := newBlockScope(s0)
	s2 := newObjScope(s1, obj)

	testScopeOk(test, s2, `
Obj defs:{}
Block defs:{}
Func  defs:{} captures:{} parentCaptures:{} numLocals:0
`)

	if obj.LocalThisIndex != -1 {
		test.Error("LocalThisIndex is wrong", obj.LocalThisIndex, -1)
	}
}

func TestThisObjScope(test *testing.T) {

	obj2 := &ast.ObjExpr{nil, nil, nil, -1, nil}
	obj3 := &ast.ObjExpr{nil, nil, nil, -1, nil}

	s0 := newFuncScope(nil)
	s1 := newBlockScope(s0)
	s2 := newObjScope(s1, obj2)
	s3 := newObjScope(s2, obj3)

	s0.put("a", false)
	s1.put("b", false)
	s3.this()

	testScopeOk(test, s3, `
Obj defs:{this: (2,true,false)}
Obj defs:{}
Block defs:{b: (1,false,false)}
Func  defs:{a: (0,false,false)} captures:{} parentCaptures:{} numLocals:3
`)

	if obj2.LocalThisIndex != -1 {
		test.Error("LocalThisIndex is wrong", obj2.LocalThisIndex, -1)
	}
	if obj3.LocalThisIndex != 2 {
		test.Error("LocalThisIndex is wrong", obj3.LocalThisIndex, 2)
	}
}

func TestMethodScope(test *testing.T) {

	obj2 := &ast.ObjExpr{nil, nil, nil, -1, nil}

	s0 := newFuncScope(nil)
	s1 := newBlockScope(s0)
	s2 := newObjScope(s1, obj2)
	s3 := newFuncScope(s2)
	s4 := newBlockScope(s3)

	s4.this()
	// simulate encountering 'this' again within the s4 block
	s4.this()

	testScopeOk(test, s4, `
Block defs:{}
Func  defs:{} captures:{this: (0,true,true)} parentCaptures:{this: (0,true,false)} numLocals:0
Obj defs:{this: (0,true,false)}
Block defs:{}
Func  defs:{} captures:{} parentCaptures:{} numLocals:1
`)

	if obj2.LocalThisIndex != 0 {
		test.Error("LocalThisIndex is wrong", obj2.LocalThisIndex, 0)
	}
}
