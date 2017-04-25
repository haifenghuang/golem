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
	"golem/core/comp"
	"golem/core/fn"
)

//---------------------------------------------------------------
// An execution environment, a.k.a 'stack frame'.

type frame struct {
	function fn.BytecodeFunc
	locals   []*fn.Ref
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
	mod    *fn.Module
	frames []*frame
}

func NewInterpreter(mod *fn.Module) *Interpreter {
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
	curFunc := fn.NewBytecodeFunc(tpl)

	// go
	return inp.invoke(curFunc, locals)
}

func (inp *Interpreter) invoke(curFunc fn.BytecodeFunc, locals []*fn.Ref) (g.Value, *ErrorStack) {

	pool := inp.mod.Pool
	defs := inp.mod.ObjDefs
	opc := curFunc.Template().OpCodes

	// stack and instruction pointer
	s := []g.Value{}
	ip := 0

	// loop over giant switch
	for {
		n := len(s) - 1

		switch opc[ip] {

		case fn.INVOKE:

			idx := index(opc, ip)
			params := s[n-idx+1:]

			switch fn := s[n-idx].(type) {
			case fn.BytecodeFunc:

				/////////////////////////////////
				// invoke a bytecode-defined func

				s = s[:n-idx]
				ip += 3

				// save the execution environment
				inp.frames = append(inp.frames, &frame{curFunc, locals, s, ip})

				// create a new execution environment
				curFunc = fn
				locals = newLocals(curFunc.Template().NumLocals, params)
				opc = curFunc.Template().OpCodes
				s = []g.Value{}
				ip = 0

			case fn.NativeFunc:

				/////////////////////////////////
				// invoke a natively-defined func

				val, err := fn.Invoke(params)
				if err != nil {
					return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
				}

				s = s[:n-idx]
				s = append(s, val)
				ip += 3

			default:
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Func'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

		case fn.RETURN:

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
			opc = curFunc.Template().OpCodes
			s = fr.stack
			ip = fr.instPtr

			// push the result
			s = append(s, result)

		case fn.NEW_FUNC:

			// push a function
			idx := index(opc, ip)
			tpl := inp.mod.Templates[idx]
			nf := fn.NewBytecodeFunc(tpl)
			s = append(s, nf)
			ip += 3

		case fn.FUNC_LOCAL:

			// get function from stack
			fn, ok := s[n].(fn.BytecodeFunc)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'BytecodeFunc'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// push a local onto the captures of the function
			idx := index(opc, ip)
			fn.PushCapture(locals[idx])
			ip += 3

		case fn.FUNC_CAPTURE:

			// get function from stack
			fn, ok := s[n].(fn.BytecodeFunc)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'BytecodeFunc'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// push a capture onto the captures of the function
			idx := index(opc, ip)
			fn.PushCapture(curFunc.GetCapture(idx))
			ip += 3

		case fn.NEW_OBJ:
			s = append(s, comp.NewObj())
			ip++

		case fn.INIT_OBJ:

			// look up ObjDef
			def := defs[index(opc, ip)]
			size := len(def.Keys)

			// get obj and values
			obj, ok := s[n-size].(comp.Obj)
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

		case fn.NEW_LIST:

			size := index(opc, ip)
			vals := make([]g.Value, size)
			copy(vals, s[n-size+1:])

			s = s[:n-size+1]
			s = append(s, comp.NewList(vals))
			ip += 3

		case fn.NEW_TUPLE:

			size := index(opc, ip)
			vals := make([]g.Value, size)
			copy(vals, s[n-size+1:])

			s = s[:n-size+1]
			s = append(s, comp.NewTuple(vals))
			ip += 3

		case fn.NEW_DICT:

			size := index(opc, ip)
			entries := make([]*g.HEntry, 0, size)

			numVals := size * 2
			for i := n - numVals + 1; i <= n; i += 2 {
				entries = append(entries, &g.HEntry{s[i], s[i+1]})
			}

			s = s[:n-numVals+1]
			s = append(s, comp.NewDict(g.NewHashMap(entries)))
			ip += 3

		case fn.GET_FIELD:

			idx := index(opc, ip)
			key, ok := pool[idx].(g.Str)
			if !ok {
				panic("Invalid GET_FIELD Key")
			}

			// get obj from stack
			obj, ok := s[n].(comp.Obj)
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

		case fn.PUT_FIELD:

			idx := index(opc, ip)
			key, ok := pool[idx].(g.Str)
			if !ok {
				panic("Invalid PUT_FIELD Key")
			}

			// get obj from stack
			obj, ok := s[n-1].(comp.Obj)
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

		case fn.INC_FIELD:

			idx := index(opc, ip)
			key, ok := pool[idx].(g.Str)
			if !ok {
				panic("Invalid GET_FIELD Key")
			}

			// get obj from stack
			obj, ok := s[n-1].(comp.Obj)
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

		case fn.GET_INDEX:

			// get Getable from stack
			gtb, ok := s[n-1].(g.Getable)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Getable'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// get index from stack
			idx := s[n]

			result, err := gtb.Get(idx)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n-1] = result
			s = s[:n]
			ip++

		case fn.SET_INDEX:

			// get Indexable from stack
			ibl, ok := s[n-2].(g.Indexable)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Getable'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// get index from stack
			idx := s[n-1]

			// get value from stack
			val := s[n]

			err := ibl.Set(idx, val)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n-2] = val
			s = s[:n-1]
			ip++

		case fn.INC_INDEX:

			// get Indexable from stack
			ibl, ok := s[n-2].(g.Indexable)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Getable'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// get index from stack
			idx := s[n-1]

			// get value from stack
			val := s[n]

			before, err := ibl.Get(idx)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			after, err := before.Add(val)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			err = ibl.Set(idx, after)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n-2] = before
			s = s[:n-1]
			ip++

		case fn.SLICE:

			// get Sliceable from stack
			slb, ok := s[n-2].(g.Sliceable)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Sliceable'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// get indices from stack
			from := s[n-1]
			to := s[n]

			result, err := slb.Slice(from, to)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n-2] = result
			s = s[:n-1]
			ip++

		case fn.SLICE_FROM:

			// get Sliceable from stack
			slb, ok := s[n-1].(g.Sliceable)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Sliceable'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// get index from stack
			from := s[n]

			result, err := slb.SliceFrom(from)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n-1] = result
			s = s[:n]
			ip++

		case fn.SLICE_TO:

			// get Sliceable from stack
			slb, ok := s[n-1].(g.Sliceable)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Sliceable'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			// get index from stack
			to := s[n]

			result, err := slb.SliceTo(to)
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n-1] = result
			s = s[:n]
			ip++

		case fn.LOAD_NULL:
			s = append(s, g.NULL)
			ip++
		case fn.LOAD_TRUE:
			s = append(s, g.TRUE)
			ip++
		case fn.LOAD_FALSE:
			s = append(s, g.FALSE)
			ip++
		case fn.LOAD_ZERO:
			s = append(s, g.ZERO)
			ip++
		case fn.LOAD_ONE:
			s = append(s, g.ONE)
			ip++
		case fn.LOAD_NEG_ONE:
			s = append(s, g.NEG_ONE)
			ip++

		case fn.LOAD_BUILTIN:
			idx := index(opc, ip)
			s = append(s, fn.Builtins[idx])
			ip += 3

		case fn.LOAD_CONST:
			idx := index(opc, ip)
			s = append(s, pool[idx])
			ip += 3

		case fn.LOAD_LOCAL:
			idx := index(opc, ip)
			s = append(s, locals[idx].Val)
			ip += 3

		case fn.LOAD_CAPTURE:
			idx := index(opc, ip)
			s = append(s, curFunc.GetCapture(idx).Val)
			ip += 3

		case fn.STORE_LOCAL:
			idx := index(opc, ip)
			locals[idx].Val = s[n]
			s = s[:n]
			ip += 3

		case fn.STORE_CAPTURE:
			idx := index(opc, ip)
			curFunc.GetCapture(idx).Val = s[n]
			s = s[:n]
			ip += 3

		case fn.JUMP:
			ip = index(opc, ip)

		case fn.JUMP_TRUE:
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

		case fn.JUMP_FALSE:
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

		case fn.EQ:
			b, err := s[n-1].Eq(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = b
			ip++

		case fn.NE:
			b, err := s[n-1].Eq(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = b.Not()
			ip++

		case fn.LT:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = g.MakeBool(val.IntVal() < 0)
			ip++

		case fn.LTE:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = g.MakeBool(val.IntVal() <= 0)
			ip++

		case fn.GT:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = g.MakeBool(val.IntVal() > 0)
			ip++

		case fn.GTE:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = g.MakeBool(val.IntVal() >= 0)
			ip++

		case fn.CMP:
			val, err := s[n-1].Cmp(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case fn.HAS:

			// get obj from stack
			obj, ok := s[n-1].(comp.Obj)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Obj'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			val, err := obj.Has(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case fn.ADD:
			val, err := s[n-1].Add(s[n])
			if err != nil {
				return nil, &ErrorStack{err, inp.stringFrames(curFunc, locals, s, ip)}
			}
			s = s[:n]
			s[n-1] = val
			ip++

		case fn.NOT:
			b, ok := s[n].(g.Bool)
			if !ok {
				return nil, &ErrorStack{
					g.TypeMismatchError("Expected 'Bool'"),
					inp.stringFrames(curFunc, locals, s, ip)}
			}

			s[n] = b.Not()
			ip++

		case fn.SUB:
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

		case fn.MUL:
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

		case fn.DIV:
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

		case fn.NEGATE:
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

		case fn.REM:
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

		case fn.BIT_AND:
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

		case fn.BIT_OR:
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

		case fn.BIT_XOR:
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

		case fn.LEFT_SHIFT:
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

		case fn.RIGHT_SHIFT:
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

		case fn.COMPLEMENT:
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

		case fn.DUP:
			s = append(s, s[n])
			ip++

		default:
			panic("Invalid opcode")
		}
	}
}

// Create a stack of string-representations from the current stack of execution frames
func (inp *Interpreter) stringFrames(
	curFunc fn.BytecodeFunc,
	locals []*fn.Ref,
	valueStack []g.Value,
	instPtr int) []string {

	n := len(inp.frames)
	stack := make([]string, n+1)

	lineNum := curFunc.Template().LineNumber(instPtr)
	stack = append(stack, fmt.Sprintf("    at line %d", lineNum))

	for i := n - 1; i >= 0; i-- {
		tp := inp.frames[i].function.Template()
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

func newLocals(numLocals int, params []g.Value) []*fn.Ref {
	p := len(params)
	locals := make([]*fn.Ref, numLocals, numLocals)
	for i := 0; i < numLocals; i++ {
		if i < p {
			locals[i] = &fn.Ref{params[i]}
		} else {
			locals[i] = &fn.Ref{g.NULL}
		}
	}
	return locals
}
