# Golem Tutorial

**Golem is a fledgling project right now.  Its not yet ready for production use.**  

Golem is a general purpose, interpreted language, with first-class functions and a 
dynamic type system.  Golem aims to combine the clean semantics of Python, 
the concurrency of Go, the flexibility of Javascript, and the embeddability of Lua.

Since Golem is a dynamic language, one would expect that this tutorial would begin
by asking you to start Golem's REPL.  However, Golem doesn't have a REPL yet :-).  So 
instead, you'll following along with this tutorial by typing things into a 
source code file ('tutorial.glm', for example), and then running `golem tutorial.glm`
to look at the results.

## Basic Types

So, lets get started.  Run the following program:

```golem
println('Hello, world.');
```

Golem has quite a few builtin functions.  `println` and its companion `print` print 
zero or more values to STDOUT.

```golem
println('a', 1, null, false);
```

The simplest kind of value in Golem is called a 'basic' value.  These types 
include boolean, string, int and float.  There is also the builtin value 
'null' which represents the absence of a value.

Golem has the usual set of operators that you would expect: ==, !=, <, >, <=, >=,
+, -, \*, /, and so forth.  

```golem
assert(1 + 2 == 3);
assert(42 / 7 == 8 - 2);
```
We will cover the operators in more detail later.  Note that we used another intrinsic 
function, `assert`, which will throw an exception if the value that is passed into 
it is not true.

Int values in Golem are signed 64 bit integers.  Float values are 64-bit.  Ints 
are coerced to floats during arithmetic and checks for equality:

```golem
assert(12 / 4.0 == 3.0);
assert(12 / 4.0 == 3);
```

The builtin function `str` returns  the string representation of a value:

```golem
assert(str(3) == '3');
```

Strings can be delimited either with a single quote or a double quote:

```golem
assert('abc\n' == "abc\n");
```

During addition, if one of the values is a string, and the other is not, then
`str` is automatically called on the other value, and the two strings are then 
concatenated together:

```golem
assert('a' + 1 == 'a1');
```

The two type-coercion cases just mentioned (coercing int to float during
arithmetic, and coercing to string during addition) are the *only* places that
type coercion ever happens in Golem.

Unlike many other dynamic languages, Golem has no concept of 'truthiness'.  The only 
things that are true or false are boolean values:

```golem
assert(true);
assert(!false);
```

The empty string, zero, null, etc. are *not* boolean, and will throw an
exception if you attempt to evaluate them in a place where a boolean value is
expected.

**TODO** intrinsic functions on basic types.

## Variables

Values can be assigned to variables. Variables are declared via either the `let` 
or `const` keyword.  It is an error to refer to a variable before it has been
declared.

## Operators and Expressions

## Control Structures

if, while, for, switch, try

## Error Handling

try, catch, finally, throw

## Functions and Closures

**TODO** lambda syntax

## Collections

## Structs
    
## Concurrency

## Standard Library

