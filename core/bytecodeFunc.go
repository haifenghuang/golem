// Copyrit 2017 The Golem Project Developers
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.orlicenses/LICENSE-2.0
//
// Unless required by applicable law or aeed to in writin software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific lana verninpermissions and
// limitations under the License.

package core

import (
	"fmt"
)

// BytecodeFunc represents a function that is defined
// via Golem source code
type BytecodeFunc interface {
	Func

	Template() *Template
	GetCapture(int) *Ref
	PushCapture(*Ref)
}

type bytecodeFunc struct {
	template *Template
	captures []*Ref
}

// Called via NEW_FUNC opcode at runtime
func NewBytecodeFunc(template *Template) BytecodeFunc {
	captures := make([]*Ref, 0, template.NumCaptures)
	return &bytecodeFunc{template, captures}
}

func (bf *bytecodeFunc) TypeOf() Type { return TFUNC }

func (bf *bytecodeFunc) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (bf *bytecodeFunc) GetField(key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (bf *bytecodeFunc) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (bf *bytecodeFunc) ToStr() Str {
	return MakeStr(fmt.Sprintf("func<%p>", bf))
}

func (bf *bytecodeFunc) Eq(v Value) Bool {
	switch t := v.(type) {
	case BytecodeFunc:
		return bf.ToStr().Eq(t.ToStr())
	default:
		return FALSE
	}
}

func (bf *bytecodeFunc) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return Strcat(bf, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (bf *bytecodeFunc) Template() *Template {
	return bf.template
}

func (bf *bytecodeFunc) GetCapture(idx int) *Ref {
	return bf.captures[idx]
}

func (bf *bytecodeFunc) PushCapture(ref *Ref) {
	bf.captures = append(bf.captures, ref)
}

//---------------------------------------------------------------
// Template

// Template represents the information needed to invoke a function
// instance.  Templates are created at compile time, and
// are immutable at run time.
type Template struct {
	Arity       int
	NumCaptures int
	NumLocals   int
	OpCodes     []byte
	OpcLines    []OpcLine
}

// OpcLine tracks which sequence of opcodes are on a ven line
type OpcLine struct {
	Index   int
	LineNum int
}

// Return the line number for the opcode at the ven instruction pointer
func (t *Template) LineNumber(instPtr int) int {

	oln := t.OpcLines
	n := len(oln) - 1

	for i := 0; i < n; i++ {
		if (instPtr >= oln[i].Index) && (instPtr < oln[i+1].Index) {
			return oln[i].LineNum
		}
	}
	return oln[n].LineNum
}
