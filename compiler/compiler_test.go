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
			t.Error(mod, " != ", expect)
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
			g.Int(int64(-1)),
			g.Int(int64(0)),
			g.Int(int64(0)),
			g.Int(int64(1)),
			g.Int(int64(2))},
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
					g.ADD,
					g.LOAD_CONST, 0, 3,
					g.ADD,
					g.LOAD_CONST, 0, 4,
					g.ADD,
					g.LOAD_CONST, 0, 5,
					g.ADD,
					g.RETURN}}}})

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
					g.RETURN}}}})

	mod = NewCompiler(newAnalyzer("null / true + false;")).Compile()
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
					g.RETURN}}}})

	mod = NewCompiler(newAnalyzer("'a' * 1.23e4;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Str("a"),
			g.Float(float64(12300))},
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
					g.RETURN}}}})

	mod = NewCompiler(newAnalyzer("'a' == true;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Str("a")},
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
					g.RETURN}}}})

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
					g.RETURN}}}})

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
					g.RETURN}}}})

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
					g.RETURN}}}})

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
					g.RETURN}}}})

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
					g.RETURN}}}})
}

func TestAssignment(t *testing.T) {

	mod := NewCompiler(newAnalyzer("let a = 1;const b = 2;a = 3;")).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(1),
			g.Int(2),
			g.Int(3)},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.STORE_LOCAL, 0, 1,
					g.LOAD_CONST, 0, 2,
					g.DUP,
					g.STORE_LOCAL, 0, 0,
					g.RETURN}}}})
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
					g.RETURN}}}})

	source = "let a = 1; if (false) { let b = 2; } else { let c = 3; } let d = 4;"
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(1),
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
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_FALSE,
					g.JUMP_FALSE, 0, 20,
					g.LOAD_CONST, 0, 1,
					g.STORE_LOCAL, 0, 1,
					g.JUMP, 0, 26,
					g.LOAD_CONST, 0, 2,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_CONST, 0, 3,
					g.STORE_LOCAL, 0, 3,
					g.RETURN}}}})
}

func TestWhile(t *testing.T) {
	source := "let a = 1; while (0 < 1) { let b = 2; }"
	mod := NewCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(1),
			g.Int(0),
			g.Int(1),
			g.Int(2)},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.LOAD_CONST, 0, 2,
					g.LT,
					g.JUMP_FALSE, 0, 26,
					g.LOAD_CONST, 0, 3,
					g.STORE_LOCAL, 0, 1,
					g.JUMP, 0, 7,
					g.RETURN}}}})

	// source = "a = 'z'; while (0 < 1) { break; continue; b = 2; } c = 3;"
	// anl := newAnalyzer(source)
	// mod = NewCompiler(anl).Compile()
	// fmt.Println("----------------------------")
	// fmt.Println(source)
	// fmt.Println("----------------------------")
	// fmt.Println(ast.DumpModule(anl.Module()))
	// fmt.Println(mod)

	source = "let a = 'z'; while (0 < 1) { break; continue; let b = 2; } let c = 3;"
	mod = NewCompiler(newAnalyzer(source)).Compile()
	ok(t, mod, &g.Module{
		[]g.Value{
			g.Str("z"),
			g.Int(0),
			g.Int(1),
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
					g.LOAD_CONST, 0, 1,
					g.LOAD_CONST, 0, 2,
					g.LT,
					g.JUMP_FALSE, 0, 32,
					g.JUMP, 0, 32,
					g.JUMP, 0, 7,
					g.LOAD_CONST, 0, 3,
					g.STORE_LOCAL, 0, 1,
					g.JUMP, 0, 7,
					g.LOAD_CONST, 0, 4,
					g.STORE_LOCAL, 0, 2,
					g.RETURN}}}})
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
					g.RETURN}}}})

	source = "let a = 1; return a - 2; a = 3;"
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(1),
			g.Int(2),
			g.Int(3)},
		nil,
		[]*g.ObjDef{},
		[]*g.Template{
			&g.Template{
				0, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.SUB,
					g.RETURN,
					g.LOAD_CONST, 0, 2,
					g.DUP,
					g.STORE_LOCAL, 0, 0,
					g.RETURN}}}})
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
					g.RETURN}},
			&g.Template{0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 0,
					g.RETURN}},
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
					g.RETURN}},
			&g.Template{1, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 1,
					g.MUL,
					g.RETURN}}}})

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
			g.Int(1),
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
					g.LOAD_CONST, 0, 0,
					g.INVOKE, 0, 1,
					g.LOAD_LOCAL, 0, 2,
					g.LOAD_CONST, 0, 1,
					g.LOAD_CONST, 0, 2,
					g.INVOKE, 0, 2,
					g.RETURN}},
			&g.Template{0, 0, 0,
				[]byte{
					g.LOAD_NULL,
					g.RETURN}},
			&g.Template{1, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_LOCAL, 0, 0,
					g.RETURN}},
			&g.Template{2, 0, 3,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CONST, 0, 3,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_LOCAL, 0, 1,
					g.MUL,
					g.LOAD_LOCAL, 0, 2,
					g.MUL,
					g.RETURN}}}})
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
					g.RETURN}},
			&g.Template{1, 0, 1,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 2,
					g.FUNC_LOCAL, 0, 0,
					g.RETURN,
					g.RETURN}},
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
					g.RETURN}}}})

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
					g.RETURN}},
			&g.Template{1, 1, 1,
				[]byte{
					g.LOAD_NULL,
					g.NEW_FUNC, 0, 2,
					g.FUNC_LOCAL, 0, 0,
					g.FUNC_CAPTURE, 0, 0,
					g.RETURN,
					g.RETURN}},
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
					g.RETURN}}}})
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
			g.Int(0),
			g.Int(1),
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
					g.LOAD_CONST, 0, 0,
					g.INIT_OBJ, 0, 1,
					g.STORE_LOCAL, 0, 1,
					g.NEW_OBJ,
					g.LOAD_CONST, 0, 1,
					g.LOAD_CONST, 0, 2,
					g.INIT_OBJ, 0, 2,
					g.STORE_LOCAL, 0, 2,
					g.NEW_OBJ,
					g.LOAD_CONST, 0, 3,
					g.LOAD_CONST, 0, 4,
					g.NEW_OBJ,
					g.LOAD_CONST, 0, 5,
					g.INIT_OBJ, 0, 4,
					g.INIT_OBJ, 0, 3,
					g.STORE_LOCAL, 0, 3,
					g.RETURN}}}})

	source = `
let x = obj { a: 0 };
let y = x.a;
x.a = 3;
`
	anl = newAnalyzer(source)
	mod = NewCompiler(anl).Compile()

	ok(t, mod, &g.Module{
		[]g.Value{
			g.Int(0),
			g.Str("a"),
			g.Int(3),
			g.Str("a")},
		nil,
		[]*g.ObjDef{
			&g.ObjDef{[]string{"a"}}},
		[]*g.Template{
			&g.Template{0, 0, 2,
				[]byte{
					g.LOAD_NULL,
					g.NEW_OBJ,
					g.LOAD_CONST, 0, 0,
					g.INIT_OBJ, 0, 0,
					g.STORE_LOCAL, 0, 0,
					g.LOAD_LOCAL, 0, 0,
					g.SELECT, 0, 1,
					g.STORE_LOCAL, 0, 1,
					g.LOAD_LOCAL, 0, 0,
					g.LOAD_CONST, 0, 2,
					g.PUT, 0, 3,
					g.RETURN}}}})

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
			g.Str("plus"),
			g.Str("minus"),
			g.Str("x"),
			g.Str("y"),
			g.Str("x"),
			g.Str("y")},
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
					g.SELECT, 0, 2,
					g.INVOKE, 0, 0,
					g.STORE_LOCAL, 0, 2,
					g.LOAD_LOCAL, 0, 1,
					g.SELECT, 0, 3,
					g.INVOKE, 0, 0,
					g.STORE_LOCAL, 0, 3,
					g.RETURN}},
			&g.Template{0, 1, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CAPTURE, 0, 0,
					g.SELECT, 0, 4,
					g.LOAD_CAPTURE, 0, 0,
					g.SELECT, 0, 5,
					g.ADD,
					g.RETURN,
					g.RETURN}},
			&g.Template{0, 1, 0,
				[]byte{
					g.LOAD_NULL,
					g.LOAD_CAPTURE, 0, 0,
					g.SELECT, 0, 6,
					g.LOAD_CAPTURE, 0, 0,
					g.SELECT, 0, 7,
					g.SUB,
					g.RETURN,
					g.RETURN}}}})
}
