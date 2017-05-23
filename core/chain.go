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

import (
//"fmt"
)

//---------------------------------------------------------------
// chainStruct

// TODO replace with a more efficient data structure
type chainStruct struct {
	fields map[string]Struct
}

func NewChain(structs []Struct) Struct {
	if len(structs) < 2 {
		panic("invalid struct")
	}

	ch := &chainStruct{make(map[string]Struct)}

	// visit backwards so that earlier structs supersede later ones
	for i := len(structs) - 1; i >= 0; i-- {
		stc := structs[i]
		for _, k := range stc.keys() {
			ch.fields[k] = stc
		}
	}

	return &_struct{ch}
}

func (ch *chainStruct) has(key Str) Bool {
	_, has := ch.fields[key.String()]
	return MakeBool(has)
}

func (ch *chainStruct) get(key Str) (Value, Error) {
	stc, ok := ch.fields[key.String()]
	if ok {
		return stc.GetField(key)
	} else {
		return nil, NoSuchFieldError(key.String())
	}
}

func (ch *chainStruct) put(key Str, val Value) Error {
	stc, ok := ch.fields[key.String()]
	if ok {
		return stc.PutField(key, val)
	} else {
		return NoSuchFieldError(key.String())
	}
}

func (ch *chainStruct) keys() []string {

	keys := make([]string, len(ch.fields), len(ch.fields))
	idx := 0
	for k := range ch.fields {
		keys[idx] = k
		idx++
	}
	return keys
}
