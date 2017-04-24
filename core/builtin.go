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
)

var Builtins = []NativeFunc{
	&_print{&_nativeFunc{}},
	&_println{&_nativeFunc{}},
	&_str{&_nativeFunc{}},
	&_len{&_nativeFunc{}}}

type _print struct{ *_nativeFunc }
type _println struct{ *_nativeFunc }
type _str struct{ *_nativeFunc }
type _len struct{ *_nativeFunc }

func (builtin *_print) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError(1, len(values))
	}

	s, err := values[0].ToStr()
	if err != nil {
		return nil, err
	}

	fmt.Print(s.String())
	return NULL, nil
}

func (builtin *_println) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError(1, len(values))
	}

	s, err := values[0].ToStr()
	if err != nil {
		return nil, err
	}

	fmt.Println(s.String())
	return NULL, nil
}

func (builtin *_str) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError(1, len(values))
	}

	return values[0].ToStr()
}

func (builtin *_len) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError(1, len(values))
	}

	if ln, ok := values[0].(Lenable); ok {
		return ln.Len()
	} else {
		return nil, TypeMismatchError("Expected Lenable Type")
	}
}
