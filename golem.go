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

package main

import (
	"fmt"
	"golem/analyzer"
	"golem/compiler"
	//g "golem/core"
	"golem/interpreter"
	"golem/parser"
	"golem/scanner"
	"io/ioutil"
	"os"
)

func main() {

	// read source
	filename := os.Args[1]
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	source := string(buf)

	// parse
	scanner := scanner.NewScanner(source)
	parser := parser.NewParser(scanner)
	exprMod, err := parser.ParseModule()
	if err != nil {
		panic(err.Error())
	}

	// analyze
	anl := analyzer.NewAnalyzer(exprMod)
	errors := anl.Analyze()
	if len(errors) > 0 {
		panic(fmt.Sprintf("%v", errors))
	}

	// compile
	cmp := compiler.NewCompiler(anl)
	mod := cmp.Compile()

	// interpret
	intp := interpreter.NewInterpreter(mod)
	_, errTrace := intp.Init()
	if errTrace != nil {
		fmt.Printf("%v\n", errTrace.Error)
		fmt.Printf("%v\n", errTrace.StackTrace)
	}
}
