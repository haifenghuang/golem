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
	"math"
	"reflect"
)

//---------------------------------------------------------------
// rng

type rng struct {
	from  int64
	to    int64
	step  int64
	count int64
}

func NewRange(from int64, to int64, step int64) (Range, Error) {

	switch {

	case step == 0:
		return nil, InvalidArgumentError("step cannot be 0")

	case ((step > 0) && (from > to)) || ((step < 0) && (from < to)):
		return &rng{from, to, step, 0}, nil

	default:
		count := int64(math.Floor(float64(to-from) / float64(step)))
		return &rng{from, to, step, count}, nil
	}
}

func (r *rng) compositeMarker() {}

func (r *rng) TypeOf() (Type, Error) {
	return TRANGE, nil
}

func (r *rng) ToStr() (Str, Error) {
	return MakeStr(fmt.Sprintf("range<%d, %d, %d>", r.from, r.to, r.step)), nil
}

func (r *rng) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (r *rng) Eq(v Value) (Bool, Error) {
	switch t := v.(type) {
	case *rng:
		return MakeBool(reflect.DeepEqual(r, t)), nil
	default:
		return FALSE, nil
	}
}

func (r *rng) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (r *rng) Add(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return Strcat(r, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (r *rng) Get(index Value) (Value, Error) {
	idx, err := ParseIndex(index, int(r.count))
	if err != nil {
		return nil, err
	}
	return MakeInt(r.from + idx.IntVal()*r.step), nil
}

func (r *rng) Len() (Int, Error) {
	return MakeInt(r.count), nil
}

func (r *rng) Slice(from Value, to Value) (Value, Error) {

	f, err := ParseIndex(from, int(r.count))
	if err != nil {
		return nil, err
	}

	t, err := ParseIndex(to, int(r.count+1))
	if err != nil {
		return nil, err
	}

	// TODO do we want a different error here?
	if t.IntVal() < f.IntVal() {
		return nil, IndexOutOfBoundsError()
	}

	return NewRange(
		r.from+f.IntVal()*r.step,
		r.from+t.IntVal()*r.step,
		r.step)
}

func (r *rng) SliceFrom(from Value) (Value, Error) {
	return r.Slice(from, MakeInt(int64(r.count)))
}

func (r *rng) SliceTo(to Value) (Value, Error) {
	return r.Slice(ZERO, to)
}

func (r *rng) From() Int { return MakeInt(r.from) }
func (r *rng) To() Int   { return MakeInt(r.to) }
func (r *rng) Step() Int { return MakeInt(r.step) }

//---------------------------------------------------------------
// Iterator

type rangeIterator struct {
	Obj
	r *rng
	n int64
}

func (r *rng) NewIterator() Iterator {

	// TODO make this immutable
	obj := NewObj()

	iter := &rangeIterator{obj, r, -1}

	obj.Init(
		&ObjDef{[]string{
			"nextValue",
			"getValue"}},
		[]Value{
			&nativeIterNext{&nativeFunc{}, iter},
			&nativeIterGet{&nativeFunc{}, iter}})

	return iter
}

func (i *rangeIterator) IterNext() Bool {
	i.n++
	return MakeBool(i.n < i.r.count)
}

func (i *rangeIterator) IterGet() (Value, Error) {

	if (i.n >= 0) && (i.n < i.r.count) {
		return MakeInt(i.r.from + i.n*i.r.step), nil
	} else {
		return nil, NoSuchElementError()
	}
}
