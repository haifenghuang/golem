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

import ()

// Parse an index value.
// The value must be between 0 (inclusive) and max (exclusive).
func ParseIndex(val Value, max int) (Int, Error) {
	if i, ok := val.(Int); ok {
		n := int(i.IntVal())
		if (n < 0) || (n >= max) {
			return nil, IndexOutOfBoundsError()
		} else {
			return i, nil
		}
	} else {
		return nil, TypeMismatchError("Expected 'Int'")
	}
}

func Strcat(a Value, b Value) (Str, Error) {

	ra := valToRunes(a)
	rb := valToRunes(b)
	result := make(str, 0, len(ra)+len(rb))

	result = append(result, runesCopy(ra)...)
	result = append(result, runesCopy(rb)...)

	return result, nil
}

// copy to avoid memory leaks
func runesCopy(s []rune) []rune {
	c := make([]rune, len(s))
	copy(c, s)
	return c
}

func valToRunes(v Value) []rune {
	if sv, ok := v.(str); ok {
		return sv
	} else {
		return v.ToStr().Runes()
	}
}

func valuesEq(as []Value, bs []Value) Bool {

	if len(as) != len(bs) {
		return FALSE
	}

	for i, a := range as {
		if a.Eq(bs[i]) == FALSE {
			return FALSE
		}
	}

	return TRUE
}
