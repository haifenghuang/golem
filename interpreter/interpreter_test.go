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

package interpreter

import (
	//"fmt"
	"golem/analyzer"
	"golem/compiler"
	g "golem/core"
	"golem/parser"
	"golem/scanner"
	"reflect"
	"testing"
)

func ok_expr(t *testing.T, source string, expect g.Value) {
	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod)

	result, errStack := intp.Init()
	if errStack != nil {
		panic(errStack)
	}

	b, err := result.Eq(expect)
	if err != nil {
		panic(err)
	}
	if !b.BoolVal() {
		t.Error(result, " != ", expect)
	}
}

func ok_ref(t *testing.T, ref *g.Ref, expect g.Value) {
	b, err := ref.Val.Eq(expect)
	if err != nil {
		panic(err)
	}
	if !b.BoolVal() {
		t.Error(ref.Val, " != ", expect)
	}
}

func ok_mod(t *testing.T, source string, expectResult g.Value, expectLocals []*g.Ref) {
	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod)

	result, errStack := intp.Init()
	if errStack != nil {
		panic(errStack)
	}

	b, err := result.Eq(expectResult)
	if err != nil {
		panic(err)
	}
	if !b.BoolVal() {
		t.Error(result, " != ", expectResult)
	}

	if !reflect.DeepEqual(mod.Locals, expectLocals) {
		t.Error(mod.Locals, " != ", expectLocals)
	}
}

func fail_expr(t *testing.T, source string, expect string) {

	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod)

	result, errStack := intp.Init()
	if result != nil {
		panic(result)
	}

	if errStack.Err.Error() != expect {
		t.Error(errStack.Err.Error(), " != ", expect)
	}
}

func fail(t *testing.T, source string, expect *ErrorStack) {

	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod)

	result, errStack := intp.Init()
	if result != nil {
		panic(result)
	}

	if !reflect.DeepEqual(errStack, expect) {
		t.Error(errStack, " != ", expect)
	}
}

func newCompiler(source string) compiler.Compiler {
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

	return compiler.NewCompiler(anl)
}

func interpret(mod *g.Module) {
	intp := NewInterpreter(mod)
	_, err := intp.Init()
	if err != nil {
		panic(err)
	}
}

func TestExpressions(t *testing.T) {

	ok_expr(t, "(2 + 3) * -4 / 10;", g.MakeInt(-2))

	ok_expr(t, "(2*2*2*2 + 2*3*(8 - 1) + 2) / (17 - 2*2*2 - -1);", g.MakeInt(6))

	ok_expr(t, "true + 'a';", g.MakeStr("truea"))
	ok_expr(t, "'a' + true;", g.MakeStr("atrue"))

	fail_expr(t, "true + null;", "TypeMismatch: Expected Number Type")
	fail_expr(t, "1 + null;", "TypeMismatch: Expected Number Type")
	fail_expr(t, "null + 1;", "NullValue")

	ok_expr(t, "true == 'a';", g.FALSE)
	ok_expr(t, "3 * 7 + 4 == 5 * 5;", g.TRUE)
	ok_expr(t, "1 != 1;", g.FALSE)
	ok_expr(t, "1 != 2;", g.TRUE)

	ok_expr(t, "!false;", g.TRUE)
	ok_expr(t, "!true;", g.FALSE)
	fail_expr(t, "!null;", "TypeMismatch: Expected 'Bool'")

	fail_expr(t, "!'a';", "TypeMismatch: Expected 'Bool'")
	fail_expr(t, "!1;", "TypeMismatch: Expected 'Bool'")
	fail_expr(t, "!1.0;", "TypeMismatch: Expected 'Bool'")

	ok_expr(t, "1 < 2;", g.TRUE)
	ok_expr(t, "1 <= 2;", g.TRUE)
	ok_expr(t, "1 > 2;", g.FALSE)
	ok_expr(t, "1 >= 2;", g.FALSE)

	ok_expr(t, "2 < 2;", g.FALSE)
	ok_expr(t, "2 <= 2;", g.TRUE)
	ok_expr(t, "2 > 2;", g.FALSE)
	ok_expr(t, "2 >= 2;", g.TRUE)

	ok_expr(t, "1 <=> 2;", g.MakeInt(-1))
	ok_expr(t, "2 <=> 2;", g.MakeInt(0))
	ok_expr(t, "2 <=> 1;", g.MakeInt(1))

	ok_expr(t, "true  && true;", g.TRUE)
	ok_expr(t, "true  && false;", g.FALSE)
	ok_expr(t, "false && true;", g.FALSE)
	ok_expr(t, "false && 12;", g.FALSE)
	fail_expr(t, "12  && false;", "TypeMismatch: Expected 'Bool'")

	ok_expr(t, "true  || true;", g.TRUE)
	ok_expr(t, "true  || false;", g.TRUE)
	ok_expr(t, "false || true;", g.TRUE)
	ok_expr(t, "false || false;", g.FALSE)
	ok_expr(t, "true  || 12;", g.TRUE)
	fail_expr(t, "12  || true;", "TypeMismatch: Expected 'Bool'")

	ok_expr(t, "~0;", g.MakeInt(-1))

	ok_expr(t, "8 % 2;", g.MakeInt(8%2))
	ok_expr(t, "8 & 2;", g.MakeInt(int64(8)&int64(2)))
	ok_expr(t, "8 | 2;", g.MakeInt(8|2))
	ok_expr(t, "8 ^ 2;", g.MakeInt(8^2))
	ok_expr(t, "8 << 2;", g.MakeInt(8<<2))
	ok_expr(t, "8 >> 2;", g.MakeInt(8>>2))

	ok_expr(t, "[true][0];", g.TRUE)
	ok_expr(t, "'abc'[1];", g.MakeStr("b"))
	fail_expr(t, "[true][2];", "IndexOutOfBounds")
	fail_expr(t, "'abc'[-1];", "IndexOutOfBounds")

	ok_expr(t, "'abc'[1:];", g.MakeStr("bc"))
	ok_expr(t, "'abc'[:1];", g.MakeStr("a"))
	ok_expr(t, "'abcd'[1:3];", g.MakeStr("bc"))
	ok_expr(t, "'abcd'[1:1];", g.MakeStr(""))

	ok_expr(t, "[6,7,8][1:];", g.NewList([]g.Value{g.MakeInt(7), g.MakeInt(8)}))
	ok_expr(t, "[6,7,8][:1];", g.NewList([]g.Value{g.MakeInt(6)}))
	ok_expr(t, "[6,7,8,9][1:3];", g.NewList([]g.Value{g.MakeInt(7), g.MakeInt(8)}))
	ok_expr(t, "[6,7,8,9][1:1];", g.NewList([]g.Value{}))

	ok_expr(t, "obj{a: 1}['a'];", g.ONE)
	ok_expr(t, "obj{a: 1} has 'a';", g.TRUE)
	ok_expr(t, "obj{a: 1} has 'b';", g.FALSE)

	fail_expr(t, "obj{a: 1}[0];", "TypeMismatch: Expected 'Str'")
	fail_expr(t, "obj{a: this has 1};", "UninitializedObj: Obj is not yet initialized")
}

func TestAssignment(t *testing.T) {
	ok_mod(t, `
let a = 1;
const B = 2;
a = a + B;
`,
		g.MakeInt(3),
		[]*g.Ref{
			&g.Ref{g.MakeInt(3)},
			&g.Ref{g.MakeInt(2)}})

	ok_mod(t, `
let a = 1;
a = a + 41;
const B = a / 6;
let c = B + 3;
c = (c + a)/13;
`,
		g.MakeInt(4),
		[]*g.Ref{
			&g.Ref{g.MakeInt(42)},
			&g.Ref{g.MakeInt(7)},
			&g.Ref{g.MakeInt(4)}})

	ok_mod(t, `
let a = 1;
let b = a += 3;
let c = ~0;
c -= -2;
c <<= 4;
b *= 2;
`,
		g.MakeInt(8),
		[]*g.Ref{
			&g.Ref{g.MakeInt(4)},
			&g.Ref{g.MakeInt(8)},
			&g.Ref{g.MakeInt(16)}})

	ok_mod(t, `
let a = 1;
let b = 2;
a = b = 11;
b = a %= 4;
`,
		g.MakeInt(3),
		[]*g.Ref{
			&g.Ref{g.MakeInt(3)},
			&g.Ref{g.MakeInt(3)}})
}

func TestIf(t *testing.T) {

	ok_mod(t, "let a = 1; if (true) { a = 2; }",
		g.MakeInt(2),
		[]*g.Ref{&g.Ref{g.MakeInt(2)}})

	ok_mod(t, "let a = 1; if (false) { a = 2; }",
		g.NULL,
		[]*g.Ref{&g.Ref{g.MakeInt(1)}})

	ok_mod(t, "let a = 1; if (1 == 1) { a = 2; } else { a = 3; } let b = 4;",
		g.MakeInt(2),
		[]*g.Ref{
			&g.Ref{g.MakeInt(2)},
			&g.Ref{g.MakeInt(4)}})

	ok_mod(t, "let a = 1; if (1 == 2) { a = 2; } else { a = 3; } const b = 4;",
		g.MakeInt(3),
		[]*g.Ref{
			&g.Ref{g.MakeInt(3)},
			&g.Ref{g.MakeInt(4)}})
}

func TestWhile(t *testing.T) {

	//	source := `
	//a = 1;
	//while (a < 11) {
	//    if (a == 4) { a = a + 2; break; }
	//    a = a + 1;
	//}`
	//	mod := newCompiler(source).Compile()
	//	fmt.Println("----------------------------")
	//	fmt.Println(source)
	//	fmt.Println(mod)

	ok_mod(t, `
let a = 1;
while (a < 3) {
    a = a + 1;
}`,
		g.MakeInt(3),
		[]*g.Ref{&g.Ref{g.MakeInt(3)}})

	ok_mod(t, `
let a = 1;
while (a < 11) {
    if (a == 4) { a = a + 2; break; }
    a = a + 1;
}`,
		g.MakeInt(6),
		[]*g.Ref{&g.Ref{g.MakeInt(6)}})

	ok_mod(t, `
let a = 1;
let b = 0;
while (a < 11) {
    a = a + 1;
    if (a > 5) { continue; }
    b = b + 1;
}`,
		g.MakeInt(11),
		[]*g.Ref{
			&g.Ref{g.MakeInt(11)},
			&g.Ref{g.MakeInt(4)}})

	ok_mod(t, `
let a = 1;
return a + 2;
let b = 5;`,
		g.MakeInt(3),
		[]*g.Ref{
			&g.Ref{g.MakeInt(1)},
			&g.Ref{g.NULL}})
}

func TestFunc(t *testing.T) {

	source := `
let a = fn() { };
let b = fn(x) { x; };
let c = fn(x, y) { let z = 4; x * y * z; };
let d = a();
let e = b(1);
let f = c(b(2), 3);
`
	mod := newCompiler(source).Compile()
	interpret(mod)
	ok_ref(t, mod.Locals[3], g.NULL)
	ok_ref(t, mod.Locals[4], g.MakeInt(1))
	ok_ref(t, mod.Locals[5], g.MakeInt(24))

	source = `
let fibonacci = fn(n) {
    let x = 0;
    let y = 1;
    let i = 1;
    while i < n {
        let z = x + y;
        x = y;
        y = z;
        i = i + 1;
    }
    return y;
};
let a = fibonacci(1);
let b = fibonacci(2);
let c = fibonacci(3);
let d = fibonacci(4);
let e = fibonacci(5);
let f = fibonacci(6);
`
	mod = newCompiler(source).Compile()
	interpret(mod)
	ok_ref(t, mod.Locals[1], g.MakeInt(1))
	ok_ref(t, mod.Locals[2], g.MakeInt(1))
	ok_ref(t, mod.Locals[3], g.MakeInt(2))
	ok_ref(t, mod.Locals[4], g.MakeInt(3))
	ok_ref(t, mod.Locals[5], g.MakeInt(5))
	ok_ref(t, mod.Locals[6], g.MakeInt(8))

	source = `
let foo = fn(n) {
    let bar = fn(x) {
        return x * (x - 1);
    };
    return bar(n) + bar(n-1);
};
let a = foo(5);
`
	mod = newCompiler(source).Compile()
	interpret(mod)
	ok_ref(t, mod.Locals[1], g.MakeInt(32))
}

func TestCapture(t *testing.T) {

	source := `
const accumGen = fn(n) {
    return fn(i) {
        n = n + i;
        return n;
    };
};
const a = accumGen(3);
let x = a(2);
let y = a(7);
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[2], g.MakeInt(5))
	ok_ref(t, mod.Locals[3], g.MakeInt(12))

	source = `
let z = 2;
const accumGen = fn(n) {
    return fn(i) {
        n = n + i;
        n = n + z;
        return n;
    };
};
const a = accumGen(3);
let x = a(2);
z = 0;
let y = a(1);
`
	mod = newCompiler(source).Compile()

	interpret(mod)

	ok_ref(t, mod.Locals[0], g.MakeInt(0))
	ok_ref(t, mod.Locals[3], g.MakeInt(7))
	ok_ref(t, mod.Locals[4], g.MakeInt(8))

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

}

func newObj(fields map[string]g.Value) g.Obj {
	o := g.NewObj()
	def := &g.ObjDef{[]string{}}
	values := []g.Value{}
	for k, v := range fields {
		def.Keys = append(def.Keys, k)
		values = append(values, v)
	}
	o.Init(def, values)
	return o
}

func TestObj(t *testing.T) {

	source := `
let w = obj {};
let x = obj { a: 0 };
let y = obj { a: 1, b: 2 };
let z = obj { a: 3, b: 4, c: obj { d: 5 } };
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	ok_ref(t, mod.Locals[0], newObj(map[string]g.Value{}))
	ok_ref(t, mod.Locals[1], newObj(map[string]g.Value{"a": g.MakeInt(0)}))
	ok_ref(t, mod.Locals[2], newObj(map[string]g.Value{"a": g.MakeInt(1), "b": g.MakeInt(2)}))
	ok_ref(t, mod.Locals[3],
		newObj(map[string]g.Value{"a": g.MakeInt(3), "b": g.MakeInt(4), "c": newObj(map[string]g.Value{"d": g.MakeInt(5)})}))

	source = `
let x = obj { a: 5 };
let y = x.a;
x.a = 6;
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], newObj(map[string]g.Value{"a": g.MakeInt(6)}))
	ok_ref(t, mod.Locals[1], g.MakeInt(5))

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
	mod = newCompiler(source).Compile()
	interpret(mod)

	ok_ref(t, mod.Locals[2], g.MakeInt(13))
	ok_ref(t, mod.Locals[3], g.MakeInt(3))

	source = `
let a = null;
a = obj { x: 8 }.x = 5;
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	ok_ref(t, mod.Locals[0], g.MakeInt(5))

	source = `
let a = obj { x: 8 };
a['x'] = 3;
let b = a['x']++;
let c = a['x'];
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], newObj(map[string]g.Value{"x": g.MakeInt(4)}))
	ok_ref(t, mod.Locals[1], g.MakeInt(3))
	ok_ref(t, mod.Locals[2], g.MakeInt(4))
}

func TestErrStack(t *testing.T) {

	source := `
let divide = fn(x, y) {
    return x / y;
};
let a = divide(3, 0);
`
	fail(t, source, &ErrorStack{
		g.DivideByZeroError(),
		[]string{
			"    at line 3",
			"    at line 5"}})

	source = `
let foo = fn(n) { n + n; };
let a = foo(5, 6);
	`
	fail(t, source, &ErrorStack{
		g.ArityMismatchError("1", 2),
		[]string{
			"    at line 3"}})
}

func TestPostfix(t *testing.T) {

	source := `
let a = 10;
let b = 20;
let c = a++;
let d = b--;
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	ok_ref(t, mod.Locals[0], g.MakeInt(11))
	ok_ref(t, mod.Locals[1], g.MakeInt(19))
	ok_ref(t, mod.Locals[2], g.MakeInt(10))
	ok_ref(t, mod.Locals[3], g.MakeInt(20))

	source = `
let a = obj { x: 10 };
let b = obj { y: 20 };
let c = a.x++;
let d = b.y--;
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], newObj(map[string]g.Value{"x": g.MakeInt(11)}))
	ok_ref(t, mod.Locals[1], newObj(map[string]g.Value{"y": g.MakeInt(19)}))
	ok_ref(t, mod.Locals[2], g.MakeInt(10))
	ok_ref(t, mod.Locals[3], g.MakeInt(20))
}

func TestTernaryIf(t *testing.T) {

	source := `
let a = true ? 3 : 4;
let b = false ? 5 : 6;
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], g.MakeInt(3))
	ok_ref(t, mod.Locals[1], g.MakeInt(6))
}

func TestList(t *testing.T) {

	source := `
let a = [];
let b = [true];
let c = [false,22];
let d = b[0];
b[0] = 33;
let e = c[1]++;
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], g.NewList([]g.Value{}))
	ok_ref(t, mod.Locals[1], g.NewList([]g.Value{g.MakeInt(33)}))
	ok_ref(t, mod.Locals[2], g.NewList([]g.Value{g.FALSE, g.MakeInt(23)}))
	ok_ref(t, mod.Locals[3], g.TRUE)
	ok_ref(t, mod.Locals[4], g.MakeInt(22))
}

func TestDict(t *testing.T) {

	source := `
let a = dict { 'x': 1, 'y': 2 };
let b = a['x'];
let c = a['z'];
a['x'] = -1;
let d = a['x'];
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[1], g.ONE)
	ok_ref(t, mod.Locals[2], g.NULL)
	ok_ref(t, mod.Locals[3], g.NEG_ONE)
}

func newRange(from int64, to int64, step int64) g.Range {
	r, err := g.NewRange(from, to, step)
	if err != nil {
		panic("invalid range")
	}
	return r
}

func TestBuiltin(t *testing.T) {

	source := `
let a = len([4,5,6]);
let b = str([4,5,6]);
let c = range(0, 5);
let d = range(0, 5, 2);
print();
println();
print(a);
println(b);
print(a,b);
println(a,b);
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], g.MakeInt(3))
	ok_ref(t, mod.Locals[1], g.MakeStr("[ 4, 5, 6 ]"))
	ok_ref(t, mod.Locals[2], newRange(0, 5, 1))
	ok_ref(t, mod.Locals[3], newRange(0, 5, 2))
}

func TestTuple(t *testing.T) {

	source := `
let a = (4,5);
let b = a[0];
let c = a[1];
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], g.NewTuple([]g.Value{g.MakeInt(4), g.MakeInt(5)}))
	ok_ref(t, mod.Locals[1], g.MakeInt(4))
	ok_ref(t, mod.Locals[2], g.MakeInt(5))
}

func TestDecl(t *testing.T) {

	source := `
let a, b = 0;
const c = 1, d;
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], g.NULL)
	ok_ref(t, mod.Locals[1], g.ZERO)
	ok_ref(t, mod.Locals[2], g.ONE)
	ok_ref(t, mod.Locals[3], g.NULL)
}
