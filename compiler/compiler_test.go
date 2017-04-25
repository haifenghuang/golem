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

package compiler

import (
	//"fmt"
	"golem/analyzer"
	//"golem/ast"
	g "golem/core"
	"golem/core/comp"
	"golem/core/fn"
	"golem/parser"
	"golem/scanner"
	"reflect"
	"testing"
)

func ok(t *testing.T, mod *fn.Module, expect *fn.Module) {

	if !reflect.DeepEqual(mod.Pool, expect.Pool) {
		t.Error(mod, " != ", expect)
	}

	if len(mod.Templates) != len(expect.Templates) {
		t.Error(mod.Templates, " != ", expect.Templates)
	}

	for i := 0; i < len(mod.Templates); i++ {

		mt := mod.Templates[i]
		et := expect.Templates[i]

		if (mt.Arity != et.Arity) || (mt.NumCaptures != et.NumCaptures) || (mt.NumLocals != et.NumLocals) {
			t.Error(mod, " != ", expect)
		}

		if !reflect.DeepEqual(mt.OpCodes, et.OpCodes) {
			t.Error("OpCodes: ", mod, " != ", expect)
		}

		// checking OpcLines is optional
		if et.OpcLines != nil {
			if !reflect.DeepEqual(mt.OpcLines, et.OpcLines) {
				t.Error("OpcLines: ", mod, " != ", expect)
			}
		}
	}
}

func newAnalyzer(source string) analyzer.Analyzer {

	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner)
	mod, err := parser.ParseModule()
	if err != nil {
		panic("oops")
	}

	anl := analyzer.NewAnalyzer(mod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		panic("oops")
	}
	return anl
}

func TestExpression(t *testing.T) {

	mod := NewCompiler(newAnalyzer("-2 + -1 + -0 + 0 + 1 + 2;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(int64(-2)),
			g.MakeInt(int64(2))},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.LOAD_NEG_ONE,
					fn.ADD,
					fn.LOAD_ZERO,
					fn.ADD,
					fn.LOAD_ZERO,
					fn.ADD,
					fn.LOAD_ONE,
					fn.ADD,
					fn.LOAD_CONST, 0, 1,
					fn.ADD,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{16, 0}}}}})

	mod = NewCompiler(newAnalyzer("(2 + 3) * -4 / 10;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(int64(2)),
			g.MakeInt(int64(3)),
			g.MakeInt(int64(-4)),
			g.MakeInt(int64(10))},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.LOAD_CONST, 0, 1,
					fn.ADD,
					fn.LOAD_CONST, 0, 2,
					fn.MUL,
					fn.LOAD_CONST, 0, 3,
					fn.DIV,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{16, 0}}}}})

	mod = NewCompiler(newAnalyzer("null / true + \nfalse;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_NULL,
					fn.LOAD_TRUE,
					fn.DIV,
					fn.LOAD_FALSE,
					fn.ADD,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{4, 2},
					fn.OpcLine{5, 1},
					fn.OpcLine{6, 0}}}}})

	mod = NewCompiler(newAnalyzer("'a' * 1.23e4;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeStr("a"),
			g.MakeFloat(float64(12300))},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.LOAD_CONST, 0, 1,
					fn.MUL,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{8, 0}}}}})

	mod = NewCompiler(newAnalyzer("'a' == true;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeStr("a")},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.LOAD_TRUE,
					fn.EQ,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{6, 0}}}}})

	mod = NewCompiler(newAnalyzer("true != false;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_TRUE, fn.LOAD_FALSE, fn.NE,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{4, 0}}}}})

	mod = NewCompiler(newAnalyzer("true > false; true >= false;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_TRUE, fn.LOAD_FALSE, fn.GT,
					fn.LOAD_TRUE, fn.LOAD_FALSE, fn.GTE,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{7, 0}}}}})

	mod = NewCompiler(newAnalyzer("true < false; true <= false; true <=> false;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_TRUE, fn.LOAD_FALSE, fn.LT,
					fn.LOAD_TRUE, fn.LOAD_FALSE, fn.LTE,
					fn.LOAD_TRUE, fn.LOAD_FALSE, fn.CMP,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{10, 0}}}}})

	mod = NewCompiler(newAnalyzer("let a = 2 && 3;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(int64(2)),
			g.MakeInt(int64(3))},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.JUMP_FALSE, 0, 17,
					fn.LOAD_CONST, 0, 1,
					fn.JUMP_FALSE, 0, 17,
					fn.LOAD_TRUE,
					fn.JUMP, 0, 18,
					fn.LOAD_FALSE,
					fn.STORE_LOCAL, 0, 0,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{21, 0}}}}})

	mod = NewCompiler(newAnalyzer("let a = 2 || 3;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(int64(2)),
			g.MakeInt(int64(3))},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.JUMP_TRUE, 0, 13,
					fn.LOAD_CONST, 0, 1,
					fn.JUMP_FALSE, 0, 17,
					fn.LOAD_TRUE,
					fn.JUMP, 0, 18,
					fn.LOAD_FALSE,
					fn.STORE_LOCAL, 0, 0,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{21, 0}}}}})
}

func TestAssignment(t *testing.T) {

	mod := NewCompiler(newAnalyzer("let a = 1;\nconst b = \n2;a = 3;")).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(2),
			g.MakeInt(3)},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 2,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_ONE,
					fn.STORE_LOCAL, 0, 0,
					fn.LOAD_CONST, 0, 0,
					fn.STORE_LOCAL, 0, 1,
					fn.LOAD_CONST, 0, 1,
					fn.DUP,
					fn.STORE_LOCAL, 0, 0,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{5, 3},
					fn.OpcLine{8, 2},
					fn.OpcLine{11, 3},
					fn.OpcLine{18, 0}}}}})
}

func TestShift(t *testing.T) {

	a := 0x1234
	high, low := byte((a>>8)&0xFF), byte(a&0xFF)

	if high != 0x12 || low != 0x34 {
		panic("shift")
	}

	var b int = int(high)<<8 + int(low)
	if b != a {
		panic("shift")
	}
}

func TestIf(t *testing.T) {

	source := "if (3 == 2) { let a = 42; }"
	anl := newAnalyzer(source)
	mod := NewCompiler(anl).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(3),
			g.MakeInt(2),
			g.MakeInt(42)},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.LOAD_CONST, 0, 1,
					fn.EQ,
					fn.JUMP_FALSE, 0, 17,
					fn.LOAD_CONST, 0, 2,
					fn.STORE_LOCAL, 0, 0,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{17, 0}}}}})

	source = `let a = 1;
		if (false) {
		    let b = 2;
		} else {
		    let c = 3;
		}
		let d = 4;`

	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(2),
			g.MakeInt(3),
			g.MakeInt(4)},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 4,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_ONE,
					fn.STORE_LOCAL, 0, 0,
					fn.LOAD_FALSE,
					fn.JUMP_FALSE, 0, 18,
					fn.LOAD_CONST, 0, 0,
					fn.STORE_LOCAL, 0, 1,
					fn.JUMP, 0, 24,
					fn.LOAD_CONST, 0, 1,
					fn.STORE_LOCAL, 0, 2,
					fn.LOAD_CONST, 0, 2,
					fn.STORE_LOCAL, 0, 3,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{5, 2},
					fn.OpcLine{9, 3},
					fn.OpcLine{15, 4},
					fn.OpcLine{18, 5},
					fn.OpcLine{24, 7},
					fn.OpcLine{30, 0}}}}})
}

func TestWhile(t *testing.T) {

	source := "let a = 1; while (0 < 1) { let b = 2; }"
	mod := NewCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(2)},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 2,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_ONE,
					fn.STORE_LOCAL, 0, 0,
					fn.LOAD_ZERO,
					fn.LOAD_ONE,
					fn.LT,
					fn.JUMP_FALSE, 0, 20,
					fn.LOAD_CONST, 0, 0,
					fn.STORE_LOCAL, 0, 1,
					fn.JUMP, 0, 5,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{20, 0}}}}})

	source = "let a = 'z'; while (0 < 1) \n{ break; continue; let b = 2; } let c = 3;"
	mod = NewCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeStr("z"),
			g.MakeInt(2),
			g.MakeInt(3)},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 3,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.STORE_LOCAL, 0, 0,
					fn.LOAD_ZERO,
					fn.LOAD_ONE,
					fn.LT,
					fn.JUMP_FALSE, 0, 28,
					fn.JUMP, 0, 28,
					fn.JUMP, 0, 7,
					fn.LOAD_CONST, 0, 1,
					fn.STORE_LOCAL, 0, 1,
					fn.JUMP, 0, 7,
					fn.LOAD_CONST, 0, 2,
					fn.STORE_LOCAL, 0, 2,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{13, 2},
					fn.OpcLine{34, 0}}}}})
}

func TestReturn(t *testing.T) {

	source := "return;"
	anl := newAnalyzer(source)
	mod := NewCompiler(anl).Compile()

	ok(t, mod, &fn.Module{
		[]g.Value{},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.RETURN,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{2, 0}}}}})

	source = "let a = 1; return a \n- 2; a = 3;"
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()

	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(2),
			g.MakeInt(3)},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_ONE,
					fn.STORE_LOCAL, 0, 0,
					fn.LOAD_LOCAL, 0, 0,
					fn.LOAD_CONST, 0, 0,
					fn.SUB,
					fn.RETURN,
					fn.LOAD_CONST, 0, 1,
					fn.DUP,
					fn.STORE_LOCAL, 0, 0,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 1},
					fn.OpcLine{8, 2},
					fn.OpcLine{12, 1},
					fn.OpcLine{13, 2},
					fn.OpcLine{20, 0}}}}})
}

func TestFunc(t *testing.T) {

	source := `
let a = fn() { 42; };
let b = fn(x) {
    let c = fn(y) {
        y * 7;
    };
    x * x + c(x);
};
`
	anl := newAnalyzer(source)
	mod := NewCompiler(anl).Compile()

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(42),
			g.MakeInt(7)},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{0, 0, 2,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_FUNC, 0, 1,
					fn.STORE_LOCAL, 0, 0,
					fn.NEW_FUNC, 0, 2,
					fn.STORE_LOCAL, 0, 1,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 2},
					fn.OpcLine{7, 3},
					fn.OpcLine{13, 0}}},
			&fn.Template{0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 2},
					fn.OpcLine{4, 0}}},
			&fn.Template{1, 0, 2,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_FUNC, 0, 3,
					fn.STORE_LOCAL, 0, 1,
					fn.LOAD_LOCAL, 0, 0,
					fn.LOAD_LOCAL, 0, 0,
					fn.MUL,
					fn.LOAD_LOCAL, 0, 1,
					fn.LOAD_LOCAL, 0, 0,
					fn.INVOKE, 0, 1,
					fn.ADD,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 4},
					fn.OpcLine{7, 7},
					fn.OpcLine{24, 0}}},
			&fn.Template{1, 0, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_LOCAL, 0, 0,
					fn.LOAD_CONST, 0, 1,
					fn.MUL,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 5},
					fn.OpcLine{8, 0}}}}})

	source = `
let a = fn() { };
let b = fn(x) { x; };
let c = fn(x, y) { let z = 4; x * y * z; };
a();
b(1);
c(2, 3);
`
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(2),
			g.MakeInt(3),
			g.MakeInt(4)},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{0, 0, 3,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_FUNC, 0, 1,
					fn.STORE_LOCAL, 0, 0,
					fn.NEW_FUNC, 0, 2,
					fn.STORE_LOCAL, 0, 1,
					fn.NEW_FUNC, 0, 3,
					fn.STORE_LOCAL, 0, 2,
					fn.LOAD_LOCAL, 0, 0,
					fn.INVOKE, 0, 0,
					fn.LOAD_LOCAL, 0, 1,
					fn.LOAD_ONE,
					fn.INVOKE, 0, 1,
					fn.LOAD_LOCAL, 0, 2,
					fn.LOAD_CONST, 0, 0,
					fn.LOAD_CONST, 0, 1,
					fn.INVOKE, 0, 2,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 2},
					fn.OpcLine{7, 3},
					fn.OpcLine{13, 4},
					fn.OpcLine{19, 5},
					fn.OpcLine{25, 6},
					fn.OpcLine{32, 7},
					fn.OpcLine{44, 0}}},

			&fn.Template{0, 0, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0}}},

			&fn.Template{1, 0, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_LOCAL, 0, 0,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 3},
					fn.OpcLine{4, 0}}},

			&fn.Template{2, 0, 3,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 2,
					fn.STORE_LOCAL, 0, 2,
					fn.LOAD_LOCAL, 0, 0,
					fn.LOAD_LOCAL, 0, 1,
					fn.MUL,
					fn.LOAD_LOCAL, 0, 2,
					fn.MUL,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 4},
					fn.OpcLine{18, 0}}}}})
}

func TestCapture(t *testing.T) {

	source := `
const accumGen = fn(n) {
    return fn(i) {
        n = n + i;
        return n;
    };
};`

	anl := newAnalyzer(source)
	mod := NewCompiler(anl).Compile()

	ok(t, mod, &fn.Module{
		[]g.Value{},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{0, 0, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_FUNC, 0, 1,
					fn.STORE_LOCAL, 0, 0,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 2},
					fn.OpcLine{7, 0}}},
			&fn.Template{1, 0, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_FUNC, 0, 2,
					fn.FUNC_LOCAL, 0, 0,
					fn.RETURN,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 3},
					fn.OpcLine{8, 0}}},
			&fn.Template{1, 1, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CAPTURE, 0, 0,
					fn.LOAD_LOCAL, 0, 0,
					fn.ADD,
					fn.DUP,
					fn.STORE_CAPTURE, 0, 0,
					fn.LOAD_CAPTURE, 0, 0,
					fn.RETURN,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 4},
					fn.OpcLine{12, 5},
					fn.OpcLine{16, 0}}}}})

	source = `
let z = 2;
const accumGen = fn(n) {
    return fn(i) {
        n = n + i + z;
        return n;
    };
};`

	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()
	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(2)},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{0, 0, 2,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.STORE_LOCAL, 0, 0,
					fn.NEW_FUNC, 0, 1,
					fn.FUNC_LOCAL, 0, 0,
					fn.STORE_LOCAL, 0, 1,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 2},
					fn.OpcLine{7, 3},
					fn.OpcLine{16, 0}}},
			&fn.Template{1, 1, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_FUNC, 0, 2,
					fn.FUNC_LOCAL, 0, 0,
					fn.FUNC_CAPTURE, 0, 0,
					fn.RETURN,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 4},
					fn.OpcLine{11, 0}}},
			&fn.Template{1, 2, 1,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CAPTURE, 0, 0,
					fn.LOAD_LOCAL, 0, 0,
					fn.ADD,
					fn.LOAD_CAPTURE, 0, 1,
					fn.ADD,
					fn.DUP,
					fn.STORE_CAPTURE, 0, 0,
					fn.LOAD_CAPTURE, 0, 0,
					fn.RETURN,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 5},
					fn.OpcLine{16, 6},
					fn.OpcLine{20, 0}}}}})
}

func TestObj(t *testing.T) {

	source := `
let w = obj {};
let x = obj { a: 0 };
let y = obj { a: 1, b: 2 };
let z = obj { a: 3, b: 4, c: obj { d: 5 } };
`

	anl := newAnalyzer(source)
	mod := NewCompiler(anl).Compile()
	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(2),
			g.MakeInt(3),
			g.MakeInt(4),
			g.MakeInt(5)},
		nil,
		[]*comp.ObjDef{
			&comp.ObjDef{[]string{}},
			&comp.ObjDef{[]string{"a"}},
			&comp.ObjDef{[]string{"a", "b"}},
			&comp.ObjDef{[]string{"a", "b", "c"}},
			&comp.ObjDef{[]string{"d"}}},
		[]*fn.Template{
			&fn.Template{0, 0, 4,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_OBJ,
					fn.INIT_OBJ, 0, 0,
					fn.STORE_LOCAL, 0, 0,
					fn.NEW_OBJ,
					fn.LOAD_ZERO,
					fn.INIT_OBJ, 0, 1,
					fn.STORE_LOCAL, 0, 1,
					fn.NEW_OBJ,
					fn.LOAD_ONE,
					fn.LOAD_CONST, 0, 0,
					fn.INIT_OBJ, 0, 2,
					fn.STORE_LOCAL, 0, 2,
					fn.NEW_OBJ,
					fn.LOAD_CONST, 0, 1,
					fn.LOAD_CONST, 0, 2,
					fn.NEW_OBJ,
					fn.LOAD_CONST, 0, 3,
					fn.INIT_OBJ, 0, 4,
					fn.INIT_OBJ, 0, 3,
					fn.STORE_LOCAL, 0, 3,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 2},
					fn.OpcLine{8, 3},
					fn.OpcLine{16, 4},
					fn.OpcLine{27, 5},
					fn.OpcLine{47, 0}}}}})

	source = `
let x = obj { a: 0 };
let y = x.a;
x.a = 3;
`
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()

	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeStr("a"),
			g.MakeInt(3),
			g.MakeStr("a")},
		nil,
		[]*comp.ObjDef{
			&comp.ObjDef{[]string{"a"}}},
		[]*fn.Template{
			&fn.Template{0, 0, 2,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_OBJ,
					fn.LOAD_ZERO,
					fn.INIT_OBJ, 0, 0,
					fn.STORE_LOCAL, 0, 0,
					fn.LOAD_LOCAL, 0, 0,
					fn.GET_FIELD, 0, 0,
					fn.STORE_LOCAL, 0, 1,
					fn.LOAD_LOCAL, 0, 0,
					fn.LOAD_CONST, 0, 1,
					fn.PUT_FIELD, 0, 2,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 2},
					fn.OpcLine{9, 3},
					fn.OpcLine{18, 4},
					fn.OpcLine{27, 0}}}}})

	source = `
let a = obj {
    x: 8,
    y: 5,
    plus:  fn() { return this.x + this.y; },
    minus: fn() { return this.x - this.y; }
};
let b = a.plus();
let c = a.minus();
`
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()
	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(8),
			g.MakeInt(5),
			g.MakeStr("plus"),
			g.MakeStr("minus"),
			g.MakeStr("x"),
			g.MakeStr("y"),
			g.MakeStr("x"),
			g.MakeStr("y")},
		nil,
		[]*comp.ObjDef{
			&comp.ObjDef{[]string{"x", "y", "plus", "minus"}}},
		[]*fn.Template{
			&fn.Template{0, 0, 4,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_OBJ,
					fn.DUP,
					fn.STORE_LOCAL, 0, 0,
					fn.LOAD_CONST, 0, 0,
					fn.LOAD_CONST, 0, 1,
					fn.NEW_FUNC, 0, 1,
					fn.FUNC_LOCAL, 0, 0,
					fn.NEW_FUNC, 0, 2,
					fn.FUNC_LOCAL, 0, 0,
					fn.INIT_OBJ, 0, 0,
					fn.STORE_LOCAL, 0, 1,
					fn.LOAD_LOCAL, 0, 1,
					fn.GET_FIELD, 0, 2,
					fn.INVOKE, 0, 0,
					fn.STORE_LOCAL, 0, 2,
					fn.LOAD_LOCAL, 0, 1,
					fn.GET_FIELD, 0, 3,
					fn.INVOKE, 0, 0,
					fn.STORE_LOCAL, 0, 3,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 2},
					fn.OpcLine{6, 3},
					fn.OpcLine{9, 4},
					fn.OpcLine{12, 5},
					fn.OpcLine{18, 6},
					fn.OpcLine{24, 7},
					fn.OpcLine{27, 2},
					fn.OpcLine{30, 8},
					fn.OpcLine{42, 9},
					fn.OpcLine{54, 0}}},
			&fn.Template{0, 1, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CAPTURE, 0, 0,
					fn.GET_FIELD, 0, 4,
					fn.LOAD_CAPTURE, 0, 0,
					fn.GET_FIELD, 0, 5,
					fn.ADD,
					fn.RETURN,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 5},
					fn.OpcLine{15, 0}}},
			&fn.Template{0, 1, 0,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CAPTURE, 0, 0,
					fn.GET_FIELD, 0, 6,
					fn.LOAD_CAPTURE, 0, 0,
					fn.GET_FIELD, 0, 7,
					fn.SUB,
					fn.RETURN,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 6},
					fn.OpcLine{15, 0}}}}})
}

func TestPostfix(t *testing.T) {

	source := `
let a = 10;
let b = 20;
let c = a++;
let d = b--;
`
	anl := newAnalyzer(source)
	mod := NewCompiler(anl).Compile()
	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println("----------------------------")
	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(int64(10)),
			g.MakeInt(int64(20))},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 4,
				[]byte{
					fn.LOAD_NULL,
					fn.LOAD_CONST, 0, 0,
					fn.STORE_LOCAL, 0, 0,
					fn.LOAD_CONST, 0, 1,
					fn.STORE_LOCAL, 0, 1,
					fn.LOAD_LOCAL, 0, 0,
					fn.DUP,
					fn.LOAD_ONE,
					fn.ADD,
					fn.STORE_LOCAL, 0, 0,
					fn.STORE_LOCAL, 0, 2,
					fn.LOAD_LOCAL, 0, 1,
					fn.DUP,
					fn.LOAD_NEG_ONE,
					fn.ADD,
					fn.STORE_LOCAL, 0, 1,
					fn.STORE_LOCAL, 0, 3,
					fn.RETURN},
				[]fn.OpcLine{
					fn.OpcLine{0, 0},
					fn.OpcLine{1, 2},
					fn.OpcLine{7, 3},
					fn.OpcLine{13, 4},
					fn.OpcLine{25, 5},
					fn.OpcLine{37, 0}}}}})

	source = `
let a = obj { x: 10 };
let b = obj { y: 20 };
let c = a.x++;
let d = b.y--;
`
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()
	//fmt.Println("----------------------------")
	//fmt.Println(source)
	////fmt.Println("----------------------------")
	////fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//fmt.Println(mod)

	ok(t, mod, &fn.Module{
		[]g.Value{
			g.MakeInt(10),
			g.MakeInt(20),
			g.MakeStr("x"),
			g.MakeStr("y")},
		nil,
		[]*comp.ObjDef{},
		[]*fn.Template{
			&fn.Template{
				0, 0, 4,
				[]byte{
					fn.LOAD_NULL,
					fn.NEW_OBJ,
					fn.LOAD_CONST, 0, 0,
					fn.INIT_OBJ, 0, 0,
					fn.STORE_LOCAL, 0, 0,
					fn.NEW_OBJ,
					fn.LOAD_CONST, 0, 1,
					fn.INIT_OBJ, 0, 1,
					fn.STORE_LOCAL, 0, 1,
					fn.LOAD_LOCAL, 0, 0,
					fn.LOAD_ONE,
					fn.INC_FIELD, 0, 2,
					fn.STORE_LOCAL, 0, 2,
					fn.LOAD_LOCAL, 0, 1,
					fn.LOAD_NEG_ONE,
					fn.INC_FIELD, 0, 3,
					fn.STORE_LOCAL, 0, 3,
					fn.RETURN},
				nil}}})
}

//func TestList(t *testing.T) {
//
//	source := `
//let a = [];
//let b = [1];
//let c = [1,2,b];
//`
//	anl := newAnalyzer(source)
//	mod := NewCompiler(anl).Compile()
//	fmt.Println("----------------------------")
//	fmt.Println(source)
//	//fmt.Println("----------------------------")
//	//fmt.Printf("%s\n", ast.Dump(anl.Module()))
//	fmt.Println(mod)
//}
