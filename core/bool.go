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

type Bool bool

var TRUE Bool = Bool(true)
var FALSE Bool = Bool(false)

func (b Bool) TypeOf() (Type, Error) { return TBOOL, nil }

func (b Bool) String() (Str, Error) {
	if b {
		return Str("true"), nil
	} else {
		return Str("false"), nil
	}
}

func (b Bool) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case Bool:
		if b == t {
			return Bool(true), nil
		} else {
			return Bool(false), nil
		}
	default:
		return Bool(false), nil
	}
}

func (b Bool) Cmp(v Value) (Int, Error) { return Int(0), ExpectedCmpError() }

func (b Bool) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat([]Value{b, t})

	default:
		return nil, ExpectedNumberError()
	}
}

func (b Bool) Sub(v Value) (Number, Error) { return nil, ExpectedNumberError() }
func (b Bool) Mul(v Value) (Number, Error) { return nil, ExpectedNumberError() }
func (b Bool) Div(v Value) (Number, Error) { return nil, ExpectedNumberError() }

func (b Bool) Negate() (Number, Error) { return Int(0), ExpectedNumberError() }

func (b Bool) Not() (Bool, Error) {
	return !b, nil
}

func (b Bool) Select(key string) (Value, Error) { return nil, ExpectedObjError() }
func (b Bool) Put(key string, val Value) Error  { return ExpectedObjError() }
