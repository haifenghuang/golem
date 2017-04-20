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

package interpreter

import (
	"fmt"
	g "golem/core"
)

//---------------------------------------------------------------
// An execution environment, a.k.a 'stack frame'.

type frame struct {
	function *g.Func
	locals   []*g.Ref
	stack    []g.Value
	instPtr  int
}

//---------------------------------------------------------------
// A ErrorStack is returned when there is an unrecoverable error

type ErrorStack struct {
	Err   g.Error
	Stack []string
}

//---------------------------------------------------------------
// The Interpreter

type Interpreter struct {
	mod    *g.Module
	frames []*frame
}

func NewInterpreter(mod *g.Module) *Interpreter {
	tpl := mod.Templates[0]
	if tpl.Arity != 0 || tpl.NumCaptures != 0 {
		panic("TODO")
	}

	return &Interpreter{mod, []*frame{}}
}

func (inp *Interpreter) Init() (g.Value, *ErrorStack) {

	// use the zeroth template
	tpl := inp.mod.Templates[0]

	// create empty locals
	locals := newLocals(tpl.NumLocals, nil)
	inp.mod.Locals = locals

	// make func
	curFunc := g.NewFunc(tpl)

	// go
	return inp.invoke(curFunc, locals)
}

func (inp *Interpreter) invoke(curFunc *g.Func, locals []*g.Ref) (g.Value, *ErrorStack) {

	pool := inp.mod.Pool
	defs := inp.mod.ObjDefs
	opc := curFunc.Template.OpCodes

	// stack and instruction pointer
	s := []g.Value{}
	ip := 0

	// loop over giant switch
	for {
		n := len(s) - 1

		switch opc[ip] {

		case g.INVOKE:

			idx := index(opc, ip)

			// get function from stack
			fn, ok := s[n-idx].(*g.Func)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Func'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}
			if fn.Template.Arity != idx {
				return nil, &ErrorStack{
					g.ArityMismatchError(fn.Template.Arity, idx),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// get params from stack
			params := s[n-idx+1:]

			// remove function and params from stack
			s = s[:n-idx]

			// move the instruction pointer
			ip += 3

			// save the execution environment
			inp.frames = append(inp.frames, &frame{curFunc, locals, s, ip})

			// create a new execution environment
			curFunc = fn
			locals = newLocals(curFunc.Template.NumLocals, params)
			opc = curFunc.Template.OpCodes
			s = []g.Value{}
			ip = 0

		case g.RETURN:

			// get result from top of stack
			result := s[n]

			// if there is no frame to pop, then we are done
			if len(inp.frames) == 0 {
				return result, nil
			}

			// pop the frame
			lf := len(inp.frames) - 1
			fr := inp.frames[lf]
			inp.frames = inp.frames[:lf]

			// restore the execution environment
			curFunc = fr.function
			locals = fr.locals
			opc = curFunc.Template.OpCodes
			s = fr.stack
			ip = fr.instPtr

			// push the result
			s = append(s, result)

		case g.NEW_FUNC:

			// push a function
			idx := index(opc, ip)
			tpl := inp.mod.Templates[idx]
			nf := g.NewFunc(tpl)
			s = append(s, nf)
			ip += 3

		case g.FUNC_LOCAL:

			// get function from stack
			fn, ok := s[n].(*g.Func)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Func'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// push a local onto the captures of the function
			idx := index(opc, ip)
			fn.Captures = append(fn.Captures, locals[idx])
			ip += 3

		case g.FUNC_CAPTURE:

			// get function from stack
			fn, ok := s[n].(*g.Func)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Func'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// push a capture onto the captures of the function
			idx := index(opc, ip)
			fn.Captures = append(fn.Captures, curFunc.Captures[idx])
			ip += 3

		case g.NEW_OBJ:
			s = append(s, g.NewObj())
			ip++

		case g.INIT_OBJ:

			// look up ObjDef
			def := defs[index(opc, ip)]
			size := len(def.Keys)

			// get obj and values
			obj, ok := s[n-size].(g.Obj)
			if !ok {
				panic("Invalid INIT_OBJ")
			}
			vals := s[n-size+1:]

			// initialize object
			obj.Init(def, vals)

			// pop values
			s = s[:n-size+1]

			// done
			ip += 3

		case g.GET_FIELD:

			idx := index(opc, ip)
			key, ok := pool[idx].(g.Str)
			if !ok {
				panic("Invalid GET_FIELD Key")
			}

			// get obj from stack
			obj, ok := s[n].(g.Obj)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Obj'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			result, err := obj.GetField(key)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n] = result
			ip += 3

		case g.PUT_FIELD:

			idx := index(opc, ip)
			key, ok := pool[idx].(g.Str)
			if !ok {
				panic("Invalid PUT_FIELD Key")
			}

			// get obj from stack
			obj, ok := s[n-1].(g.Obj)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Obj'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// get value from stack
			value := s[n]

			err := obj.PutField(key, value)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n-1] = value
			s = s[:n]
			ip += 3

		case g.INC_FIELD:

			idx := index(opc, ip)
			key, ok := pool[idx].(g.Str)
			if !ok {
				panic("Invalid GET_FIELD Key")
			}

			// get obj from stack
			obj, ok := s[n-1].(g.Obj)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Obj'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// get value from stack
			value := s[n]

			before, err := obj.GetField(key)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			after, err := before.Add(value)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			err = obj.PutField(key, after)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n-1] = before
			s = s[:n]
			ip += 3

		case g.LOAD_NULL:
			s = append(s, g.NULL)
			ip++
		case g.LOAD_TRUE:
			s = append(s, g.TRUE)
			ip++
		case g.LOAD_FALSE:
			s = append(s, g.FALSE)
			ip++
		case g.LOAD_ZERO:
			s = append(s, g.ZERO)
			ip++
		case g.LOAD_ONE:
			s = append(s, g.ONE)
			ip++
		case g.LOAD_NEG_ONE:
			s = append(s, g.NEG_ONE)
			ip++

		case g.LOAD_CONST:
			idx := index(opc, ip)
			s = append(s, pool[idx])
			ip += 3

		case g.LOAD_LOCAL:
			idx := index(opc, ip)
			s = append(s, locals[idx].Val)
			ip += 3

		case g.LOAD_CAPTURE:
			idx := index(opc, ip)
			s = append(s, curFunc.Captures[idx].Val)
			ip += 3

		case g.STORE_LOCAL:
			idx := index(opc, ip)
			locals[idx].Val = s[n]
			s = s[:n]
			ip += 3

		case g.STORE_CAPTURE:
			idx := index(opc, ip)
			curFunc.Captures[idx].Val = s[n]
			s = s[:n]
			ip += 3

		case g.JUMP:
			ip = index(opc, ip)

		case g.JUMP_TRUE:
			b, ok := s[n].(g.Bool)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Bool'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			s = s[:n]
			if b.BoolVal() {
				ip = index(opc, ip)
			} else {
				ip += 3
			}

		case g.JUMP_FALSE:
			b, ok := s[n].(g.Bool)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Bool'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			s = s[:n]
			if b.BoolVal() {
				ip += 3
			} else {
				ip = index(opc, ip)
			}

		case g.EQ:
			b, err := s[n-1].Eq(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = b
			ip++

		case g.NE:
			b, err := s[n-1].Eq(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = b.Not()
			ip++

		case g.LT:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = g.MakeBool(val.IntVal() < 0)
			ip++

		case g.LTE:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = g.MakeBool(val.IntVal() <= 0)
			ip++

		case g.GT:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = g.MakeBool(val.IntVal() > 0)
			ip++

		case g.GTE:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = g.MakeBool(val.IntVal() >= 0)
			ip++

		case g.CMP:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.ADD:
			val, err := s[n-1].Add(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.NOT:
			b, ok := s[n].(g.Bool)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Bool'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n] = b.Not()
			ip++

		case g.SUB:
			z, ok := s[n-1].(g.Number)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected Number Type"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.Sub(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.MUL:
			z, ok := s[n-1].(g.Number)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected Number Type"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.Mul(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.DIV:
			z, ok := s[n-1].(g.Number)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected Number Type"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.Div(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.NEGATE:
			z, ok := s[n-1].(g.Number)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected Number Type"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.Negate()
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s[n] = val
			ip++

		case g.REM:
			z, ok := s[n-1].(g.Int)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Int'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.Rem(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.BIT_AND:
			z, ok := s[n-1].(g.Int)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Int'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.BitAnd(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.BIT_OR:
			z, ok := s[n-1].(g.Int)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Int'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.BitOr(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.BIT_XOR:
			z, ok := s[n-1].(g.Int)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Int'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.BitXOr(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.LEFT_SHIFT:
			z, ok := s[n-1].(g.Int)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Int'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.LeftShift(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.RIGHT_SHIFT:
			z, ok := s[n-1].(g.Int)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Int'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.RightShift(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case g.COMPLEMENT:
			z, ok := s[n].(g.Int)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Int'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := z.Complement()
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s[n] = val
			ip++

		case g.DUP:
			s = append(s, s[n])
			ip++

		default:
			panic("Invalid opcode")
		}
	}
}

// Create a stack of string-representations from the current stack of execution frames
func (inp *Interpreter) stringFrames(
	curFunc *g.Func,
	locals []*g.Ref,
	valueStack []g.Value,
	instPtr int) []string {

	n := len(inp.frames)
	stack := make([]string, n+1)

	lineNum := curFunc.Template.LineNumber(instPtr)
	stack = append(stack, fmt.Sprintf("    at line %d", lineNum))

	for i := n - 1; i >= 0; i-- {
		tp := inp.frames[i].function.Template
		lineNum := tp.LineNumber(inp.frames[i].instPtr)
		stack = append(stack, fmt.Sprintf("    at line %d", lineNum))
	}

	return stack
}

func index(opcodes []byte, ip int) int {
	high := opcodes[ip+1]
	low := opcodes[ip+2]
	return int(high)<<8 + int(low)
}

func newLocals(numLocals int, params []g.Value) []*g.Ref {
	p := len(params)
	locals := make([]*g.Ref, numLocals, numLocals)
	for i := 0; i < numLocals; i++ {
		if i < p {
			locals[i] = &g.Ref{params[i]}
		} else {
			locals[i] = &g.Ref{g.NULL}
		}
	}
	return locals
}
