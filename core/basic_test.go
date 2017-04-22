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
	//"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, flag bool) {
	if !flag {
		panic("adfadfa")
		t.Error("assertion failure")
	}
}

func ok(t *testing.T, val Value, err Error, expect Value) {

	assert(t, err == nil)

	if !reflect.DeepEqual(val, expect) {
		t.Error(val, " != ", expect)
	}
}

func fail(t *testing.T, val Value, err Error, expect string) {

	assert(t, val == nil)

	if err.Error() != expect {
		t.Error(err.Error(), " != ", expect)
	}
}

func okType(t *testing.T, val Value, expected Type) {
	tp, err := val.TypeOf()
	assert(t, tp == expected)
	assert(t, err == nil)
}

func TestNull(t *testing.T) {
	okType(t, NULL, TNULL)

	s, err := NULL.String()
	ok(t, s, err, MakeStr("null"))

	b, err := NULL.Eq(NULL)
	ok(t, b, err, TRUE)
	b, err = NULL.Eq(TRUE)
	ok(t, b, err, FALSE)
	i, err := NULL.Cmp(TRUE)
	fail(t, i, err, "NullValue")

	v, err := NULL.Add(MakeStr("a"))
	ok(t, v, err, MakeStr("nulla"))
}

func TestBool(t *testing.T) {

	s, err := TRUE.String()
	ok(t, s, err, MakeStr("true"))
	s, err = FALSE.String()
	ok(t, s, err, MakeStr("false"))

	okType(t, TRUE, TBOOL)
	okType(t, FALSE, TBOOL)

	assert(t, TRUE.BoolVal())
	assert(t, !FALSE.BoolVal())

	b, err := TRUE.Eq(TRUE)
	ok(t, b, err, TRUE)
	b, err = FALSE.Eq(FALSE)
	ok(t, b, err, TRUE)
	b, err = TRUE.Eq(FALSE)
	ok(t, b, err, FALSE)
	b, err = FALSE.Eq(TRUE)
	ok(t, b, err, FALSE)
	b, err = FALSE.Eq(MakeStr("a"))
	ok(t, b, err, FALSE)

	i, err := TRUE.Cmp(MakeInt(1))
	fail(t, i, err, "TypeMismatch: Expected Comparable Type")

	val := TRUE.Not()
	ok(t, val, nil, FALSE)
	val = FALSE.Not()
	ok(t, val, nil, TRUE)

	v, err := TRUE.Add(MakeInt(1))
	fail(t, v, err, "TypeMismatch: Expected Number Type")

	v, err = TRUE.Add(MakeStr("a"))
	ok(t, v, err, MakeStr("truea"))
}

func TestStr(t *testing.T) {
	a := MakeStr("a")
	b := MakeStr("b")

	var v Value

	v, err := a.String()
	ok(t, v, err, MakeStr("a"))
	v, err = b.String()
	ok(t, v, err, MakeStr("b"))

	okType(t, a, TSTR)
	v, err = a.Eq(b)
	ok(t, v, err, FALSE)
	v, err = b.Eq(a)
	ok(t, v, err, FALSE)
	v, err = a.Eq(a)
	ok(t, v, err, TRUE)
	v, err = a.Eq(MakeStr("a"))
	ok(t, v, err, TRUE)

	v, err = a.Cmp(MakeInt(1))
	fail(t, v, err, "TypeMismatch: Expected Comparable Type")
	v, err = a.Cmp(a)
	ok(t, v, err, MakeInt(0))
	v, err = a.Cmp(b)
	ok(t, v, err, MakeInt(-1))
	v, err = b.Cmp(a)
	ok(t, v, err, MakeInt(1))

	v, err = a.Add(MakeInt(1))
	ok(t, v, err, MakeStr("a1"))
	v, err = a.Add(NULL)
	ok(t, v, err, MakeStr("anull"))

	ab := MakeStr("ab")
	v, err = ab.Get(MakeInt(0))
	ok(t, v, err, a)
	v, err = ab.Get(MakeInt(1))
	ok(t, v, err, b)

	v, err = ab.Get(MakeInt(-1))
	fail(t, v, err, "IndexOutOfBounds")

	v, err = ab.Get(MakeInt(2))
	fail(t, v, err, "IndexOutOfBounds")

	v, err = MakeStr("").Len()
	ok(t, v, err, ZERO)

	v, err = MakeStr("a").Len()
	ok(t, v, err, ONE)

	v, err = MakeStr("abcde").Len()
	ok(t, v, err, MakeInt(5))
}

func TestRunes(t *testing.T) {

	runes := str{}
	const nihongo = "日本語"
	for _, r := range nihongo {
		runes = append(runes, r)
	}
	assert(t, string(runes) == nihongo)

	assert(t, runesEq(str{}, str{}))
	assert(t, runesEq(str{'a'}, str{'a'}))
	assert(t, runesEq(str{'a', 'b'}, str{'a', 'b'}))

	assert(t, !runesEq(str{}, str{'a'}))
	assert(t, !runesEq(str{'a'}, str{}))
	assert(t, !runesEq(str{'a'}, str{'b'}))
	assert(t, !runesEq(str{'c', 'b'}, str{'a', 'b'}))

	assert(t, runesCmp(str{}, str{}) == 0)
	assert(t, runesCmp(str{'a'}, str{'a'}) == 0)
	assert(t, runesCmp(str{'a', 'b'}, str{'a', 'b'}) == 0)

	assert(t, runesCmp(str{'a', 'b'}, str{'a', 'z'}) == -1)
	assert(t, runesCmp(str{'c', 'b'}, str{'a', 'b'}) == 1)
	assert(t, runesCmp(str{}, str{'a', 'b'}) == -2)
	assert(t, runesCmp(str{'a', 'b', 'c', 'd', 'e'}, str{'a', 'b'}) == 3)
}

func TestInt(t *testing.T) {
	a := MakeInt(0)
	b := MakeInt(1)

	s, err := a.String()
	ok(t, s, err, MakeStr("0"))
	s, err = b.String()
	ok(t, s, err, MakeStr("1"))

	okType(t, a, TINT)

	z, err := a.Eq(b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(a)
	ok(t, z, err, TRUE)
	z, err = a.Eq(MakeInt(0))
	ok(t, z, err, TRUE)
	z, err = a.Eq(MakeFloat(0.0))
	ok(t, z, err, TRUE)

	n, err := a.Cmp(TRUE)
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = a.Cmp(a)
	ok(t, n, err, MakeInt(0))
	n, err = a.Cmp(b)
	ok(t, n, err, MakeInt(-1))
	n, err = b.Cmp(a)
	ok(t, n, err, MakeInt(1))

	f := MakeFloat(0.0)
	g := MakeFloat(1.0)
	n, err = a.Cmp(f)
	ok(t, n, err, MakeInt(0))
	n, err = a.Cmp(g)
	ok(t, n, err, MakeInt(-1))
	n, err = g.Cmp(a)
	ok(t, n, err, MakeInt(1))

	val, err := a.Negate()
	ok(t, val, err, MakeInt(0))

	val, err = b.Negate()
	ok(t, val, err, MakeInt(-1))

	val, err = MakeInt(3).Sub(MakeInt(2))
	ok(t, val, err, MakeInt(1))
	val, err = MakeInt(3).Sub(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(1.0))
	val, err = MakeInt(3).Sub(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Sub(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Sub(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeInt(3).Mul(MakeInt(2))
	ok(t, val, err, MakeInt(6))
	val, err = MakeInt(3).Mul(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(6.0))
	val, err = MakeInt(3).Mul(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Mul(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Mul(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeInt(3).Div(MakeInt(2))
	ok(t, val, err, MakeInt(1))
	val, err = MakeInt(3).Div(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(1.5))
	val, err = MakeInt(3).Div(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Div(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeInt(3).Div(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeInt(3).Div(MakeInt(0))
	fail(t, val, err, "DivideByZero")
	val, err = MakeInt(3).Div(MakeFloat(0.0))
	fail(t, val, err, "DivideByZero")

	v1, err := MakeInt(3).Add(MakeInt(2))
	ok(t, v1, err, MakeInt(5))
	v1, err = MakeInt(3).Add(MakeFloat(2.0))
	ok(t, v1, err, MakeFloat(5.0))
	v1, err = MakeInt(3).Add(MakeStr("a"))
	ok(t, v1, err, MakeStr("3a"))
	v2, err := MakeInt(3).Add(FALSE)
	fail(t, v2, err, "TypeMismatch: Expected Number Type")
	v2, err = MakeInt(3).Add(NULL)
	fail(t, v2, err, "TypeMismatch: Expected Number Type")

	v1, err = MakeInt(7).Rem(MakeInt(3))
	ok(t, v1, err, MakeInt(1))
	v1, err = MakeInt(8).BitAnd(MakeInt(41))
	ok(t, v1, err, MakeInt(8&41))
	v1, err = MakeInt(8).BitOr(MakeInt(41))
	ok(t, v1, err, MakeInt(8|41))
	v1, err = MakeInt(8).BitXOr(MakeInt(41))
	ok(t, v1, err, MakeInt(8^41))
	v1, err = MakeInt(1).LeftShift(MakeInt(3))
	ok(t, v1, err, MakeInt(8))
	v1, err = MakeInt(8).RightShift(MakeInt(3))
	ok(t, v1, err, MakeInt(1))

	v1, err = MakeInt(8).RightShift(MakeStr("a"))
	fail(t, v1, err, "TypeMismatch: Expected 'Int'")

	v1, err = MakeInt(8).RightShift(MakeInt(-1))
	fail(t, v1, err, "InvalidArgument: Shift count cannot be less than zero")
	v1, err = MakeInt(8).LeftShift(MakeInt(-1))
	fail(t, v1, err, "InvalidArgument: Shift count cannot be less than zero")

	v1, err = MakeInt(0).Complement()
	ok(t, v1, err, MakeInt(-1))
}

func TestFloat(t *testing.T) {
	a := MakeFloat(0.1)
	b := MakeFloat(1.2)

	s, err := a.String()
	ok(t, s, err, MakeStr("0.1"))
	s, err = b.String()
	ok(t, s, err, MakeStr("1.2"))

	okType(t, a, TFLOAT)
	z, err := a.Eq(b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(a)
	ok(t, z, err, TRUE)
	z, err = a.Eq(MakeFloat(0.1))
	ok(t, z, err, TRUE)

	f := MakeFloat(0.0)
	g := MakeFloat(1.0)
	i := MakeInt(0)
	j := MakeInt(1)
	n, err := f.Cmp(MakeStr("f"))
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = f.Cmp(f)
	ok(t, n, err, MakeInt(0))
	n, err = f.Cmp(g)
	ok(t, n, err, MakeInt(-1))
	n, err = g.Cmp(f)
	ok(t, n, err, MakeInt(1))
	n, err = f.Cmp(i)
	ok(t, n, err, MakeInt(0))
	n, err = f.Cmp(j)
	ok(t, n, err, MakeInt(-1))
	n, err = j.Cmp(f)
	ok(t, n, err, MakeInt(1))

	z, err = MakeFloat(1.0).Eq(MakeInt(1))
	ok(t, z, err, TRUE)

	val, err := a.Negate()
	ok(t, val, err, MakeFloat(-0.1))

	val, err = MakeFloat(3.3).Sub(MakeInt(2))
	ok(t, val, err, MakeFloat(float64(3.3)-float64(int64(2))))
	val, err = MakeFloat(3.3).Sub(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(float64(3.3)-float64(2.0)))
	val, err = MakeFloat(3.3).Sub(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Sub(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Sub(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeFloat(3.3).Mul(MakeInt(2))
	ok(t, val, err, MakeFloat(float64(3.3)*float64(int64(2))))
	val, err = MakeFloat(3.3).Mul(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(float64(3.3)*float64(2.0)))
	val, err = MakeFloat(3.3).Mul(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Mul(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Mul(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeFloat(3.3).Div(MakeInt(2))
	ok(t, val, err, MakeFloat(float64(3.3)/float64(int64(2))))
	val, err = MakeFloat(3.3).Div(MakeFloat(2.0))
	ok(t, val, err, MakeFloat(float64(3.3)/float64(2.0)))
	val, err = MakeFloat(3.3).Div(MakeStr("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Div(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = MakeFloat(3.3).Div(NULL)
	fail(t, val, err, "TypeMismatch: Expected Number Type")

	val, err = MakeFloat(3.3).Div(MakeInt(0))
	fail(t, val, err, "DivideByZero")
	val, err = MakeFloat(3.3).Div(MakeFloat(0.0))
	fail(t, val, err, "DivideByZero")

	v1, err := MakeFloat(3.3).Add(MakeInt(2))
	ok(t, v1, err, MakeFloat(float64(3.3)+float64(int64(2))))
	v1, err = MakeFloat(3.3).Add(MakeFloat(2.0))
	ok(t, v1, err, MakeFloat(float64(3.3)+float64(2.0)))
	v1, err = MakeFloat(3.3).Add(MakeStr("a"))
	ok(t, v1, err, MakeStr("3.3a"))
	v1, err = MakeFloat(3.3).Add(FALSE)
	fail(t, v1, err, "TypeMismatch: Expected Number Type")
	v1, err = MakeFloat(3.3).Add(NULL)
	fail(t, v1, err, "TypeMismatch: Expected Number Type")
}
