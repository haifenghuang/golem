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

package analyzer

import (
	"fmt"
	"golem/ast"
	"golem/parser"
	"golem/scanner"
	"testing"
)

func ok(t *testing.T, anl Analyzer, errors []error, dump string) {

	if len(errors) != 0 {
		t.Error(errors)
	}

	if "\n"+ast.Dump(anl.Module()) != dump {
		t.Error("\n"+ast.Dump(anl.Module()), " != ", dump)
	}
}

func fail(t *testing.T, errors []error, expect string) {

	if fmt.Sprintf("%v", errors) != expect {
		t.Error(errors, " != ", expect)
	}
}

func dump(source string) {
	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner)
	mod, err := parser.ParseModule()
	if err != nil {
		panic("analyzer_test: could not parse")
	}
	fmt.Println(ast.Dump(mod))
}

func newAnalyzer(source string) Analyzer {
	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner)
	mod, err := parser.ParseModule()
	if err != nil {
		panic("analyzer_test: could not parse")
	}
	return NewAnalyzer(mod)
}

func TestFlat(t *testing.T) {

	anl := newAnalyzer("let a = 1; const b = 2; a = b + 3;")
	errors := anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   BasicExpr(INT,"1")
.   .   Const
.   .   .   IdentExpr(b,(1,true,false))
.   .   .   BasicExpr(INT,"2")
.   .   Assignment
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   BinaryExpr("+")
.   .   .   .   IdentExpr(b,(1,true,false))
.   .   .   .   BasicExpr(INT,"3")
`)

	errors = newAnalyzer("a;").Analyze()
	fail(t, errors, "[Symbol 'a' is not defined]")

	errors = newAnalyzer("let a = 1;const a = 1;").Analyze()
	fail(t, errors, "[Symbol 'a' is already defined]")

	errors = newAnalyzer("const a = 1;a = 1;").Analyze()
	fail(t, errors, "[Symbol 'a' is constant]")

	errors = newAnalyzer("a = a;").Analyze()
	fail(t, errors, "[Symbol 'a' is not defined Symbol 'a' is not defined]")
}

func TestNested(t *testing.T) {

	source := `
let a = 1; 
if (true) { 
    a = 2; 
    const b = 2; 
} else { 
    a = 3; 
    let b = 3; 
}`
	anl := newAnalyzer(source)
	//errors := anl.Analyze()
	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	errors := anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:3 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   BasicExpr(INT,"1")
.   .   If
.   .   .   BasicExpr(TRUE,"true")
.   .   .   Block
.   .   .   .   Assignment
.   .   .   .   .   IdentExpr(a,(0,false,false))
.   .   .   .   .   BasicExpr(INT,"2")
.   .   .   .   Const
.   .   .   .   .   IdentExpr(b,(1,true,false))
.   .   .   .   .   BasicExpr(INT,"2")
.   .   .   Block
.   .   .   .   Assignment
.   .   .   .   .   IdentExpr(a,(0,false,false))
.   .   .   .   .   BasicExpr(INT,"3")
.   .   .   .   Let
.   .   .   .   .   IdentExpr(b,(2,false,false))
.   .   .   .   .   BasicExpr(INT,"3")
`)
}

func TestWhile(t *testing.T) {

	anl := newAnalyzer("while true { 1 + 2; }")
	errors := anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:0 numCaptures:0 parentCaptures:[])
.   Block
.   .   While
.   .   .   BasicExpr(TRUE,"true")
.   .   .   Block
.   .   .   .   BinaryExpr("+")
.   .   .   .   .   BasicExpr(INT,"1")
.   .   .   .   .   BasicExpr(INT,"2")
`)

	anl = newAnalyzer("while true { 1 + 2; break; continue; }")
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:0 numCaptures:0 parentCaptures:[])
.   Block
.   .   While
.   .   .   BasicExpr(TRUE,"true")
.   .   .   Block
.   .   .   .   BinaryExpr("+")
.   .   .   .   .   BasicExpr(INT,"1")
.   .   .   .   .   BasicExpr(INT,"2")
.   .   .   .   Break
.   .   .   .   Continue
`)

	errors = newAnalyzer("break;").Analyze()
	fail(t, errors, "['break' outside of loop]")

	errors = newAnalyzer("continue;").Analyze()
	fail(t, errors, "['continue' outside of loop]")
}

func TestPureFunction(t *testing.T) {
	source := `
let a = 1;
let b = fn(x) {
    let c = fn(y, z) {
        if (y < 33) {
            return y + z + 5;
        } else {
            let b = 42;
        }
    };
    return c(3);
};`

	anl := newAnalyzer(source)
	errors := anl.Analyze()

	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(0,false,false))
.   .   .   BasicExpr(INT,"1")
.   .   Let
.   .   .   IdentExpr(b,(1,false,false))
.   .   .   FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   .   .   .   IdentExpr(x,(0,false,false))
.   .   .   .   Block
.   .   .   .   .   Let
.   .   .   .   .   .   IdentExpr(c,(1,false,false))
.   .   .   .   .   .   FnExpr(numLocals:3 numCaptures:0 parentCaptures:[])
.   .   .   .   .   .   .   IdentExpr(y,(0,false,false))
.   .   .   .   .   .   .   IdentExpr(z,(1,false,false))
.   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   If
.   .   .   .   .   .   .   .   .   BinaryExpr("<")
.   .   .   .   .   .   .   .   .   .   IdentExpr(y,(0,false,false))
.   .   .   .   .   .   .   .   .   .   BasicExpr(INT,"33")
.   .   .   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   .   .   Return
.   .   .   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   .   .   .   IdentExpr(y,(0,false,false))
.   .   .   .   .   .   .   .   .   .   .   .   .   IdentExpr(z,(1,false,false))
.   .   .   .   .   .   .   .   .   .   .   .   BasicExpr(INT,"5")
.   .   .   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   .   .   Let
.   .   .   .   .   .   .   .   .   .   .   IdentExpr(b,(2,false,false))
.   .   .   .   .   .   .   .   .   .   .   BasicExpr(INT,"42")
.   .   .   .   .   Return
.   .   .   .   .   .   InvokeExpr
.   .   .   .   .   .   .   IdentExpr(c,(1,false,false))
.   .   .   .   .   .   .   BasicExpr(INT,"3")
`)
}

func TestCaptureFunction(t *testing.T) {

	source := `
const accumGen = fn(n) {
    return fn(i) {
        n = n + i;
        return n;
    };
};
`
	anl := newAnalyzer(source)
	errors := anl.Analyze()

	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   Block
.   .   Const
.   .   .   IdentExpr(accumGen,(0,true,false))
.   .   .   FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   .   .   .   IdentExpr(n,(0,false,false))
.   .   .   .   Block
.   .   .   .   .   Return
.   .   .   .   .   .   FnExpr(numLocals:1 numCaptures:1 parentCaptures:[(0,false,false)])
.   .   .   .   .   .   .   IdentExpr(i,(0,false,false))
.   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   Assignment
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   .   IdentExpr(i,(0,false,false))
.   .   .   .   .   .   .   .   Return
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
`)

	source = `
let z = 2;
const accumGen = fn(n) {
    return fn(i) {
        n = n + i;
        n = n + z;
        return n;
    };
};
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()

	//fmt.Println(source)
	//fmt.Println(ast.Dump(anl.Module()))
	//fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(z,(0,false,false))
.   .   .   BasicExpr(INT,"2")
.   .   Const
.   .   .   IdentExpr(accumGen,(1,true,false))
.   .   .   FnExpr(numLocals:1 numCaptures:1 parentCaptures:[(0,false,false)])
.   .   .   .   IdentExpr(n,(0,false,false))
.   .   .   .   Block
.   .   .   .   .   Return
.   .   .   .   .   .   FnExpr(numLocals:1 numCaptures:2 parentCaptures:[(0,false,false), (0,false,true)])
.   .   .   .   .   .   .   IdentExpr(i,(0,false,false))
.   .   .   .   .   .   .   Block
.   .   .   .   .   .   .   .   Assignment
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   .   IdentExpr(i,(0,false,false))
.   .   .   .   .   .   .   .   Assignment
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
.   .   .   .   .   .   .   .   .   .   IdentExpr(z,(1,false,true))
.   .   .   .   .   .   .   .   Return
.   .   .   .   .   .   .   .   .   IdentExpr(n,(0,false,true))
`)
}

func TestObj(t *testing.T) {

	errors := newAnalyzer("this;").Analyze()
	fail(t, errors, "['this' outside of loop]")

	source := `
obj{ };
`
	anl := newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:0 numCaptures:0 parentCaptures:[])
.   Block
.   .   ObjExpr([],-1)
`)

	source = `
obj{ a: 1 };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:0 numCaptures:0 parentCaptures:[])
.   Block
.   .   ObjExpr([a],-1)
.   .   .   BasicExpr(INT,"1")
`)

	source = `
obj{ a: this };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   Block
.   .   ObjExpr([a],0)
.   .   .   ThisExpr((0,true,false))
`)

	source = `
obj{ a: obj { b: this } };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   Block
.   .   ObjExpr([a],-1)
.   .   .   ObjExpr([b],0)
.   .   .   .   ThisExpr((0,true,false))
`)

	source = `
obj{ a: obj { b: 1 }, c: this.a };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:1 numCaptures:0 parentCaptures:[])
.   Block
.   .   ObjExpr([a, c],0)
.   .   .   ObjExpr([b],-1)
.   .   .   .   BasicExpr(INT,"1")
.   .   .   FieldExpr(a)
.   .   .   .   ThisExpr((0,true,false))
`)

	source = `
obj{ a: obj { b: this }, c: this };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   ObjExpr([a, c],1)
.   .   .   ObjExpr([b],0)
.   .   .   .   ThisExpr((0,true,false))
.   .   .   ThisExpr((1,true,false))
`)

	source = `
obj{ a: this, b: obj { c: this } };
`
	anl = newAnalyzer(source)
	errors = anl.Analyze()
	ok(t, anl, errors, `
FnExpr(numLocals:2 numCaptures:0 parentCaptures:[])
.   Block
.   .   ObjExpr([a, b],0)
.   .   .   ThisExpr((0,true,false))
.   .   .   ObjExpr([c],1)
.   .   .   .   ThisExpr((1,true,false))
`)

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
	errors = anl.Analyze()

	//	fmt.Println(source)
	//	fmt.Println(ast.Dump(anl.Module()))
	//	fmt.Println(errors)

	ok(t, anl, errors, `
FnExpr(numLocals:4 numCaptures:0 parentCaptures:[])
.   Block
.   .   Let
.   .   .   IdentExpr(a,(1,false,false))
.   .   .   ObjExpr([x, y, plus, minus],0)
.   .   .   .   BasicExpr(INT,"8")
.   .   .   .   BasicExpr(INT,"5")
.   .   .   .   FnExpr(numLocals:0 numCaptures:1 parentCaptures:[(0,true,false)])
.   .   .   .   .   Block
.   .   .   .   .   .   Return
.   .   .   .   .   .   .   BinaryExpr("+")
.   .   .   .   .   .   .   .   FieldExpr(x)
.   .   .   .   .   .   .   .   .   ThisExpr((0,true,true))
.   .   .   .   .   .   .   .   FieldExpr(y)
.   .   .   .   .   .   .   .   .   ThisExpr((0,true,true))
.   .   .   .   FnExpr(numLocals:0 numCaptures:1 parentCaptures:[(0,true,false)])
.   .   .   .   .   Block
.   .   .   .   .   .   Return
.   .   .   .   .   .   .   BinaryExpr("-")
.   .   .   .   .   .   .   .   FieldExpr(x)
.   .   .   .   .   .   .   .   .   ThisExpr((0,true,true))
.   .   .   .   .   .   .   .   FieldExpr(y)
.   .   .   .   .   .   .   .   .   ThisExpr((0,true,true))
.   .   Let
.   .   .   IdentExpr(b,(2,false,false))
.   .   .   InvokeExpr
.   .   .   .   FieldExpr(plus)
.   .   .   .   .   IdentExpr(a,(1,false,false))
.   .   Let
.   .   .   IdentExpr(c,(3,false,false))
.   .   .   InvokeExpr
.   .   .   .   FieldExpr(minus)
.   .   .   .   .   IdentExpr(a,(1,false,false))
`)
}
