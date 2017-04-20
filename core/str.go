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
	"bytes"
)

type Str string

func (s Str) TypeOf() (Type, Error) { return TSTR, nil }

func (s Str) String() (Str, Error) { return s, nil }

func (s Str) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {

	case Str:
		return MakeBool(s == t), nil

	default:
		return FALSE, nil
	}
}

func (s Str) Cmp(v Value) (Int, Error) {
	switch t := v.(type) {

	case Str:
		if s < t {
			return -1, nil
		} else if s > t {
			return 1, nil
		} else {
			return 0, nil
		}

	default:
		return 0, TypeMismatchError("Expected Comparable Type")
	}
}

func (s Str) Add(v Value) (Value, Error) {
	return strcat([]Value{s, v})
}

func (s Str) Sub(v Value) (Number, Error)    { return Int(0), TypeMismatchError("Expected Number Type") }
func (s Str) Mul(v Value) (Number, Error)    { return Int(0), TypeMismatchError("Expected Number Type") }
func (s Str) Div(v Value) (Number, Error)    { return Int(0), TypeMismatchError("Expected Number Type") }
func (s Str) Rem(v Value) (Int, Error)       { return Int(0), TypeMismatchError("Expected 'Int'") }
func (s Str) BitAnd(v Value) (Int, Error)    { return Int(0), TypeMismatchError("Expected 'Int'") }
func (s Str) BitOr(v Value) (Int, Error)     { return Int(0), TypeMismatchError("Expected 'Int'") }
func (s Str) BitXOr(v Value) (Int, Error)    { return Int(0), TypeMismatchError("Expected 'Int'") }
func (s Str) LeftShift(v Value) (Int, Error) { return Int(0), TypeMismatchError("Expected 'Int'") }
func (s Str) RightShift(Value) (Int, Error)  { return Int(0), TypeMismatchError("Expected 'Int'") }

func (s Str) Negate() (Number, Error)  { return Int(0), TypeMismatchError("Expected Number Type") }
func (s Str) Complement() (Int, Error) { return Int(0), TypeMismatchError("Expected 'Int'") }

//-----------------------------------

func strcat(a []Value) (Str, Error) {
	var buf bytes.Buffer
	for _, v := range a {
		s, err := v.String()
		if err != nil {
			return Str(""), err
		}
		buf.WriteString(string(s))
	}
	return Str(buf.String()), nil
}

func (s Str) GetField(key string) (Value, Error)   { return nil, TypeMismatchError("Expected 'Obj'") }
func (s Str) PutField(key string, val Value) Error { return TypeMismatchError("Expected 'Obj'") }
