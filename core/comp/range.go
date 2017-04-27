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

package comp

import (
	"fmt"
	g "golem/core"
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

func NewRange(from int64, to int64, step int64) (Range, g.Error) {

	switch {

	case step == 0:
		return nil, g.InvalidArgumentError("step cannot be 0")

	case ((step > 0) && (from > to)) || ((step < 0) && (from < to)):
		return &rng{from, to, step, 0}, nil

	default:
		count := int64(math.Floor(float64(to-from) / float64(step)))
		return &rng{from, to, step, count}, nil
	}
}

func (r *rng) compositeMarker() {}

func (r *rng) TypeOf() (g.Type, g.Error) {
	return g.TRANGE, nil
}

func (r *rng) ToStr() (g.Str, g.Error) {
	return g.MakeStr(fmt.Sprintf("range<%d, %d, %d>", r.from, r.to, r.step)), nil
}

func (r *rng) HashCode() (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Hashable Type")
}

func (r *rng) Eq(v g.Value) (g.Bool, g.Error) {
	switch t := v.(type) {
	case *rng:
		return g.MakeBool(reflect.DeepEqual(r, t)), nil
	default:
		return g.FALSE, nil
	}
}

func (r *rng) Cmp(v g.Value) (g.Int, g.Error) {
	return nil, g.TypeMismatchError("Expected Comparable Type")
}

func (r *rng) Add(v g.Value) (g.Value, g.Error) {
	switch t := v.(type) {

	case g.Str:
		return g.Strcat(r, t)

	default:
		return nil, g.TypeMismatchError("Expected Number Type")
	}
}

func (r *rng) Get(index g.Value) (g.Value, g.Error) {
	idx, err := g.ParseIndex(index, int(r.count))
	if err != nil {
		return nil, err
	}
	return g.MakeInt(r.from + idx.IntVal()*r.step), nil
}

func (r *rng) Len() (g.Int, g.Error) {
	return g.MakeInt(r.count), nil
}

func (r *rng) Slice(from g.Value, to g.Value) (g.Value, g.Error) {

	f, err := g.ParseIndex(from, int(r.count))
	if err != nil {
		return nil, err
	}

	t, err := g.ParseIndex(to, int(r.count+1))
	if err != nil {
		return nil, err
	}

	// TODO do we want a different error here?
	if t.IntVal() < f.IntVal() {
		return nil, g.IndexOutOfBoundsError()
	}

	return NewRange(
		r.from+f.IntVal()*r.step,
		r.from+t.IntVal()*r.step,
		r.step)
}

func (r *rng) SliceFrom(from g.Value) (g.Value, g.Error) {
	return r.Slice(from, g.MakeInt(int64(r.count)))
}

func (r *rng) SliceTo(to g.Value) (g.Value, g.Error) {
	return r.Slice(g.ZERO, to)
}

func (r *rng) From() g.Int { return g.MakeInt(r.from) }
func (r *rng) To() g.Int   { return g.MakeInt(r.to) }
func (r *rng) Step() g.Int { return g.MakeInt(r.step) }
