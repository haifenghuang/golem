// Copyright 2017 The Golem Project Developers //
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

//import (
//	"bytes"
//)
//
////---------------------------------------------------------------
//// chain
//
//// TODO replace with a more efficient data structure
//type chain struct {
//	keySet  map[string]bool
//	structs []Struct
//}
//
////func NewChain(structs []Struct) Struct {
////}
//
//func newChain(structs []Struct) *chain {
//	if len(structs) < 2 {
//		panic("invalid chain")
//	}
//
//	ch := &chain{make(map[string]bool), structs}
//
//	for _, s := range structs {
//		for _, k := range s.keys() {
//			ch.keySet[k] = true
//		}
//	}
//
//	return ch
//}
//
//func (ch *chain) compositeMarker() {}
//
//func (ch *chain) TypeOf() Type { return TSTRUCT }
//
//func (ch *chain) ToStr() Str {
//
//	var buf bytes.Buffer
//	buf.WriteString("struct {")
//	idx := 0
//	for k := range ch.keySet {
//		if idx > 0 {
//			buf.WriteString(",")
//		}
//		idx++
//		buf.WriteString(" ")
//		buf.WriteString(k)
//		buf.WriteString(": ")
//
//		v, err := ch.GetField(str(k))
//		if err != nil {
//			panic("invalid chain")
//		}
//		buf.WriteString(v.ToStr().String())
//	}
//	buf.WriteString(" }")
//	return MakeStr(buf.String())
//}
//
//func (ch *chain) HashCode() (Int, Error) {
//	// TODO $hash()
//	return nil, TypeMismatchError("Expected Hashable Type")
//}
//
//func (ch *chain) GetField(key Str) (Value, Error) {
//	for _, s := range ch.structs {
//		v, err := s.GetField(key)
//		if err != nil {
//			if err.Kind() != "NoSuchField" {
//				return nil, err
//			}
//		} else {
//			return v, nil
//		}
//	}
//
//	return nil, NoSuchFieldError(key.String())
//}
//
//func (ch *chain) keys() []string {
//
//	keys := make([]string, len(ch.keySet), len(ch.keySet))
//	idx := 0
//	for k := range ch.keySet {
//		keys[idx] = k
//		idx++
//	}
//	return keys
//}
