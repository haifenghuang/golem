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

func assert(t *testing.T, flag Bool) {
	if !flag {
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
	ok(t, s, err, Str("null"))

	b, err := NULL.Eq(NULL)
	ok(t, b, err, TRUE)
	b, err = NULL.Eq(TRUE)
	ok(t, b, err, FALSE)

	i, err := NULL.Cmp(TRUE)
	fail(t, i, err, "NullValue")

	n, err := NULL.Negate()
	fail(t, n, err, "NullValue")

	val, err := NULL.Not()
	fail(t, val, err, "NullValue")

	num, err := NULL.Sub(Int(1))
	fail(t, num, err, "NullValue")

	num, err = NULL.Mul(Int(1))
	fail(t, num, err, "NullValue")

	num, err = NULL.Div(Int(1))
	fail(t, num, err, "NullValue")

	v, err := NULL.Add(Int(1))
	fail(t, v, err, "NullValue")

	v, err = NULL.Add(Str("a"))
	ok(t, v, err, Str("nulla"))
}

func TestBool(t *testing.T) {

	s, err := TRUE.String()
	ok(t, s, err, Str("true"))
	s, err = FALSE.String()
	ok(t, s, err, Str("false"))

	okType(t, TRUE, TBOOL)
	okType(t, FALSE, TBOOL)

	assert(t, TRUE)
	assert(t, !FALSE)

	b, err := TRUE.Eq(TRUE)
	ok(t, b, err, TRUE)
	b, err = FALSE.Eq(FALSE)
	ok(t, b, err, TRUE)
	b, err = TRUE.Eq(FALSE)
	ok(t, b, err, FALSE)
	b, err = FALSE.Eq(TRUE)
	ok(t, b, err, FALSE)
	b, err = FALSE.Eq(Str("a"))
	ok(t, b, err, FALSE)

	i, err := TRUE.Cmp(Int(1))
	fail(t, i, err, "TypeMismatch: Expected Comparable Type")

	n, err := TRUE.Negate()
	fail(t, n, err, "TypeMismatch: Expected Number Type")

	val, err := TRUE.Not()
	ok(t, val, err, FALSE)
	val, err = FALSE.Not()
	ok(t, val, err, TRUE)

	n, err = TRUE.Sub(Int(1))
	fail(t, n, err, "TypeMismatch: Expected Number Type")

	n, err = TRUE.Mul(Int(1))
	fail(t, n, err, "TypeMismatch: Expected Number Type")

	n, err = TRUE.Div(Int(1))
	fail(t, n, err, "TypeMismatch: Expected Number Type")

	v, err := TRUE.Add(Int(1))
	fail(t, v, err, "TypeMismatch: Expected Number Type")

	v, err = TRUE.Add(Str("a"))
	ok(t, v, err, Str("truea"))
}

func TestStr(t *testing.T) {
	a := Str("a")
	b := Str("b")

	s, err := a.String()
	ok(t, s, err, Str("a"))
	s, err = b.String()
	ok(t, s, err, Str("b"))

	okType(t, a, TSTR)
	z, err := a.Eq(b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(a)
	ok(t, z, err, TRUE)
	z, err = a.Eq(Str("a"))
	ok(t, z, err, TRUE)

	i, err := a.Cmp(Int(1))
	fail(t, i, err, "TypeMismatch: Expected Comparable Type")
	i, err = a.Cmp(a)
	ok(t, i, err, Int(0))
	i, err = a.Cmp(b)
	ok(t, i, err, Int(-1))
	i, err = b.Cmp(a)
	ok(t, i, err, Int(1))

	n, err := a.Negate()
	fail(t, n, err, "TypeMismatch: Expected Number Type")

	val, err := a.Not()
	fail(t, val, err, "TypeMismatch: Expected 'Bool'")

	n, err = a.Sub(Int(1))
	fail(t, n, err, "TypeMismatch: Expected Number Type")

	n, err = a.Mul(Int(1))
	fail(t, n, err, "TypeMismatch: Expected Number Type")

	n, err = a.Div(Int(1))
	fail(t, n, err, "TypeMismatch: Expected Number Type")

	v, err := a.Add(Int(1))
	ok(t, v, err, Str("a1"))
	v, err = a.Add(NULL)
	ok(t, v, err, Str("anull"))
}

func TestInt(t *testing.T) {
	a := Int(0)
	b := Int(1)

	s, err := a.String()
	ok(t, s, err, Str("0"))
	s, err = b.String()
	ok(t, s, err, Str("1"))

	okType(t, a, TINT)

	z, err := a.Eq(b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(a)
	ok(t, z, err, TRUE)
	z, err = a.Eq(Int(0))
	ok(t, z, err, TRUE)
	z, err = a.Eq(Float(0.0))
	ok(t, z, err, TRUE)

	n, err := a.Cmp(TRUE)
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = a.Cmp(a)
	ok(t, n, err, Int(0))
	n, err = a.Cmp(b)
	ok(t, n, err, Int(-1))
	n, err = b.Cmp(a)
	ok(t, n, err, Int(1))

	f := Float(0.0)
	g := Float(1.0)
	n, err = a.Cmp(f)
	ok(t, n, err, Int(0))
	n, err = a.Cmp(g)
	ok(t, n, err, Int(-1))
	n, err = g.Cmp(a)
	ok(t, n, err, Int(1))

	val, err := a.Negate()
	ok(t, val, err, Int(0))

	val, err = b.Negate()
	ok(t, val, err, Int(-1))

	nv, err := a.Not()
	fail(t, nv, err, "TypeMismatch: Expected 'Bool'")

	val, err = Int(3).Sub(Int(2))
	ok(t, val, err, Int(1))
	val, err = Int(3).Sub(Float(2.0))
	ok(t, val, err, Float(1.0))
	val, err = Int(3).Sub(Str("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Int(3).Sub(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Int(3).Sub(NULL)
	fail(t, val, err, "NullValue")

	val, err = Int(3).Mul(Int(2))
	ok(t, val, err, Int(6))
	val, err = Int(3).Mul(Float(2.0))
	ok(t, val, err, Float(6.0))
	val, err = Int(3).Mul(Str("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Int(3).Mul(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Int(3).Mul(NULL)
	fail(t, val, err, "NullValue")

	val, err = Int(3).Div(Int(2))
	ok(t, val, err, Int(1))
	val, err = Int(3).Div(Float(2.0))
	ok(t, val, err, Float(1.5))
	val, err = Int(3).Div(Str("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Int(3).Div(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Int(3).Div(NULL)
	fail(t, val, err, "NullValue")

	val, err = Int(3).Div(Int(0))
	fail(t, val, err, "DivideByZero")
	val, err = Int(3).Div(Float(0.0))
	fail(t, val, err, "DivideByZero")

	v1, err := Int(3).Add(Int(2))
	ok(t, v1, err, Int(5))
	v1, err = Int(3).Add(Float(2.0))
	ok(t, v1, err, Float(5.0))
	v1, err = Int(3).Add(Str("a"))
	ok(t, v1, err, Str("3a"))
	v2, err := Int(3).Add(FALSE)
	fail(t, v2, err, "TypeMismatch: Expected Number Type")
	v2, err = Int(3).Add(NULL)
	fail(t, v2, err, "NullValue")

	v1, err = Int(7).Rem(Int(3))
	ok(t, v1, err, Int(1))
	v1, err = Int(8).BitAnd(Int(41))
	ok(t, v1, err, Int(8&41))
	v1, err = Int(8).BitOr(Int(41))
	ok(t, v1, err, Int(8|41))
	v1, err = Int(8).BitXOr(Int(41))
	ok(t, v1, err, Int(8^41))
	v1, err = Int(1).LeftShift(Int(3))
	ok(t, v1, err, Int(8))
	v1, err = Int(8).RightShift(Int(3))
	ok(t, v1, err, Int(1))

	v1, err = Int(8).RightShift(Str("a"))
	fail(t, v1, err, "TypeMismatch: Expected 'Int'")

	v1, err = Int(8).RightShift(Int(-1))
	fail(t, v1, err, "InvalidArgument: Shift count cannot be less than zero")
	v1, err = Int(8).LeftShift(Int(-1))
	fail(t, v1, err, "InvalidArgument: Shift count cannot be less than zero")

	v1, err = Int(0).Complement()
	ok(t, v1, err, Int(-1))
}

func TestFloat(t *testing.T) {
	a := Float(0.1)
	b := Float(1.2)

	s, err := a.String()
	ok(t, s, err, Str("0.1"))
	s, err = b.String()
	ok(t, s, err, Str("1.2"))

	okType(t, a, TFLOAT)
	z, err := a.Eq(b)
	ok(t, z, err, FALSE)
	z, err = b.Eq(a)
	ok(t, z, err, FALSE)
	z, err = a.Eq(a)
	ok(t, z, err, TRUE)
	z, err = a.Eq(Float(0.1))
	ok(t, z, err, TRUE)

	f := Float(0.0)
	g := Float(1.0)
	i := Int(0)
	j := Int(1)
	n, err := f.Cmp(Str("f"))
	fail(t, n, err, "TypeMismatch: Expected Comparable Type")
	n, err = f.Cmp(f)
	ok(t, n, err, Int(0))
	n, err = f.Cmp(g)
	ok(t, n, err, Int(-1))
	n, err = g.Cmp(f)
	ok(t, n, err, Int(1))
	n, err = f.Cmp(i)
	ok(t, n, err, Int(0))
	n, err = f.Cmp(j)
	ok(t, n, err, Int(-1))
	n, err = j.Cmp(f)
	ok(t, n, err, Int(1))

	z, err = Float(1.0).Eq(Int(1))
	ok(t, z, err, TRUE)

	val, err := a.Negate()
	ok(t, val, err, Float(-0.1))

	nv, err := a.Not()
	fail(t, nv, err, "TypeMismatch: Expected 'Bool'")

	val, err = Float(3.3).Sub(Int(2))
	ok(t, val, err, Float(float64(3.3)-float64(int64(2))))
	val, err = Float(3.3).Sub(Float(2.0))
	ok(t, val, err, Float(float64(3.3)-float64(2.0)))
	val, err = Float(3.3).Sub(Str("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Float(3.3).Sub(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Float(3.3).Sub(NULL)
	fail(t, val, err, "NullValue")

	val, err = Float(3.3).Mul(Int(2))
	ok(t, val, err, Float(float64(3.3)*float64(int64(2))))
	val, err = Float(3.3).Mul(Float(2.0))
	ok(t, val, err, Float(float64(3.3)*float64(2.0)))
	val, err = Float(3.3).Mul(Str("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Float(3.3).Mul(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Float(3.3).Mul(NULL)
	fail(t, val, err, "NullValue")

	val, err = Float(3.3).Div(Int(2))
	ok(t, val, err, Float(float64(3.3)/float64(int64(2))))
	val, err = Float(3.3).Div(Float(2.0))
	ok(t, val, err, Float(float64(3.3)/float64(2.0)))
	val, err = Float(3.3).Div(Str("a"))
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Float(3.3).Div(FALSE)
	fail(t, val, err, "TypeMismatch: Expected Number Type")
	val, err = Float(3.3).Div(NULL)
	fail(t, val, err, "NullValue")

	val, err = Float(3.3).Div(Int(0))
	fail(t, val, err, "DivideByZero")
	val, err = Float(3.3).Div(Float(0.0))
	fail(t, val, err, "DivideByZero")

	v1, err := Float(3.3).Add(Int(2))
	ok(t, v1, err, Float(float64(3.3)+float64(int64(2))))
	v1, err = Float(3.3).Add(Float(2.0))
	ok(t, v1, err, Float(float64(3.3)+float64(2.0)))
	v1, err = Float(3.3).Add(Str("a"))
	ok(t, v1, err, Str("3.3a"))
	v1, err = Float(3.3).Add(FALSE)
	fail(t, v1, err, "TypeMismatch: Expected Number Type")
	v1, err = Float(3.3).Add(NULL)
	fail(t, v1, err, "NullValue")
}
