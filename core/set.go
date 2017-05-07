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

type set struct {
	hashMap *HashMap
	//add      *setAdd
	//addAll   *setAddAll
	//clear    *setClear
	//isEmpty  *setIsEmpty
	//contains *setContains
}

func NewSet(values []Value) Set {

	hashMap := EmptyHashMap()
	for _, v := range values {
		hashMap.Put(v, TRUE)
	}
	s := &set{hashMap}
	//s := &set{hashMap, nil, nil, nil, nil, nil}

	//s.addAll = &setAdd{&nativeFunc{}, s}
	//s.addAll = &setAddAll{&nativeFunc{}, s}
	//s.clear = &setClear{&nativeFunc{}, s}
	//s.isEmpty = &setIsEmpty{&nativeFunc{}, s}
	//s.contains = &setContains{&nativeFunc{}, s}

	return s
}

func (s *set) compositeMarker() {}

func (s *set) TypeOf() Type { return TDICT }

func (s *set) ToStr() Str {
	if s.hashMap.Len().IntVal() == 0 {
		return MakeStr("set {}")
	}

	var buf bytes.Buffer
	buf.WriteString("set {")
	idx := 0
	itr := s.hashMap.Iterator()

	for itr.Next() {
		entry := itr.Get()
		if idx > 0 {
			buf.WriteString(",")
		}
		idx++

		buf.WriteString(" ")
		s := entry.Key.ToStr()
		buf.WriteString(s.String())
	}

	buf.WriteString(" }")
	return MakeStr(buf.String())
}

func (s *set) HashCode() (Int, Error) {
	return nil, TypeMismatchError("Expected Hashable Type")
}

func (s *set) Eq(v Value) Bool {
	switch t := v.(type) {
	case *set:
		return MakeBool(reflect.DeepEqual(s.hashMap, t.hashMap))
	default:
		return FALSE
	}
}

func (s *set) Cmp(v Value) (Int, Error) {
	return nil, TypeMismatchError("Expected Comparable Type")
}

func (s *set) Plus(v Value) (Value, Error) {
	switch t := v.(type) {

	case Str:
		return Strcat(s, t)

	default:
		return nil, TypeMismatchError("Expected Number Type")
	}
}

func (s *set) Len() Int {
	return s.hashMap.Len()
}

//func (s *set) Clear() {
//	s.hashMap = EmptyHashMap()
//}
//
//func (s *set) IsEmpty() Bool {
//	return MakeBool(s.hashMap.Len().IntVal() == 0)
//}
//
//func (s *set) Contains(key Value) (Bool, Error) {
//	return s.hashMap.Contains(key)
//}
//
//func (s *set) AddAll(val Value) Error {
//	if ibl, ok := val.(Iterable); ok {
//		itr := ibl.NewIterator()
//		for itr.IterNext().BoolVal() {
//			v, err := itr.IterGet()
//			if err != nil {
//				return err
//			}
//			if tp, ok := v.(tuple); ok {
//				if len(tp) == 2 {
//					s.hashMap.Put(tp[0], tp[1])
//				} else {
//					return TupleLengthError(2, len(tp))
//				}
//			} else {
//				return TypeMismatchError("Expected Tuple")
//			}
//		}
//		return nil
//	} else {
//		return TypeMismatchError("Expected Iterable Type")
//	}
//}

////---------------------------------------------------------------
//// Iterator
//
//type setIterator struct {
//	Obj
//	s       *set
//	itr     *HIterator
//	hasNext bool
//}

func (s *set) NewIterator() Iterator {
	panic("NewIterator")

	//	next := &nativeIterNext{nativeFunc{}, nil}
	//	get := &nativeIterGet{nativeFunc{}, nil}
	//	// TODO make this immutable
	//	obj := NewObj([]*ObjEntry{
	//		&ObjEntry{"nextValue", next},
	//		&ObjEntry{"getValue", get}})
	//
	//	itr := &setIterator{obj, s, s.hashMap.Iterator(), false}
	//
	//	next.itr = itr
	//	get.itr = itr
	//	return itr
}

//func (i *setIterator) IterNext() Bool {
//	i.hasNext = i.itr.Next()
//	return MakeBool(i.hasNext)
//}
//
//func (i *setIterator) IterGet() (Value, Error) {
//
//	if i.hasNext {
//		entry := i.itr.Get()
//		return NewTuple([]Value{entry.Key, entry.Value}), nil
//	} else {
//		return nil, NoSuchElementError()
//	}
//}

//--------------------------------------------------------------
// intrinsic functions

func (s *set) GetField(key Str) (Value, Error) {
	//	switch key.String() {
	//	case "addAll":
	//		return s.addAll, nil
	//	case "clear":
	//		return s.clear, nil
	//	case "isEmpty":
	//		return s.isEmpty, nil
	//	case "contains":
	//		return s.contains, nil
	//	default:
	return nil, NoSuchFieldError(key.String())
	//	}
}

//type setAddAll struct {
//	*nativeFunc
//	s *set
//}
//
//type setClear struct {
//	*nativeFunc
//	s *set
//}
//
//type setIsEmpty struct {
//	*nativeFunc
//	s *set
//}
//
//type setContains struct {
//	*nativeFunc
//	s *set
//}
//
//func (f *setAddAll) Invoke(values []Value) (Value, Error) {
//	if len(values) != 1 {
//		return nil, ArityMismatchError("1", len(values))
//	}
//
//	err := f.s.AddAll(values[0])
//	if err != nil {
//		return nil, err
//	} else {
//		return f.s, nil
//	}
//}
//
//func (f *setClear) Invoke(values []Value) (Value, Error) {
//	if len(values) != 0 {
//		return nil, ArityMismatchError("0", len(values))
//	}
//	f.s.Clear()
//	return f.s, nil
//}
//
//func (f *setIsEmpty) Invoke(values []Value) (Value, Error) {
//	if len(values) != 0 {
//		return nil, ArityMismatchError("0", len(values))
//	}
//	return f.s.IsEmpty(), nil
//}
//
//func (f *setContains) Invoke(values []Value) (Value, Error) {
//	if len(values) != 1 {
//		return nil, ArityMismatchError("1", len(values))
//	}
//	return f.s.Contains(values[0])
//}
