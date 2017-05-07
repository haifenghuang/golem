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
	"reflect"
)

type dict struct {
	hashMap     *HashMap
	addAll      *dictAddAll
	clear       *dictClear
	isEmpty     *dictIsEmpty
	containsKey *dictContainsKey
}

func NewDict(entries []*HEntry) Dict {

	hashMap := NewHashMap(entries)

	d := &dict{hashMap, nil, nil, nil, nil}

	d.addAll = &dictAddAll{&nativeFunc{}, d}
	d.clear = &dictClear{&nativeFunc{}, d}
	d.isEmpty = &dictIsEmpty{&nativeFunc{}, d}
	d.containsKey = &dictContainsKey{&nativeFunc{}, d}

	return d
}

func (d *dict) compositeMarker() {}

func (d *dict) TypeOf() Type { return TDICT }

func (d *dict) ToStr() Str {
	if d.hashMap.Len().IntVal() == 0 {
		return MakeStr("dict {}")
	}

	var buf bytes.Buffer
	buf.WriteString("dict {")
	idx := 0
	itr := d.hashMap.Iterator()

	for itr.Next() {
		entry := itr.Get()
		if idx > 0 {
			buf.WriteString(",")
		}
		idx++

		buf.WriteString(" ")
		s := entry.Key.ToStr()
		buf.WriteString(s.String())

		buf.WriteString(": ")
		s = entry.Value.ToStr()
		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return MakeStr(buf.String())
}

func (d *dict) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (d *dict) Eq(v Value) Bool {
	switch t := v.(type) {
	case *dict:
		return MakeBool(reflect.DeepEqual(d.hashMap, t.hashMap))
	default:
		return FALSE
	}
}

func (d *dict) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (d *dict) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return Strcat(d, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (d *dict) Get(key Value) (Value, Error) {
	return d.hashMap.Get(key)
}

func (d *dict) Set(key Value, val Value) Error {
	return d.hashMap.Put(key, val)
}

func (d *dict) Len() Int {
	return d.hashMap.Len()
}

func (d *dict) Clear() {
	d.hashMap = EmptyHashMap()
}

func (d *dict) IsEmpty() Bool {
	return MakeBool(d.hashMap.Len().IntVal() == 0)
}

func (d *dict) ContainsKey(key Value) (Bool, Error) {
	return d.hashMap.ContainsKey(key)
}

func (d *dict) AddAll(val Value) Error {
	if ibl, ok := val.(Iterable); ok {
		itr := ibl.NewIterator()
		for itr.IterNext().BoolVal() {
			v, err := itr.IterGet()
			if err != nil {
				return err
			}
			if tp, ok := v.(tuple); ok {
				if len(tp) == 2 {
					d.hashMap.Put(tp[0], tp[1])
				} else {
					return TupleLengthError(2, len(tp))
				}
			} else {
				return TypeMismatchError("Expected Tuple")
			}
		}
		return nil
	} else {
		return TypeMismatchError("Expected Iterable Type")
	}
}

//---------------------------------------------------------------
// Iterator

type dictIterator struct {
	Obj
	d       *dict
	itr     *HIterator
	hasNext bool
}

func (d *dict) NewIterator() Iterator {

	next := &nativeIterNext{nativeFunc{}, nil}
	get := &nativeIterGet{nativeFunc{}, nil}
	// TODO make this immutable
	obj := NewObj([]*ObjEntry{
		&ObjEntry{"nextValue", next},
		&ObjEntry{"getValue", get}})

	itr := &dictIterator{obj, d, d.hashMap.Iterator(), false}

	next.itr = itr
	get.itr = itr
	return itr
}

func (i *dictIterator) IterNext() Bool {
	i.hasNext = i.itr.Next()
	return MakeBool(i.hasNext)
}

func (i *dictIterator) IterGet() (Value, Error) {

	if i.hasNext {
		entry := i.itr.Get()
		return NewTuple([]Value{entry.Key, entry.Value}), nil
	} else {
		return nil, NoSuchElementError()
	}
}

//--------------------------------------------------------------
// intrinsic functions

func (d *dict) GetField(key Str) (Value, Error) {
	switch key.String() {
	case "addAll":
		return d.addAll, nil
	case "clear":
		return d.clear, nil
	case "isEmpty":
		return d.isEmpty, nil
	case "containsKey":
		return d.containsKey, nil
	default:
		return nil, NoSuchFieldError(key.String())
	}
}

type dictAddAll struct {
	*nativeFunc
	d *dict
}

type dictClear struct {
	*nativeFunc
	d *dict
}

type dictIsEmpty struct {
	*nativeFunc
	d *dict
}

type dictContainsKey struct {
	*nativeFunc
	d *dict
}

func (f *dictAddAll) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}

	err := f.d.AddAll(values[0])
	if err != nil {
		return nil, err
	} else {
		return f.d, nil
	}
}

func (f *dictClear) Invoke(values []Value) (Value, Error) {
	if len(values) != 0 {
		return nil, ArityMismatchError("0", len(values))
	}
	f.d.Clear()
	return f.d, nil
}

func (f *dictIsEmpty) Invoke(values []Value) (Value, Error) {
	if len(values) != 0 {
		return nil, ArityMismatchError("0", len(values))
	}
	return f.d.IsEmpty(), nil
}

func (f *dictContainsKey) Invoke(values []Value) (Value, Error) {
	if len(values) != 1 {
		return nil, ArityMismatchError("1", len(values))
	}
	return f.d.ContainsKey(values[0])
}
