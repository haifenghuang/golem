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

const (
	PRINT = iota
	PRINTLN
	STR
	LEN
	RANGE
	ASSERT
)

var Builtins = []NativeFunc{
	&fnPrint{&nativeFunc{}},
	&fnPrintln{&nativeFunc{}},
	&fnStr{&nativeFunc{}},
	&fnLen{&nativeFunc{}},
	&fnRange{&nativeFunc{}},
	&fnAssert{&nativeFunc{}}}

type fnPrint struct{ *nativeFunc }
type fnPrintln struct{ *nativeFunc }
type fnStr struct{ *nativeFunc }
type fnLen struct{ *nativeFunc }
type fnRange struct{ *nativeFunc }
type fnAssert struct{ *nativeFunc }

func (fn *fnPrint) Invoke(values []Value) (Value, Error) {
	for _, v := range values {
		s, err := v.ToStr()
		if err != nil {
			return nil, err
		}
		fmt.Print(s.String())
	}

	return NULL, nil
}

func (fn *fnPrintln) Invoke(values []Value) (Value, Error) {
	for _, v := range values {
		s, err := v.ToStr()
		if err != nil {
			return nil, err
		}
		fmt.Print(s.String())
	}
	fmt.Println()

	return NULL, nil
}

func (fn *fnStr) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}

	return values[0].ToStr()
}

func (fn *fnLen) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}

	if ln, ok := values[0].(Lenable); ok {
		return ln.Len()
	} else {
		return nil, TypeMismatchError("Expected Lenable Type")
	}
}

func (fn *fnRange) Invoke(values []Value) (Value, Error) {
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

func (fn *fnAssert) Invoke(values []Value) (Value, Error) {
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
