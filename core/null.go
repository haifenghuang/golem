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

type Null struct{}

var NULL *Null = &Null{}

func (n *Null) TypeOf() (Type, Error) { return TNULL, nil }

func (n *Null) String() (Str, Error) { return Str("null"), nil }

func (n *Null) Eq(v Value) (Bool, Error) {
	switch v.(type) {
	case *Null:
		return true, nil
	default:
		return false, nil
	}
}

func (n *Null) Cmp(v Value) (Int, Error) { return 0, NullValueError() }

func (n *Null) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return strcat([]Value{n, t})

	default:
		return nil, NullValueError()
	}
}

func (n *Null) Sub(v Value) (Number, Error)    { return nil, NullValueError() }
func (n *Null) Mul(v Value) (Number, Error)    { return nil, NullValueError() }
func (n *Null) Div(v Value) (Number, Error)    { return nil, NullValueError() }
func (n *Null) Rem(v Value) (Int, Error)       { return Int(0), NullValueError() }
func (n *Null) BitAnd(v Value) (Int, Error)    { return Int(0), NullValueError() }
func (n *Null) BitOr(v Value) (Int, Error)     { return Int(0), NullValueError() }
func (n *Null) BitXOr(v Value) (Int, Error)    { return Int(0), NullValueError() }
func (n *Null) LeftShift(v Value) (Int, Error) { return Int(0), NullValueError() }
func (n *Null) RightShift(Value) (Int, Error)  { return Int(0), NullValueError() }

func (n *Null) Negate() (Number, Error)  { return Int(0), NullValueError() }
func (n *Null) Not() (Bool, Error)       { return false, NullValueError() }
func (n *Null) Complement() (Int, Error) { return Int(0), NullValueError() }

func (n *Null) GetField(key string) (Value, Error)   { return nil, NullValueError() }
func (n *Null) PutField(key string, val Value) Error { return NullValueError() }
