Overview
--------

This is a POSIX-compliant\* shell written in [Go](https://golang.org/).

\*: (not really)

TODO
----

The following are checked if they are implemented at all. This doesn't mean
the checked feature is complete or bug-free.

- [x] Simple commands
- [x] Environment variables
- [x] Parameter expansion
- [x] Strings
- [ ] Redirection
- [ ] Piping
- [ ] Variable assignment
- [ ] Control flow

Notes
-----

Links:

- [POSIX shell reference](http://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html)
- [Bash reference manual](https://tiswww.case.edu/php/chet/bash/bashref.html)
- [Rob Pike's talk on the Golang template lexer](https://www.youtube.com/watch?v=HxaD_trXwRE)
- [Golang template lexer source](https://github.com/golang/go/blob/master/src/text/template/parse/lex.go)

### For loops

There are three forms of for loops.

##### 1. For loop with a sequence.

```
for x in `seq 1 10`; do
    echo $x
done
```

##### 2. For loop with `in` but with no sequence.

This is interpreted as an empty sequence, so no loop iterations are done
(e.g. if your sequence comes from an empty shell variable, the loop does nothing).

```
for x in; do
    echo $x
done
```

##### 3. For loop without an `in` clause.

In this case, the loop iterates over `$@`, the command line arguments.

```
for x; do
    echo $x
done
```

### `&&`, `||`

There is no order of operations for And/Or operators.

>The operators "&&" and "||" shall have equal precedence and shall be evaluated with left associativity. For example, both of the following commands write solely bar to standard output:

>false && echo foo || echo bar
>true || echo foo && echo bar
