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
	"golem/parser"
	"golem/scanner"
	"reflect"
	"testing"
)

func assert(t *testing.T, flag bool) {
	if !flag {
		t.Error("assertion failure")
	}
}

func ok(t *testing.T, mod *g.BytecodeModule, expect *g.BytecodeModule) {

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

		// checking LineNumberTable is optional
		if et.LineNumberTable != nil {
			if !reflect.DeepEqual(mt.LineNumberTable, et.LineNumberTable) {
				t.Error("LineNumberTable: ", mod, " != ", expect)
			}
		}
	}
}

func newAnalyzer(source string) analyzer.Analyzer {

	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner)
	mod, err := parser.ParseModule()
	if err != nil {
		panic(err)
	}

	anl := analyzer.NewAnalyzer(mod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		panic(err)
	}
	return anl
}

func contents() g.Struct {
	stc, _ := g.NewStruct([]*g.StructEntry{})
	return stc
}

func TestExpression(t *testing.T) {

	mod := NewCompiler(newAnalyzer("-2 + -1 + -0 + 0 + 1 + 2;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(-2)),
			g.MakeInt(int64(2))},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.LOAD_NEG_ONE,
					g.PLUS,
					g.LOAD_ZERO,
					g.PLUS,
					g.LOAD_ZERO,
					g.PLUS,
					g.LOAD_ONE,
					g.PLUS,
					g.LOAD_CONST, 0, 1,
					g.PLUS,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{16, 0}},
				nil}}, contents()})

	mod = NewCompiler(newAnalyzer("(2 + 3) * -4 / 10;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(2)),
			g.MakeInt(int64(3)),
			g.MakeInt(int64(-4)),
			g.MakeInt(int64(10))},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.PLUS,
					g.LOAD_CONST, 0, 2,
					g.MUL,
					g.LOAD_CONST, 0, 3,
					g.DIV,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{16, 0}},
				nil}}, contents()})

	mod = NewCompiler(newAnalyzer("null / true + \nfalse;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_NULL,
					g.LOAD_TRUE,
					g.DIV,
					g.LOAD_FALSE,
					g.PLUS,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{4, 2},
					{5, 1},
					{6, 0}},
				nil}}, contents()})

	mod = NewCompiler(newAnalyzer("'a' * 1.23e4;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeStr("a"),
			g.MakeFloat(float64(12300))},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.MUL,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{8, 0}},
				nil}}, contents()})

	mod = NewCompiler(newAnalyzer("'a' == true;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeStr("a")},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.LOAD_TRUE,
					g.EQ,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{6, 0}},
				nil}}, contents()})

	mod = NewCompiler(newAnalyzer("true != false;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_TRUE, g.LOAD_FALSE, g.NE,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{4, 0}},
				nil}}, contents()})

	mod = NewCompiler(newAnalyzer("true > false; true >= false;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_TRUE, g.LOAD_FALSE, g.GT,
					g.LOAD_TRUE, g.LOAD_FALSE, g.GTE,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{7, 0}},
				nil}}, contents()})

	mod = NewCompiler(newAnalyzer("true < false; true <= false; true <=> false;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_TRUE, g.LOAD_FALSE, g.LT,
					g.LOAD_TRUE, g.LOAD_FALSE, g.LTE,
					g.LOAD_TRUE, g.LOAD_FALSE, g.CMP,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{10, 0}},
				nil}}, contents()})

	mod = NewCompiler(newAnalyzer("let a = 2 && 3;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(2)),
			g.MakeInt(int64(3))},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.JUMP_FALSE, 0, 17,
					g.LOAD_CONST, 0, 1,
					g.JUMP_FALSE, 0, 17,
					g.LOAD_TRUE,
					g.JUMP, 0, 18,
					g.LOAD_FALSE,
					g.STORE_LOCAL, 0, 0,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{21, 0}},
				nil}}, contents()})

	mod = NewCompiler(newAnalyzer("let a = 2 || 3;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(2)),
			g.MakeInt(int64(3))},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.JUMP_TRUE, 0, 13,
					g.LOAD_CONST, 0, 1,
					g.JUMP_FALSE, 0, 17,
					g.LOAD_TRUE,
					g.JUMP, 0, 18,
					g.LOAD_FALSE,
					g.STORE_LOCAL, 0, 0,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{21, 0}},
				nil}}, contents()})
}

func TestAssignment(t *testing.T) {

	mod := NewCompiler(newAnalyzer("let a = 1;\nconst b = \n2;a = 3;")).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2),
			g.MakeInt(3)},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_ONE,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 1,
					g.LOAD_CONST, 0, 1,
					g.DUP,
					g.STORE_LOCAL, 0, 0,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{5, 3},
					{8, 2},
					{11, 3},
					{18, 0}},
				nil}}, contents()})
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
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(3),
			g.MakeInt(2),
			g.MakeInt(42)},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.EQ,
					g.JUMP_FALSE, 0, 17,
					g.LOAD_CONST, 0, 2,
					g.STORE_LOCAL, 0, 0,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{17, 0}},
				nil}}, contents()})

	source = `let a = 1;
		if (false) {
		    let b = 2;
		} else {
		    let c = 3;
		}
		let d = 4;`

	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2),
			g.MakeInt(3),
			g.MakeInt(4)},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 4,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_ONE,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_FALSE,
					g.JUMP_FALSE, 0, 18,
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 1,
					g.JUMP, 0, 24,
					g.LOAD_CONST, 0, 1,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_CONST, 0, 2,
					g.STORE_LOCAL, 0, 3,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{5, 2},
					{9, 3},
					{15, 4},
					{18, 5},
					{24, 7},
					{30, 0}},
				nil}}, contents()})
}

func TestWhile(t *testing.T) {

	source := "let a = 1; while (0 < 1) { let b = 2; }"
	mod := NewCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2)},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_ONE,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_ZERO,
					g.LOAD_ONE,
					g.LT,
					g.JUMP_FALSE, 0, 20,
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 1,
					g.JUMP, 0, 5,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{20, 0}},
				nil}}, contents()})

	source = "let a = 'z'; while (0 < 1) \n{ break; continue; let b = 2; } let c = 3;"
	mod = NewCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeStr("z"),
			g.MakeInt(2),
			g.MakeInt(3)},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 3,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_ZERO,
					g.LOAD_ONE,
					g.LT,
					g.JUMP_FALSE, 0, 28,
					g.JUMP, 0, 28,
					g.JUMP, 0, 7,
					g.LOAD_CONST, 0, 1,
					g.STORE_LOCAL, 0, 1,
					g.JUMP, 0, 7,
					g.LOAD_CONST, 0, 2,
					g.STORE_LOCAL, 0, 2,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{13, 2},
					{34, 0}},
				nil}}, contents()})
}

func TestReturn(t *testing.T) {

	source := "return;"
	anl := newAnalyzer(source)
	mod := NewCompiler(anl).Compile()

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.RETURN,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{2, 0}},
				nil}}, contents()})

	source = "let a = 1; return a \n- 2; a = 3;"
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2),
			g.MakeInt(3)},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_ONE,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 0,
					g.SUB,
					g.RETURN,
					g.LOAD_CONST, 0, 1,
					g.DUP,
					g.STORE_LOCAL, 0, 0,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 1},
					{8, 2},
					{12, 1},
					{13, 2},
					{20, 0}},
				nil}}, contents()})
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
	//fmt.Printf("%s\n", ast.Dump(anl.BytecodeModule()))
	//fmt.Println(mod)

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(42),
			g.MakeInt(7)},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{0, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 1,
					g.STORE_LOCAL, 0, 0,
					g.NEW_FUNC, 0, 2,
					g.STORE_LOCAL, 0, 1,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 3},
					{13, 0}},
				nil},
			&g.Template{0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{4, 0}},
				nil},
			&g.Template{1, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 3,
					g.STORE_LOCAL, 0, 1,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_LOCAL, 0, 0,
					g.MUL,
					g.LOAD_LOCAL, 0, 1,
					g.LOAD_LOCAL, 0, 0,
					g.INVOKE, 0, 1,
					g.PLUS,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 4},
					{7, 7},
					{24, 0}},
				nil},
			&g.Template{1, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.MUL,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 5},
					{8, 0}},
				nil}}, contents()})

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
	//fmt.Printf("%s\n", ast.Dump(anl.BytecodeModule()))
	//fmt.Println(mod)

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2),
			g.MakeInt(3),
			g.MakeInt(4)},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{0, 0, 3,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 1,
					g.STORE_LOCAL, 0, 0,
					g.NEW_FUNC, 0, 2,
					g.STORE_LOCAL, 0, 1,
					g.NEW_FUNC, 0, 3,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_LOCAL, 0, 0,
					g.INVOKE, 0, 0,
					g.LOAD_LOCAL, 0, 1,
					g.LOAD_ONE,
					g.INVOKE, 0, 1,
					g.LOAD_LOCAL, 0, 2,
					g.LOAD_CONST, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.INVOKE, 0, 2,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 3},
					{13, 4},
					{19, 5},
					{25, 6},
					{32, 7},
					{44, 0}},
				nil},

			&g.Template{0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0}},
				nil},

			&g.Template{1, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_LOCAL, 0, 0,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 3},
					{4, 0}},
				nil},

			&g.Template{2, 0, 3,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 2,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_LOCAL, 0, 1,
					g.MUL,
					g.LOAD_LOCAL, 0, 2,
					g.MUL,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 4},
					{18, 0}},
				nil}}, contents()})
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

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{0, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 1,
					g.STORE_LOCAL, 0, 0,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 0}},
				nil},
			&g.Template{1, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 2,
					g.FUNC_LOCAL, 0, 0,
					g.RETURN,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 3},
					{8, 0}},
				nil},
			&g.Template{1, 1, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CAPTURE, 0, 0,
					g.LOAD_LOCAL, 0, 0,
					g.PLUS,
					g.DUP,
					g.STORE_CAPTURE, 0, 0,
					g.LOAD_CAPTURE, 0, 0,
					g.RETURN,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 4},
					{12, 5},
					{16, 0}},
				nil}}, contents()})

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
	//fmt.Printf("%s\n", ast.Dump(anl.BytecodeModule()))
	//fmt.Println(mod)

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(2)},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{0, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.NEW_FUNC, 0, 1,
					g.FUNC_LOCAL, 0, 0,
					g.STORE_LOCAL, 0, 1,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 3},
					{16, 0}},
				nil},
			&g.Template{1, 1, 1,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 2,
					g.FUNC_LOCAL, 0, 0,
					g.FUNC_CAPTURE, 0, 0,
					g.RETURN,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 4},
					{11, 0}},
				nil},
			&g.Template{1, 2, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CAPTURE, 0, 0,
					g.LOAD_LOCAL, 0, 0,
					g.PLUS,
					g.LOAD_CAPTURE, 0, 1,
					g.PLUS,
					g.DUP,
					g.STORE_CAPTURE, 0, 0,
					g.LOAD_CAPTURE, 0, 0,
					g.RETURN,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 5},
					{16, 6},
					{20, 0}},
				nil}}, contents()})
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
	//fmt.Printf("%s\n", ast.Dump(anl.BytecodeModule()))
	//fmt.Println(mod)

	ok(t, mod, &g.BytecodeModule{
		[]g.Basic{
			g.MakeInt(int64(10)),
			g.MakeInt(int64(20))},
		nil,
		[][]*g.StructEntryDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 4,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.STORE_LOCAL, 0, 1,
					g.LOAD_LOCAL, 0, 0,
					g.DUP,
					g.LOAD_ONE,
					g.PLUS,
					g.STORE_LOCAL, 0, 0,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_LOCAL, 0, 1,
					g.DUP,
					g.LOAD_NEG_ONE,
					g.PLUS,
					g.STORE_LOCAL, 0, 1,
					g.STORE_LOCAL, 0, 3,
					g.RETURN},
				[]g.LineNumberEntry{
					{0, 0},
					{1, 2},
					{7, 3},
					{13, 4},
					{25, 5},
					{37, 0}},
				nil}}, contents()})
}

func TestPool(t *testing.T) {
	pool := g.EmptyHashMap()

	assert(t, poolIndex(pool, g.MakeInt(4)) == 0)
	assert(t, poolIndex(pool, g.MakeStr("a")) == 1)
	assert(t, poolIndex(pool, g.MakeFloat(1.0)) == 2)
	assert(t, poolIndex(pool, g.MakeStr("a")) == 1)
	assert(t, poolIndex(pool, g.MakeInt(4)) == 0)

	slice := makePoolSlice(pool)
	assert(t, reflect.DeepEqual(
		slice,
		[]g.Basic{
			g.MakeInt(4),
			g.MakeStr("a"),
			g.MakeFloat(1.0)}))
}

func TestTry(t *testing.T) {

	source := `
let a = 1;
try {
    a++;
}
finally {
    a++;
}
assert(a == 2);
`
	anl := newAnalyzer(source)
	mod := NewCompiler(anl).Compile()
	assert(t, mod.Templates[0].ExceptionHandlers[0] ==
		g.ExceptionHandler{5, 14, -1, 14})

	source = `
try {
    try {
        3 / 0;
    } catch e2 {
        assert(1,2);
    }
} catch e {
    println(e);
}
`
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()
	//	fmt.Println("----------------------------")
	//	fmt.Println(source)
	//	fmt.Println("----------------------------")
	//	fmt.Printf("%s\n", ast.Dump(anl.Module()))
	//	fmt.Println(mod)
}
