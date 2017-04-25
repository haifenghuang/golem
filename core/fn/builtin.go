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

package fn

import (
	"fmt"
	g "golem/core"
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

func (builtin *_print) Invoke(values []g.Value) (g.Value, g.Error) {
	for _, v := range values {
		s, err := v.ToStr()
		if err != nil {
			return nil, err
		}
		fmt.Print(s.String())
	}

	return g.NULL, nil
}

func (builtin *_println) Invoke(values []g.Value) (g.Value, g.Error) {
	for _, v := range values {
		s, err := v.ToStr()
		if err != nil {
			return nil, err
		}
		fmt.Print(s.String())
	}
	fmt.Println()

	return g.NULL, nil
}

func (builtin *_str) Invoke(values []g.Value) (g.Value, g.Error) {
	if len(values) != 1 {
		return nil, g.ArityMismatchError(1, len(values))
	}

	return values[0].ToStr()
}

func (builtin *_len) Invoke(values []g.Value) (g.Value, g.Error) {
	if len(values) != 1 {
		return nil, g.ArityMismatchError(1, len(values))
	}

	if ln, ok := values[0].(g.Lenable); ok {
		return ln.Len()
	} else {
		return nil, g.TypeMismatchError("Expected Lenable Type")
	}
}
