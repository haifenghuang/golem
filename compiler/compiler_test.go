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

func ok(t *testing.T, mod *g.Module, expect *g.Module) {

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

		if !reflect.DeepEqual(mt.OpcLines, et.OpcLines) {
			t.Error("OpcLines: ", mod, " != ", expect)
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
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(int64(-2)),
			g.Int(int64(2))},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.LOAD_NEG_ONE,
					g.ADD,
					g.LOAD_ZERO,
					g.ADD,
					g.LOAD_ZERO,
					g.ADD,
					g.LOAD_ONE,
					g.ADD,
					g.LOAD_CONST, 0, 1,
					g.ADD,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{16, 0}}}}})

	mod = NewCompiler(newAnalyzer("(2 + 3) * -4 / 10;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(int64(2)),
			g.Int(int64(3)),
			g.Int(int64(-4)),
			g.Int(int64(10))},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.ADD,
					g.LOAD_CONST, 0, 2,
					g.MUL,
					g.LOAD_CONST, 0, 3,
					g.DIV,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{16, 0}}}}})

	mod = NewCompiler(newAnalyzer("null / true + \nfalse;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_NULL,
					g.LOAD_TRUE,
					g.DIV,
					g.LOAD_FALSE,
					g.ADD,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{4, 2},
					g.OpcLine{5, 1},
					g.OpcLine{6, 0}}}}})

	mod = NewCompiler(newAnalyzer("'a' * 1.23e4;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.MakeStr("a"),
			g.MakeFloat(float64(12300))},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.MUL,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{8, 0}}}}})

	mod = NewCompiler(newAnalyzer("'a' == true;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.MakeStr("a")},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.LOAD_TRUE,
					g.EQ,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{6, 0}}}}})

	mod = NewCompiler(newAnalyzer("true != false;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_TRUE, g.LOAD_FALSE, g.NE,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{4, 0}}}}})

	mod = NewCompiler(newAnalyzer("true > false; true >= false;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_TRUE, g.LOAD_FALSE, g.GT,
					g.LOAD_TRUE, g.LOAD_FALSE, g.GTE,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{7, 0}}}}})

	mod = NewCompiler(newAnalyzer("true < false; true <= false; true <=> false;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_TRUE, g.LOAD_FALSE, g.LT,
					g.LOAD_TRUE, g.LOAD_FALSE, g.LTE,
					g.LOAD_TRUE, g.LOAD_FALSE, g.CMP,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{10, 0}}}}})

	mod = NewCompiler(newAnalyzer("let a = 2 && 3;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(int64(2)),
			g.Int(int64(3))},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{21, 0}}}}})

	mod = NewCompiler(newAnalyzer("let a = 2 || 3;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(int64(2)),
			g.Int(int64(3))},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{21, 0}}}}})
}

func TestAssignment(t *testing.T) {

	mod := NewCompiler(newAnalyzer("let a = 1;\nconst b = \n2;a = 3;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(2),
			g.Int(3)},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{5, 3},
					g.OpcLine{8, 2},
					g.OpcLine{11, 3},
					g.OpcLine{18, 0}}}}})
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
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(3),
			g.Int(2),
			g.Int(42)},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{17, 0}}}}})

	source = `let a = 1;
		if (false) {
		    let b = 2;
		} else {
		    let c = 3;
		}
		let d = 4;`

	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(2),
			g.Int(3),
			g.Int(4)},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{5, 2},
					g.OpcLine{9, 3},
					g.OpcLine{15, 4},
					g.OpcLine{18, 5},
					g.OpcLine{24, 7},
					g.OpcLine{30, 0}}}}})
}

func TestWhile(t *testing.T) {

	source := "let a = 1; while (0 < 1) { let b = 2; }"
	mod := NewCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(2)},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{20, 0}}}}})

	source = "let a = 'z'; while (0 < 1) \n{ break; continue; let b = 2; } let c = 3;"
	mod = NewCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.MakeStr("z"),
			g.Int(2),
			g.Int(3)},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{13, 2},
					g.OpcLine{34, 0}}}}})
}

func TestReturn(t *testing.T) {

	source := "return;"
	anl := newAnalyzer(source)
	mod := NewCompiler(anl).Compile()

	ok(t, mod, &g.Module{
		[]g.Value{},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.RETURN,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{2, 0}}}}})

	source = "let a = 1; return a \n- 2; a = 3;"
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(2),
			g.Int(3)},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 1},
					g.OpcLine{8, 2},
					g.OpcLine{12, 1},
					g.OpcLine{13, 2},
					g.OpcLine{20, 0}}}}})
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

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(42),
			g.Int(7)},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{0, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 1,
					g.STORE_LOCAL, 0, 0,
					g.NEW_FUNC, 0, 2,
					g.STORE_LOCAL, 0, 1,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{7, 3},
					g.OpcLine{13, 0}}},
			&g.Template{0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{4, 0}}},
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
					g.ADD,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 4},
					g.OpcLine{7, 7},
					g.OpcLine{24, 0}}},
			&g.Template{1, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.MUL,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 5},
					g.OpcLine{8, 0}}}}})

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

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(2),
			g.Int(3),
			g.Int(4)},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{7, 3},
					g.OpcLine{13, 4},
					g.OpcLine{19, 5},
					g.OpcLine{25, 6},
					g.OpcLine{32, 7},
					g.OpcLine{44, 0}}},

			&g.Template{0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0}}},

			&g.Template{1, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_LOCAL, 0, 0,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 3},
					g.OpcLine{4, 0}}},

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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 4},
					g.OpcLine{18, 0}}}}})
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

	ok(t, mod, &g.Module{
		[]g.Value{},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{0, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 1,
					g.STORE_LOCAL, 0, 0,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{7, 0}}},
			&g.Template{1, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 2,
					g.FUNC_LOCAL, 0, 0,
					g.RETURN,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 3},
					g.OpcLine{8, 0}}},
			&g.Template{1, 1, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CAPTURE, 0, 0,
					g.LOAD_LOCAL, 0, 0,
					g.ADD,
					g.DUP,
					g.STORE_CAPTURE, 0, 0,
					g.LOAD_CAPTURE, 0, 0,
					g.RETURN,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 4},
					g.OpcLine{12, 5},
					g.OpcLine{16, 0}}}}})

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

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(2)},
		nil,
		[]*g.ObjDef{},
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
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{7, 3},
					g.OpcLine{16, 0}}},
			&g.Template{1, 1, 1,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 2,
					g.FUNC_LOCAL, 0, 0,
					g.FUNC_CAPTURE, 0, 0,
					g.RETURN,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 4},
					g.OpcLine{11, 0}}},
			&g.Template{1, 2, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CAPTURE, 0, 0,
					g.LOAD_LOCAL, 0, 0,
					g.ADD,
					g.LOAD_CAPTURE, 0, 1,
					g.ADD,
					g.DUP,
					g.STORE_CAPTURE, 0, 0,
					g.LOAD_CAPTURE, 0, 0,
					g.RETURN,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 5},
					g.OpcLine{16, 6},
					g.OpcLine{20, 0}}}}})
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

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(2),
			g.Int(3),
			g.Int(4),
			g.Int(5)},
		nil,
		[]*g.ObjDef{
			&g.ObjDef{[]string{}},
			&g.ObjDef{[]string{"a"}},
			&g.ObjDef{[]string{"a", "b"}},
			&g.ObjDef{[]string{"a", "b", "c"}},
			&g.ObjDef{[]string{"d"}}},
		[]*g.Template{
			&g.Template{0, 0, 4,
				[]byte{
					g.LOAD_NULL,
					g.NEW_OBJ,
					g.INIT_OBJ, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.NEW_OBJ,
					g.LOAD_ZERO,
					g.INIT_OBJ, 0, 1,
					g.STORE_LOCAL, 0, 1,
					g.NEW_OBJ,
					g.LOAD_ONE,
					g.LOAD_CONST, 0, 0,
					g.INIT_OBJ, 0, 2,
					g.STORE_LOCAL, 0, 2,
					g.NEW_OBJ,
					g.LOAD_CONST, 0, 1,
					g.LOAD_CONST, 0, 2,
					g.NEW_OBJ,
					g.LOAD_CONST, 0, 3,
					g.INIT_OBJ, 0, 4,
					g.INIT_OBJ, 0, 3,
					g.STORE_LOCAL, 0, 3,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{8, 3},
					g.OpcLine{16, 4},
					g.OpcLine{27, 5},
					g.OpcLine{47, 0}}}}})

	source = `
let x = obj { a: 0 };
let y = x.a;
x.a = 3;
`
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()

	ok(t, mod, &g.Module{
		[]g.Value{
			g.MakeStr("a"),
			g.Int(3),
			g.MakeStr("a")},
		nil,
		[]*g.ObjDef{
			&g.ObjDef{[]string{"a"}}},
		[]*g.Template{
			&g.Template{0, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.NEW_OBJ,
					g.LOAD_ZERO,
					g.INIT_OBJ, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_LOCAL, 0, 0,
					g.GET_FIELD, 0, 0,
					g.STORE_LOCAL, 0, 1,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.PUT_FIELD, 0, 2,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{9, 3},
					g.OpcLine{18, 4},
					g.OpcLine{27, 0}}}}})

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

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(8),
			g.Int(5),
			g.MakeStr("plus"),
			g.MakeStr("minus"),
			g.MakeStr("x"),
			g.MakeStr("y"),
			g.MakeStr("x"),
			g.MakeStr("y")},
		nil,
		[]*g.ObjDef{
			&g.ObjDef{[]string{"x", "y", "plus", "minus"}}},
		[]*g.Template{
			&g.Template{0, 0, 4,
				[]byte{
					g.LOAD_NULL,
					g.NEW_OBJ,
					g.DUP,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.NEW_FUNC, 0, 1,
					g.FUNC_LOCAL, 0, 0,
					g.NEW_FUNC, 0, 2,
					g.FUNC_LOCAL, 0, 0,
					g.INIT_OBJ, 0, 0,
					g.STORE_LOCAL, 0, 1,
					g.LOAD_LOCAL, 0, 1,
					g.GET_FIELD, 0, 2,
					g.INVOKE, 0, 0,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_LOCAL, 0, 1,
					g.GET_FIELD, 0, 3,
					g.INVOKE, 0, 0,
					g.STORE_LOCAL, 0, 3,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{6, 3},
					g.OpcLine{9, 4},
					g.OpcLine{12, 5},
					g.OpcLine{18, 6},
					g.OpcLine{24, 7},
					g.OpcLine{27, 2},
					g.OpcLine{30, 8},
					g.OpcLine{42, 9},
					g.OpcLine{54, 0}}},
			&g.Template{0, 1, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CAPTURE, 0, 0,
					g.GET_FIELD, 0, 4,
					g.LOAD_CAPTURE, 0, 0,
					g.GET_FIELD, 0, 5,
					g.ADD,
					g.RETURN,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 5},
					g.OpcLine{15, 0}}},
			&g.Template{0, 1, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CAPTURE, 0, 0,
					g.GET_FIELD, 0, 6,
					g.LOAD_CAPTURE, 0, 0,
					g.GET_FIELD, 0, 7,
					g.SUB,
					g.RETURN,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 6},
					g.OpcLine{15, 0}}}}})
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

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(int64(10)),
			g.Int(int64(20))},
		nil,
		[]*g.ObjDef{},
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
					g.ADD,
					g.STORE_LOCAL, 0, 0,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_LOCAL, 0, 1,
					g.DUP,
					g.LOAD_NEG_ONE,
					g.ADD,
					g.STORE_LOCAL, 0, 1,
					g.STORE_LOCAL, 0, 3,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{7, 3},
					g.OpcLine{13, 4},
					g.OpcLine{25, 5},
					g.OpcLine{37, 0}}}}})

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

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(10),
			g.Int(20),
			g.MakeStr("x"),
			g.MakeStr("y")},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 4,
				[]byte{
					g.LOAD_NULL,
					g.NEW_OBJ,
					g.LOAD_CONST, 0, 0,
					g.INIT_OBJ, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.NEW_OBJ,
					g.LOAD_CONST, 0, 1,
					g.INIT_OBJ, 0, 1,
					g.STORE_LOCAL, 0, 1,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_ONE,
					g.INC_FIELD, 0, 2,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_LOCAL, 0, 1,
					g.LOAD_NEG_ONE,
					g.INC_FIELD, 0, 3,
					g.STORE_LOCAL, 0, 3,
					g.RETURN},
				[]g.OpcLine{
					g.OpcLine{0, 0},
					g.OpcLine{1, 2},
					g.OpcLine{11, 3},
					g.OpcLine{21, 4},
					g.OpcLine{31, 5},
					g.OpcLine{41, 0}}}}})
}
