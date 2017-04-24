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
	"fmt"
)

//---------------------------------------------------------------
// Func represents an instance of a function

type _func struct {
}

func (f *_func) TypeOf() (Type, Error) { return TFUNC, nil }

func (f *_func) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *_func) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

//---------------------------------------------------------------
// BytecodeFunc represents a function that is defined
// via Golem source code

type _bytecodeFunc struct {
	*_func
	template *Template
	captures []*Ref
}

// Called via NEW_FUNC opcode at runtime
func NewBytecodeFunc(template *Template) BytecodeFunc {
	captures := make([]*Ref, 0, template.NumCaptures)
	return &_bytecodeFunc{&_func{}, template, captures}
}

func (bf *_bytecodeFunc) ToStr() (Str, Error) {
	return MakeStr(bf.bytecodeStr()), nil
}

func (bf *_bytecodeFunc) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *_bytecodeFunc:
		if bf.bytecodeStr() == t.bytecodeStr() {
			return TRUE, nil
		} else {
			return FALSE, nil
		}
	default:
		return FALSE, nil
	}
}

func (bf *_bytecodeFunc) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(bf, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (bf *_bytecodeFunc) Template() *Template {
	return bf.template
}

func (bf *_bytecodeFunc) GetCapture(idx int) *Ref {
	return bf.captures[idx]
}

func (bf *_bytecodeFunc) PushCapture(ref *Ref) {
	bf.captures = append(bf.captures, ref)
}

func (bf *_bytecodeFunc) bytecodeStr() string {
	return fmt.Sprintf("func(%p)", bf)
}

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

// Return the line number for the opcode at the given instruction pointer
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

// OpcLine tracks which sequence of opcodes are on a given line
type OpcLine struct {
	Index   int
	LineNum int
}

//---------------------------------------------------------------
// NativeFunc represents a function that is defined
// natively within Go.

type _nativeFunc struct {
	*_func
}

func (nf *_nativeFunc) ToStr() (Str, Error) {
	return MakeStr(nf.nativeStr()), nil
}

func (nf *_nativeFunc) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *_nativeFunc:
		if nf.nativeStr() == t.nativeStr() {
			return TRUE, nil
		} else {
			return FALSE, nil
		}
	default:
		return FALSE, nil
	}
}

func (nf *_nativeFunc) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat(nf, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (nf *_nativeFunc) nativeStr() string {
	return fmt.Sprintf("nativeFunc(%p)", nf)
}
