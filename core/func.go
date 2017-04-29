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

//---------------------------------------------------------------
// Template

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

//---------------------------------------------------------------
// function

type function struct {
}

func (f *function) TypeOf() (Type, Error) { return TFUNC, nil }

func (f *function) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (f *function) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

//---------------------------------------------------------------
// bytecodeFunc

type bytecodeFunc struct {
	*function
	template *Template
	captures []*Ref
}

// Called via NEW_FUNC opcode at runtime
func NewBytecodeFunc(template *Template) BytecodeFunc {
	captures := make([]*Ref, 0, template.NumCaptures)
	return &bytecodeFunc{&function{}, template, captures}
}

func (bf *bytecodeFunc) ToStr() (Str, Error) {
	return MakeStr(bf.bytecodeStr()), nil
}

func (bf *bytecodeFunc) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *bytecodeFunc:
		if bf.bytecodeStr() == t.bytecodeStr() {
			return TRUE, nil
		} else {
			return FALSE, nil
		}
	default:
		return FALSE, nil
	}
}

func (bf *bytecodeFunc) Add(v Value) (Value, Error) {
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

func (bf *bytecodeFunc) bytecodeStr() string {
	return fmt.Sprintf("func<%p>", bf)
}

//---------------------------------------------------------------
// nativeFunc

type nativeFunc struct {
	*function
}

func (nf *nativeFunc) ToStr() (Str, Error) {
	return MakeStr(nf.nativeStr()), nil
}

func (nf *nativeFunc) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *nativeFunc:
		if nf.nativeStr() == t.nativeStr() {
			return TRUE, nil
		} else {
			return FALSE, nil
		}
	default:
		return FALSE, nil
	}
}

func (nf *nativeFunc) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return Strcat(nf, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (nf *nativeFunc) nativeStr() string {
	return fmt.Sprintf("nativeFunc<%p>", nf)
}

//---------------------------------------------------------------
// nativeFuncs that operate on various Types

type nativeIterNext struct {
	*nativeFunc
	itr Iterator
}

type nativeIterGet struct {
	*nativeFunc
	itr Iterator
}

func (f *nativeIterNext) Invoke(values []Value) (Value, Error) {
	return f.itr.IterNext(), nil
}

func (f *nativeIterGet) Invoke(values []Value) (Value, Error) {
	return f.itr.IterGet()
}
