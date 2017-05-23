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
	"fmt"
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

	result, errTrace := intp.Init()
	if errTrace != nil {
		panic(errTrace)
	}

	b := result.Eq(expect)
	if !b.BoolVal() {
		t.Error(result, " != ", expect)
		panic("ok_expr")
	}
}

func ok_ref(t *testing.T, ref *g.Ref, expect g.Value) {
	b := ref.Val.Eq(expect)
	if !b.BoolVal() {
		t.Error(ref.Val, " != ", expect)
	}
}

func ok_mod(t *testing.T, source string, expectResult g.Value, expectLocals []*g.Ref) {
	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod)

	result, errTrace := intp.Init()
	if errTrace != nil {
		panic(errTrace)
	}

	b := result.Eq(expectResult)
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

	result, errTrace := intp.Init()
	if result != nil {
		panic(result)
	}

	if errTrace.Error.Error() != expect {
		t.Error(errTrace.Error.Error(), " != ", expect)
	}
}

func fail(t *testing.T, source string, expectErr g.Error, expectErrTrace []string) *g.Module {

	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod)

	result, errTrace := intp.Init()
	if result != nil {
		panic(result)
	}

	if !reflect.DeepEqual(errTrace.Error, expectErr) {
		t.Error(errTrace.Error, " != ", expectErr)
	}

	if !reflect.DeepEqual(errTrace.StackTrace, expectErrTrace) {
		t.Error(errTrace.StackTrace, " != ", expectErrTrace)
	}

	return mod
}

func failErr(t *testing.T, source string, expect g.Error) {

	mod := newCompiler(source).Compile()
	intp := NewInterpreter(mod)

	result, errTrace := intp.Init()
	if result != nil {
		panic(result)
	}

	if errTrace.Error.Error() != expect.Error() {
		t.Error(errTrace.Error, " != ", expect)
	}
}

func newStruct(entries []*g.StructEntry) g.Struct {
	stc, err := g.NewStruct(entries)
	if err != nil {
		panic("invalid struct")
	}
	return stc
}

func newCompiler(source string) compiler.Compiler {
	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner)
	mod, err := parser.ParseModule()
	if err != nil {
		panic(err.Error())
	}
	anl := analyzer.NewAnalyzer(mod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		panic(fmt.Sprintf("%v", errors))
	}

	return compiler.NewCompiler(anl)
}

func interpret(mod *g.Module) {
	intp := NewInterpreter(mod)
	_, errTrace := intp.Init()
	if errTrace != nil {
		fmt.Printf("%v\n", errTrace.Error)
		fmt.Printf("%v\n", errTrace.StackTrace)
		panic("interpreter failed")
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
	ok_expr(t, "2 <=> 2;", g.ZERO)
	ok_expr(t, "2 <=> 1;", g.ONE)

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

	ok_expr(t, "struct{a: 1}['a'];", g.ONE)
	ok_expr(t, "struct{a: 1} has 'a';", g.TRUE)
	ok_expr(t, "struct{a: 1} has 'b';", g.FALSE)

	fail_expr(t, "struct{a: 1}[0];", "TypeMismatch: Expected 'Str'")

	fail_expr(t, "struct{a: 1, a: 2};", "DuplicateField: Field 'a' is a duplicate")

	ok_expr(t, "struct{} == struct{};", g.TRUE)
	ok_expr(t, "struct{a:1} == struct{a:1};", g.TRUE)
	ok_expr(t, "struct{a:1,b:2} == struct{a:1,b:2};", g.TRUE)
	ok_expr(t, "struct{a:1} != struct{a:1,b:2};", g.TRUE)
	ok_expr(t, "struct{a:1,b:2} != struct{b:2};", g.TRUE)
	ok_expr(t, "struct{a:1,b:2} != struct{a:3,b:2};", g.TRUE)
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
		[]*g.Ref{&g.Ref{g.ONE}})

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
			&g.Ref{g.ONE},
			&g.Ref{g.NULL}})
}

func TestFunc(t *testing.T) {

	source := `
let a = fn(x) { x; };
let b = a(1);
`
	mod := newCompiler(source).Compile()

	interpret(mod)
	ok_ref(t, mod.Locals[1], g.ONE)

	source = `
let a = fn() { };
let b = fn(x) { x; };
let c = fn(x, y) { let z = 4; x * y * z; };
let d = a();
let e = b(1);
let f = c(b(2), 3);
`
	mod = newCompiler(source).Compile()

	interpret(mod)
	ok_ref(t, mod.Locals[3], g.NULL)
	ok_ref(t, mod.Locals[4], g.ONE)
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
	ok_ref(t, mod.Locals[1], g.ONE)
	ok_ref(t, mod.Locals[2], g.ONE)
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

	ok_ref(t, mod.Locals[0], g.ZERO)
	ok_ref(t, mod.Locals[3], g.MakeInt(7))
	ok_ref(t, mod.Locals[4], g.MakeInt(8))

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

}

func TestStruct(t *testing.T) {

	source := `
let w = struct {};
let x = struct { a: 0 };
let y = struct { a: 1, b: 2 };
let z = struct { a: 3, b: 4, c: struct { d: 5 } };
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	ok_ref(t, mod.Locals[0], newStruct([]*g.StructEntry{}))
	ok_ref(t, mod.Locals[1], newStruct([]*g.StructEntry{{"a", g.ZERO}}))
	ok_ref(t, mod.Locals[2], newStruct([]*g.StructEntry{{"a", g.ONE}, {"b", g.MakeInt(2)}}))
	ok_ref(t, mod.Locals[3], newStruct([]*g.StructEntry{{"a", g.MakeInt(3)}, {"b", g.MakeInt(4)},
		{"c", newStruct([]*g.StructEntry{{"d", g.MakeInt(5)}})}}))

	source = `
let x = struct { a: 5 };
let y = x.a;
x.a = 6;
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], newStruct([]*g.StructEntry{{"a", g.MakeInt(6)}}))
	ok_ref(t, mod.Locals[1], g.MakeInt(5))

	source = `
let a = struct {
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
a = struct { x: 8 }.x = 5;
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	ok_ref(t, mod.Locals[0], g.MakeInt(5))

	source = `
let a = struct { x: 8 };
a['x'] = 3;
let b = a['x']++;
let c = a['x'];
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], newStruct([]*g.StructEntry{{"x", g.MakeInt(4)}}))
	ok_ref(t, mod.Locals[1], g.MakeInt(3))
	ok_ref(t, mod.Locals[2], g.MakeInt(4))

	source = `
let a = struct { x: 8 };
assert(a has 'x');
assert(!(a has 'z'));
assert(a has 'x');
let b = struct { x: this has 'x', y: this has 'z' };
assert(b.x);
assert(!b.y);
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestChain(t *testing.T) {

	fail_expr(t, "struct (1) { y: 3, z: 4};", "TypeMismatch: Expected 'Struct'")

	source := `
let a = struct { x: 1, y: 2};
let b = struct (a) { y: 3, z: 4};
assert(b.x == 1);
assert(b.y == 3);
assert(b.z == 4);
a.x = 5;
a.y = 6;
assert(b.x == 5);
assert(b.y == 3);
assert(b.z == 4);
let c = struct (b) { w: 10};
assert(c.w == 10);
assert(c.x == 5);
assert(c.y == 3);
assert(c.z == 4);
a.x = 7;
b.z = 11;
assert(c.w == 10);
assert(c.x == 7);
assert(c.y == 3);
assert(c.z == 11);
`
	mod := newCompiler(source).Compile()
	interpret(mod)
}

func TestErrStack(t *testing.T) {

	source := `
let divide = fn(x, y) {
    return x / y;
};
let a = divide(3, 0);
`
	fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 3",
			"    at line 5"})

	source = `
let foo = fn(n) { n + n; };
let a = foo(5, 6);
	`
	fail(t, source,
		g.ArityMismatchError("1", 2),
		[]string{
			"    at line 3"})
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
let a = struct { x: 10 };
let b = struct { y: 20 };
let c = a.x++;
let d = b.y--;
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	ok_ref(t, mod.Locals[0], newStruct([]*g.StructEntry{{"x", g.MakeInt(11)}}))
	ok_ref(t, mod.Locals[1], newStruct([]*g.StructEntry{{"y", g.MakeInt(19)}}))
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

	source = `
let a = [];
a.add(1);
assert(a == [1]);
a.add(2).add([3]);
assert(a == [1,2,[3]]);
let b = [];
b.add(4);
assert(b == [4]);
assert(a.add == a.add);
assert(b.add == b.add);
assert(a.add != b.add);
assert(b.add != a.add);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = [];
a.addAll([1,2]).addAll('bc');
assert(a == [1,2,'b','c']);
let b = [];
b.addAll(range(0,3));
b.addAll(dict { 'x': 1, 'y': 2 });
assert(b == [ 0, 1, 2, ('x', 1), ('y', 2)]);
assert(a.addAll == a.addAll);
assert(b.addAll == b.addAll);
assert(a.addAll != b.addAll);
assert(b.addAll != a.addAll);
assert(a.add != a.addAll);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "let a = []; a.addAll(false);"
	failErr(t, source, g.TypeMismatchError("Expected Iterable Type"))

	source = "let a = []; a.add(3,4);"
	failErr(t, source, g.ArityMismatchError("1", 2))

	source = `
let a = [];
assert(a.isEmpty());
a.add(1);
assert(!a.isEmpty());
a.clear();
assert(a.isEmpty());
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = [];
assert(!a.contains('x'));
assert(a.indexOf('x') == -1);
a = ['z', 'x'];
assert(a.contains('x'));
assert(a.indexOf('x') == 1);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = [];
assert(a.join() == '');
assert(a.join(',') == '');
a.add(1);
assert(a.join() == '1');
assert(a.join(',') == '1');
a.add(2);
assert(a.join() == '12');
assert(a.join(',') == '1,2');
a.add('abc');
assert(a.join() == '12abc');
assert(a.join(',') == '1,2,abc');
`
	mod = newCompiler(source).Compile()
	interpret(mod)
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

	source = `
let a = dict {};
a.addAll([(1,2)]).addAll([(3,4)]);
assert(a == dict {1:2,3:4});
let b = dict {};
assert(a.addAll == a.addAll);
assert(b.addAll == b.addAll);
assert(a.addAll != b.addAll);
assert(b.addAll != a.addAll);
assert(a.clear != a.addAll);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "let a = dict{}; a.addAll(false);"
	failErr(t, source, g.TypeMismatchError("Expected Iterable Type"))
	source = "let a = dict{}; a.addAll([false]);"
	failErr(t, source, g.TypeMismatchError("Expected Tuple"))
	source = "let a = dict{}; a.addAll([(1,2,3)]);"
	failErr(t, source, g.TupleLengthError(2, 3))

	source = "let a = dict{}; a[[1,2]];"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = "let a = dict{}; a[[1,2]] = 3;"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = "let a = dict{}; a.containsKey([1,2]);"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = `
let a = dict {};
assert(a.isEmpty());
a[1] = 2;
assert(!a.isEmpty());
a.clear();
assert(a.isEmpty());
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = dict {'z': 3};
assert(a.containsKey('z'));
assert(!a.containsKey('x'));
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestSet(t *testing.T) {

	source := `
let a = set {};
a.add(1);
assert(a == set {1});
a.add(2).add(3).add(2);
assert(a == set {1,2,3});
let b = set { 4 };
b.add(4);
assert(b == set { 4 });
assert(a.add == a.add);
assert(b.add == b.add);
assert(a.add != b.add);
assert(b.add != a.add);
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = set {};
a.addAll([1,2]).addAll('bc');
assert(a == set {1,2,'b','c'});
let b = set {};
b.addAll(range(0,3));
assert(b == set { 0, 1, 2 });
assert(a.addAll == a.addAll);
assert(b.addAll == b.addAll);
assert(a.addAll != b.addAll);
assert(b.addAll != a.addAll);
assert(a.add != a.addAll);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "let a = set{}; a.addAll(false);"
	failErr(t, source, g.TypeMismatchError("Expected Iterable Type"))

	source = "let a = set{}; a.add(3,4);"
	failErr(t, source, g.ArityMismatchError("1", 2))

	source = "let a = set{}; a.add([1,2]);"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = "let a = set{}; a.contains([1,2]);"
	failErr(t, source, g.TypeMismatchError("Expected Hashable Type"))

	source = `
let a = set{};
assert(a.isEmpty());
a.add(1);
assert(!a.isEmpty());
a.clear();
assert(a.isEmpty());
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = set{};
assert(!a.contains('x'));
a = set {'z', 'x'};
assert(a.contains('x'));
`
	mod = newCompiler(source).Compile()
	interpret(mod)
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
assert(print == print);
assert(print != println);
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

	source = `
let a = assert(true);
`
	mod = newCompiler(source).Compile()
	interpret(mod)
	ok_ref(t, mod.Locals[0], g.TRUE)

	fail(t, "assert(1, 2);",
		g.ArityMismatchError("1", 2),
		[]string{
			"    at line 1"})

	fail(t, "assert(1);",
		g.TypeMismatchError("Expected 'Bool'"),
		[]string{
			"    at line 1"})

	fail(t, "assert(1 == 2);",
		g.AssertionFailedError(),
		[]string{
			"    at line 1"})
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

func TestFor(t *testing.T) {

	source := `
let a = 0;
for n in [1,2,3] {
    a += n;
}
assert(a == 6);
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
let keys = '';
let values = 0;
for (k, v)  in dict {'a': 1, 'b': 2, 'c': 3} {
    keys += k;
    values += v;
}
assert(keys == 'bac');
assert(values == 6);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let entries = '';
for e in dict {'a': 1, 'b': 2, 'c': 3} {
    entries += str(e);
}
assert(entries == '(b, 2)(a, 1)(c, 3)');
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let keys = '';
let values = 0;
for (k, v)  in [('a', 1), ('b', 2), ('c', 3)] {
    keys += k;
    values += v;
}
assert(keys == 'abc');
assert(values == 6);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = "for (k, v)  in [1, 2, 3] {}"
	fail(t, source,
		g.TypeMismatchError("Expected 'Tuple'"),
		[]string{"    at line 1"})

	source = "for (a, b, c)  in [('a', 1), ('b', 2), ('c', 3)] {}"
	fail(t, source,
		g.InvalidArgumentError("Expected Tuple of length 3"),
		[]string{"    at line 1"})
}

func TestSwitch(t *testing.T) {

	source := `
let s = '';
for i in range(0, 4) {
    switch {
    case i == 0:
        s += 'a';

    case i == 1, i == 2:
        s += 'b';

    default:
        s += 'c';
    }
}
assert(s == 'abbc');
`
	mod := newCompiler(source).Compile()
	interpret(mod)

	source = `
let s = '';
for i in range(0, 4) {
    switch {
    case i == 0, i == 1:
        s += 'a';

    case i == 2:
        s += 'b';
    }
}
assert(s == 'aab');
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let s = '';
for i in range(0, 4) {
    switch i {
    case 0, 1:
        s += 'a';

    case 2:
        s += 'b';
    }
}
assert(s == 'aab');
`
	mod = newCompiler(source).Compile()
	interpret(mod)
}

func TestGetField(t *testing.T) {

	source := "null.bogus;"
	fail(t, source,
		g.NullValueError(),
		[]string{"    at line 1"})

	err := g.NoSuchFieldError("bogus")

	failErr(t, "true.bogus;", err)
	failErr(t, "'a'.bogus;", err)
	failErr(t, "(0).bogus;", err)
	failErr(t, "(0.123).bogus;", err)

	failErr(t, "(1,2).bogus;", err)
	failErr(t, "range(1,2).bogus;", err)
	failErr(t, "[1,2].bogus;", err)
	failErr(t, "dict {'a':1}.bogus;", err)
	failErr(t, "struct {a:1}.bogus;", err)

	failErr(t, "(fn() {}).bogus;", err)
}

func TestFinally(t *testing.T) {

	source := `
let a = 1;
try {
    3 / 0;
} finally {
    a = 2;
}
try {
    3 / 0;
} finally {
    a = 3;
}
`
	mod := fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 4"})
	ok_ref(t, mod.Locals[0], g.MakeInt(2))

	source = `
let a = 1;
try {
    try {
        3 / 0;
    } finally {
        a++;
    }
} finally {
    a++;
}
`
	mod = fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 5"})
	ok_ref(t, mod.Locals[0], g.MakeInt(3))

	source = `
let a = 1;
let b = fn() { a++; };
try {
    try {
        3 / 0;
    } finally {
        a++;
        b();
    }
} finally {
    a++;
}
`
	mod = fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 6"})
	ok_ref(t, mod.Locals[0], g.MakeInt(4))

	source = `
let a = 1;
let b = fn() { 
    try {
        try {
            3 / 0;
        } finally {
            a++;
        }
    } finally {
        a++;
    }
};
try {
    b();
} finally {
    a++;
}
`
	//mod = newCompiler(source).Compile()
	//fmt.Println("----------------------------")
	//fmt.Println(source)
	//fmt.Println(mod)

	mod = fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 6",
			"    at line 15"})
	ok_ref(t, mod.Locals[0], g.MakeInt(4))

	source = `
let b = fn() { 
    try {
    } finally {
        return 1;
    }
    return 2;
};
assert(b() == 1);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
let a = 1;
let b = fn() { 
    try {
        try {
        } finally {
            return 1;
        }
    } finally {
        a = 2;
    }
};
assert(b() == 1);
assert(a == 1);
`
	mod = newCompiler(source).Compile()
	interpret(mod)

	source = `
try {
    assert(1,2,3);
} finally {
}
`
	mod = fail(t, source,
		g.ArityMismatchError("1", 3),
		[]string{
			"    at line 3"})

	source = `
try {
    assert(1,2,3);
} finally {
    1/0;
}
`
	mod = fail(t, source,
		g.DivideByZeroError(),
		[]string{
			"    at line 5"})
}
