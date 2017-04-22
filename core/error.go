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
	"fmt"
	"strings"
)

type Error interface {
	error
	Kind() string
	Msg() string
}

type serror struct {
	kind string
	msg  string
}

func (e *serror) Error() string {
	if e.msg == "" {
		return e.kind
	} else {
		return strings.Join([]string{e.kind, ": ", e.msg}, "")
	}
}

func (e *serror) Kind() string { return e.kind }
func (e *serror) Msg() string  { return e.msg }

func NullValueError() Error {
	return &serror{"NullValue", ""}
}

func TypeMismatchError(msg string) Error {
	return &serror{"TypeMismatch", msg}
}

func ArityMismatchError(expected int, actual int) Error {
	return &serror{
		"ArityMismatch",
		fmt.Sprintf("Expected %d params, got %d", expected, actual)}
}

func UninitializedObjError() Error {
	return &serror{"UninitializedObj", "Obj is not yet initialized"}
}

func DivideByZeroError() Error {
	return &serror{"DivideByZero", ""}
}

func IndexOutOfBoundsError() Error {
	return &serror{"IndexOutOfBounds", ""}
}

func NoSuchFieldError(field string) Error {
	return &serror{
		"NoSuchField",
		fmt.Sprintf("Field '%s' not found", field)}
}

func InvalidArgumentError(msg string) Error {
	return &serror{"InvalidArgument", msg}
}
