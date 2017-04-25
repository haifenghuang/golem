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
	"bytes"
	g "golem/core"
	"reflect"
)

//---------------------------------------------------------------
// An ObjDef contains the information needed to instantiate an Obj
// instance.  ObjDefs are created at compile time, and
// are immutable at run time.

type ObjDef struct {
	Keys []string
}

//---------------------------------------------------------------
// obj

type obj struct {
	// TODO replace this with a more efficient data structure
	fields map[string]g.Value
	inited bool
}

func NewObj() Obj {
	return &obj{nil, false}
}

func (o *obj) Init(def *ObjDef, vals []g.Value) {
	o.fields = make(map[string]g.Value)
	for i, k := range def.Keys {
		o.fields[k] = vals[i]
	}
	o.inited = true
}

func (o *obj) compositeMarker() {}

func (o *obj) TypeOf() (g.Type, g.Error) {
	if !o.inited {
		return g.TOBJ, g.UninitializedObjError()
	}

	return g.TOBJ, nil
}

func (o *obj) ToStr() (g.Str, g.Error) {
	if !o.inited {
		return nil, g.UninitializedObjError()
	}

	if len(o.fields) == 0 {
		return g.MakeStr("obj {}"), nil
	}

	var buf bytes.Buffer
	buf.WriteString("obj {")
	idx := 0
	for k, v := range o.fields {
		if idx > 0 {
			buf.WriteString(",")
		}
		idx = idx + 1
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString(": ")

		s, err := v.ToStr()
		if err != nil {
			return nil, err
		}
		buf.WriteString(s.String())
	}
	buf.WriteString(" }")
	return g.MakeStr(buf.String()), nil
}

func (o *obj) HashCode() (g.Int, g.Error) {
	if !o.inited {
		return nil, g.UninitializedObjError()
	}

	// TODO $hash()
	return nil, g.TypeMismatchError("Expected Hashable Type")
}

func (o *obj) Eq(v g.Value) (g.Bool, g.Error) {
	if !o.inited {
		return g.FALSE, g.UninitializedObjError()
	}

	// TODO $eq()
	switch t := v.(type) {
	case *obj:
		return g.MakeBool(reflect.DeepEqual(o.fields, t.fields)), nil
	default:
		return g.FALSE, nil
	}
}

func (o *obj) Cmp(v g.Value) (g.Int, g.Error) {
	if !o.inited {
		return nil, g.UninitializedObjError()
	}

	return nil, g.TypeMismatchError("Expected Comparable Type")
}

func (o *obj) Add(v g.Value) (g.Value, g.Error) {
	if !o.inited {
		return nil, g.UninitializedObjError()
	}

	switch t := v.(type) {

	case g.Str:
		return g.Strcat(o, t)

	default:
		return nil, g.TypeMismatchError("Expected Number Type")
	}
}

func (o *obj) Get(index g.Value) (g.Value, g.Error) {
	if !o.inited {
		return nil, g.UninitializedObjError()
	}

	if s, ok := index.(g.Str); ok {
		return o.GetField(s)
	} else {
		return nil, g.TypeMismatchError("Expected 'Str'")
	}
}

func (o *obj) Set(index g.Value, val g.Value) g.Error {
	if !o.inited {
		return g.UninitializedObjError()
	}

	if s, ok := index.(g.Str); ok {
		return o.PutField(s, val)
	} else {
		return g.TypeMismatchError("Expected 'Str'")
	}
}

func (o *obj) GetField(key g.Str) (g.Value, g.Error) {
	if !o.inited {
		return nil, g.UninitializedObjError()
	}

	v, ok := o.fields[key.String()]
	if ok {
		return v, nil
	} else {
		return nil, g.NoSuchFieldError(key.String())
	}
}

func (o *obj) PutField(key g.Str, val g.Value) g.Error {
	if !o.inited {
		return g.UninitializedObjError()
	}

	_, ok := o.fields[key.String()]
	if ok {
		o.fields[key.String()] = val
		return nil
	} else {
		return g.NoSuchFieldError(key.String())
	}
}

func (o *obj) Has(key g.Value) (g.Bool, g.Error) {
	if !o.inited {
		return nil, g.UninitializedObjError()
	}

	if s, ok := key.(g.Str); ok {
		_, has := o.fields[s.String()]
		return g.MakeBool(has), nil
	} else {
		return nil, g.TypeMismatchError("Expected 'Str'")
	}
}
