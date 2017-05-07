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

// NativeFunc represents a function that is defined
// natively within Go.
type NativeFunc interface {
	Func
	Invoke([]Value) (Value, Error)
}

// NOTE: 'nativeFunc' cannot be an empty struct, because empty structs have
// unusual semantics in Go, i.e. they all point to the same address.
//
// https://golang.org/ref/spec#Size_and_alignment_guarantees
//
// To work around that, we place an arbitrary value inside the struct, so
// that it wont be empty.
//
type nativeFunc struct {
	placeholder int
}

func (nf *nativeFunc) TypeOf() Type { return TFUNC }

func (nf *nativeFunc) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (nf *nativeFunc) GetField(key Str) (Value, Error) {
	return nil, NoSuchFieldError(key.String())
}

func (nf *nativeFunc) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (nf *nativeFunc) Eq(v Value) Bool {
	switch t := v.(type) {
	case NativeFunc:
		return nf.ToStr().Eq(t.ToStr())
	default:
		return FALSE
	}
}

func (nf *nativeFunc) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return Strcat(nf, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (nf *nativeFunc) ToStr() Str {
	return MakeStr(fmt.Sprintf("native<%p>", nf))
}

//---------------------------------------------------------------
// nativeFuncs that support iteration

type nativeIterNext struct {
	nativeFunc
	itr Iterator
}

type nativeIterGet struct {
	nativeFunc
	itr Iterator
}

func (fn *nativeIterNext) Invoke(values []Value) (Value, Error) {
	return fn.itr.IterNext(), nil
}

func (fn *nativeIterGet) Invoke(values []Value) (Value, Error) {
	return fn.itr.IterGet()
}

//---------------------------------------------------------------
// Builtins

const (
	PRINT = iota
	PRINTLN
	STR
	LEN
	RANGE
	ASSERT
)

type (
	nativePrint   struct{ *nativeFunc }
	nativePrintln struct{ *nativeFunc }
	nativeStr     struct{ *nativeFunc }
	nativeLen     struct{ *nativeFunc }
	nativeRange   struct{ *nativeFunc }
	nativeAssert  struct{ *nativeFunc }
)

var Builtins = []NativeFunc{
	&nativePrint{&nativeFunc{}},
	&nativePrintln{&nativeFunc{}},
	&nativeStr{&nativeFunc{}},
	&nativeLen{&nativeFunc{}},
	&nativeRange{&nativeFunc{}},
	&nativeAssert{&nativeFunc{}}}

func (fn *nativePrint) Invoke(values []Value) (Value, Error) {
	for _, v := range values {
		fmt.Print(v.ToStr().String())
	}

	return NULL, nil
}

func (fn *nativePrintln) Invoke(values []Value) (Value, Error) {
	for _, v := range values {
		fmt.Print(v.ToStr().String())
	}
	fmt.Println()

	return NULL, nil
}

func (fn *nativeStr) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}

	return values[0].ToStr(), nil
}

func (fn *nativeLen) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}

	if ln, ok := values[0].(Lenable); ok {
		return ln.Len(), nil
	} else {
		return nil, TypeMismatchError("Expected Lenable Type")
	}
}

func (fn *nativeRange) Invoke(values []Value) (Value, Error) {
	if len(values) < 2 || len(values) > 3 {
		return nil, ArityMismatchError("2 or 3", len(values))
	}

	from, ok := values[0].(Int)
	if !ok {
		return nil, TypeMismatchError("Expected 'Int'")
	}

	to, ok := values[1].(Int)
	if !ok {
		return nil, TypeMismatchError("Expected 'Int'")
	}

	step := ONE
	if len(values) == 3 {
		step, ok = values[2].(Int)
		if !ok {
			return nil, TypeMismatchError("Expected 'Int'")
		}
	}

	return NewRange(from.IntVal(), to.IntVal(), step.IntVal())
}

func (fn *nativeAssert) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}

	b, ok := values[0].(Bool)
	if !ok {
		return nil, TypeMismatchError("Expected 'Bool'")
	}

	if b.BoolVal() {
		return TRUE, nil
	} else {
		return nil, AssertionFailedError()
	}
}
