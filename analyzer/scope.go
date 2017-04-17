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

package analyzer

import (
	"bytes"
	"fmt"
	"golem/ast"
	"sort"
)

type scopeType int

const (
	funcType scopeType = iota
	blockType
	objType
)

func (s scopeType) String() string {
	if s == funcType {
		return "Func "
	} else if s == blockType {
		return "Block"
	} else {
		return "Obj"
	}
}

type scope struct {
	parent *scope
	defs   map[string]*ast.Variable

	scopeType scopeType
	funcScope *funcScope
	objScope  *objScope
}

type funcScope struct {
	numLocals      int
	captures       map[string]*ast.Variable
	parentCaptures map[string]*ast.Variable
}

type objScope struct {
	obj *ast.ObjExpr
}

func newFuncScope(parent *scope) *scope {
	return newScope(
		parent, funcType,
		&funcScope{
			0,
			make(map[string]*ast.Variable),
			make(map[string]*ast.Variable)},
		nil)
}

func newBlockScope(parent *scope) *scope {
	return newScope(parent, blockType, nil, nil)
}

func newObjScope(parent *scope, obj *ast.ObjExpr) *scope {
	return newScope(parent, objType, nil, &objScope{obj})
}

func newScope(
	parent *scope,
	scopeType scopeType,
	funcScope *funcScope,
	objScope *objScope) *scope {

	s := &scope{
		parent,
		make(map[string]*ast.Variable),
		scopeType,
		funcScope,
		objScope}

	return s
}

func (s *scope) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%v", s.scopeType))

	buf.WriteString(" defs:")
	buf.WriteString(mapString(s.defs))

	if s.scopeType == funcType {
		buf.WriteString(" captures:")
		buf.WriteString(mapString(s.funcScope.captures))
		buf.WriteString(" parentCaptures:")
		buf.WriteString(mapString(s.funcScope.parentCaptures))
		buf.WriteString(fmt.Sprintf(" numLocals:%d", s.funcScope.numLocals))
	}

	return buf.String()
}

func mapString(m map[string]*ast.Variable) string {

	keys := make([]string, len(m))
	i := 0
	for k, _ := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.WriteString("{")
	n := 0
	for _, k := range keys {
		if n > 0 {
			buf.WriteString(", ")
		}
		n++
		buf.WriteString(fmt.Sprintf("%v: %v", k, m[k]))
	}
	buf.WriteString("}")
	return buf.String()
}

// Define a Variable, either as a formal param for a Function,
// or via Let or Const.
func (s *scope) put(sym string, isConst bool) *ast.Variable {

	_, ok := s.defs[sym]
	if ok {
		panic("symbol is already defined")
	}
	v := &ast.Variable{incrementNumLocals(s), isConst, false}
	s.defs[sym] = v
	return v
}

// Create a variable for 'this', or return an existing 'this' variable.
func (s *scope) this() *ast.Variable {

	// find the nearest parent objScope
	os := s
	for os.scopeType != objType {
		os = os.parent
	}

	// define a 'this' variable on the objScope, if its not already defined
	v, ok := os.defs["this"]
	if !ok {
		idx := incrementNumLocals(os)
		v = &ast.Variable{idx, true, false}
		os.defs["this"] = v
		os.objScope.obj.LocalThisIndex = idx
	}

	// now call get(), from the original scope, to trigger captures in
	// any intervening functions.
	if _, ok := s.get("this"); !ok {
		panic("call to 'this' failed")
	}

	// done
	return v
}

// Increment the number of local variables in the nearest parent funcScope.
func incrementNumLocals(s *scope) int {

	for {
		if s.scopeType == funcType {
			idx := s.funcScope.numLocals
			s.funcScope.numLocals += 1
			if s.funcScope.numLocals >= (2 << 16) {
				panic("TODO wide index")
			}
			return idx
		} else {
			s = s.parent
		}
	}
}

// Get a Variable, by traversing up the scope stack.
func (s *scope) get(sym string) (*ast.Variable, bool) {

	// make a stack of every scope we encounter.
	stack := make([]*scope, 0, 4)
	z := s
	for {
		stack = append(stack, z)

		// Check to see if the symbol is already captured
		if z.scopeType == funcType {
			if v, ok := z.funcScope.captures[sym]; ok {
				return v, true
			}
		}

		// If the varable is defined in the current scope, then
		// capture anything that needs to be captured in the stack we have
		// created, and then return.
		if _, ok := z.defs[sym]; ok {
			return capture(sym, stack), true
		}

		// Still can't find it.  Go to the parent, or return false.
		if z.parent == nil {
			return nil, false
		} else {
			z = z.parent
		}
	}
}

// Create a succession of captures in the stack of scopes.
func capture(sym string, stack []*scope) *ast.Variable {

	n := len(stack)
	v := stack[n-1].defs[sym]

	if n == 1 {
		return v
	}

	// look for functions between the beginning and end, exclusive
	for a := n - 2; a >= 1; a-- {
		if stack[a].scopeType == funcType {
			s := stack[a]

			// capture a variable from the parent scope down into this parent scope
			idx := len(s.funcScope.parentCaptures)
			s.funcScope.parentCaptures[sym] = v

			// capture the variable in this scope
			v = &ast.Variable{idx, v.IsConst, true}
			s.funcScope.captures[sym] = v
		}
	}

	return v
}
