#! /usr/bin/env bash

go test golem/core
go test golem/hashmap

go test golem/scanner
go test golem/parser
go test golem/analyzer
go test golem/compiler
go test golem/interpreter

