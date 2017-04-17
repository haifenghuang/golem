#! /usr/bin/env bash

go test golem-lang/core

go test golem-lang/scanner
go test golem-lang/parser
go test golem-lang/analyzer
go test golem-lang/compiler
go test golem-lang/interpreter

