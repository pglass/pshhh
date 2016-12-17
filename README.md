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

Quickstart
----------

To build psh, use the Makefile

```
$ make psh
$ ./psh -t '/bin/echo hello world'
```

### Running the tests

To run the tests, you will need to:

1. Install dependencies managed by Glide. First, install [Glide](https://github.com/Masterminds/glide)
2. Run `glide install`
3. Run `make test`

Notes
-----

Links:

- [POSIX shell reference](http://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html)
- [Bash reference manual](https://tiswww.case.edu/php/chet/bash/bashref.html)

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

### Programs stored in environment variables

> Tilde expansions, parameter expansions, command substitutions, arithmetic expansions, and quote removals that occur within a single word expand to a single field. It is only field splitting or pathname expansion that can create multiple fields from a single word. The single exception to this rule is the expansion of the special parameter '@' within double-quotes, as described in Special Parameters.

We can execute a single command stored in an environment variable

```
bash> FOO='echo wumbo'; $FOO
wumbo
```

We cannot execute multiple commands in a variable

```
bash> FOO='echo wumbo; echo mini'; $FOO
wumbo; echo mini
```

The executed command does not expand variables

```
bash> export Y=thisisy
bash> FOO='echo $Y'; $FOO
$Y
```

We can use `eval` to make this work though

```
bash> export Y=thisisy
bash> FOO='echo $Y'; eval $FOO
$Y
```

Implementation
--------------

There are four peices:

1. The `Lexer` which converts an input string to a stream of tokens
2. The `Parser` which converts a stream of tokens to a parse tree, representing
the shell program
3. The `Interpreter` which "executes" the parse tree
4. The CLI which gets input from a user, and orchestrates the lexing, parsing,
and interpretation of that input.

### The Lexer

This is modeled after the [Golang template lexer](
https://github.com/golang/go/blob/24a088d20ad52c527f61b34217da72589e366833/src/text/template/parse/lex.go#L478),
which is described in [Rob Pike's talk on the Golang template lexer](
https://www.youtube.com/watch?v=HxaD_trXwRE).

The `Lexer` type runs a loop that advances through "state functions". In
general, state functions will,

- Read the next character(s) in the stream
- Emit token(s)
- Return the next state (which is picked up by the lexer)

This makes writing the lexer fairly straightforward and readable (see
`lex/states.go`).

Tokens are emitted onto a golang channel. I actually don't like using the
channel for a couple of reasons:

- The token channel introduces concurrency where we probably don't need it.
(There may be some tiny performance benefit, but I doubt it is noticeable.)
- It makes debug logging from the lexer and other stages appear out of order.

But channels are idiomatic for golang and it wasn't problematic.

##### Modified state functions in detail

In comparison to the [Golang template lexer](
https://github.com/golang/go/blob/master/src/text/template/parse/lex.go), the
psh lexer makes a modification to the state function signature to include a next
state:

```
// the golang template lexer's version of state functions
type stateFn func(*Lexer) stateFn

// the version of state functions in psh
type stateFn func(*Lexer, stateFn) stateFn
```

Why do I want a state function that accepts the next state? Well, it is easier
to compose state functions into chains, or to "resume" in a previous state.

For example, the psh lexer must be able to tokenize a parameter expansion both
inside and outside a string.

```
${X}                    -- lexDollarExpansion
    Dollar      "$"     -- These are the tokens for a "bare" parameter expansion
    LeftBrace   "{"
    Name        "X"
    RightBrace  "}"

"abc${X}def"
    Quote       '"'     -- lexDoubleQuotedString
    Word        "abc"

                        -- lexDollarExpansion
    Dollar      "$"     -- These tokens inside a string are identical to those above
    LeftBrace   "{"
    Name        "X"
    RightBrace  "}"

    Word        "def"
    Quote       '"'
```

And there is recursion like `"abc${X:-"${Y:-"$Z"}"}def"` where you have strings
in parameter expansions in strings in parameter expansions in strings in ...,
which I needed to support.

In order to lex these, I wrote two state functions:

- `lexDollarExpansion` - tokenizes a parameter expansion
- `lexDoubleQuotedString` - tokenizes a string, and calls out to
`lexDollarExpansion` where it needs to

Simple!!

Well, not really. There were some issues trying to use unmodified state
functions for this use case:

- I cannot test state functions for equality. If I execute one state function,
it returns the next state. But which one? I have no way to know where I'm at,
which precludes making intelligent decisions. I could maybe use hacks, like
reflection, or I could wrap state functions in a struct, but this would
muddy the code.
- A state function can return any state, which could return other states and so
on. If I want state A to delegate to state B to consume a specific portion of
the text (and then have state A pick up where B left off), I first needed to
the number of states advance through after invoking B - which is _possible_
but that information (e.g. "how many states until I've tokenized from left
brace to right brace?") is implicit in the code. Relying on that sort of
implicit information makes the code brittle (easy to break when updated),
which is undesirable.
- State functions are not "customizable". For example, they do not accept a
convenient boolean parameter that says "yes, you are inside a string" in order
for one state to sometimes return different states depending on context. Now,
[the golang template lexer uses additional variables in "state functions"](
https://github.com/golang/go/blob/master/src/text/template/parse/lex.go#L478),
but this changes the type signature of the function (it is no longer of type
`stateFn`), and it doesn't scale well since I have to forward each parameter
along to possibly many other state functions. Another solution to this is to
store a flag on the `Lexer` type, but this defeats one of the nice things
about state functions: the code is the state! Putting state variables into the
`Lexer` meant less comprehensible code.

For the solution, I wanted something that was consistent and understandable. I
wanted *all* the state functions to be of type `stateFn`. I wanted to avoid
state variables, either as parameters or as flags on `Lexer`. I wanted to avoid
testing for states. That is, I should never have to ask "which state am I at?".
I should be able to reuse state functions easily.

The solution was to change the `stateFn` type to receive a `nextState`
parameter:

```
// the golang template lexer's version of state functions
type stateFn func(*Lexer) stateFn

// the version of state functions in psh
type stateFn func(*Lexer, stateFn) stateFn
```

This lets me compose states into chains that are easy to read:

```
// lex the end of a brace expansion, and then go to `nextState`
func lexBraceExpansionEnd(lx *Lexer, nextState stateFn) stateFn {
    ...
    // "I will get to nextState, but first we need to lex an operator, then
    // some text, then another operator, and then we can go to nextState"
    return composeStates(lx, lexOperator, lexText, lexOperator, nextState)
    ...
}
```

For my string tokenizing problem, I could now say:

```
// lex the contents of a string, then go to `nextState`
func lexDoubleQuotedStringContents(lx *Lexer, nextState stateFn) stateFn {
    ...
	if c == '$' {
        // lex the dollar expansion, lex more string contents, then go to nextState
        return composeStates(lx, lexDollarExpansion, lexDoubleQuotedStringContents, nextState)
    }
    ...
}
```

### The Parser

The `Parser` is in `ast/parser.go`. It accepts a `Lexer`, consumes all tokens
from the token channel, and returns an `ast.Node`.

Rather than having a single enormous Parser class responsible for everything,
parsing different bits of the syntax tree is delegated to the different
node types. If a node implements `ast.Parselet`, then it can be used to parse a
portion of the syntax tree.

The incoming tokens retain their line/position information, which we can use
for error messages.

##### Parser: String handling

Strings felt a bit weird. For example, the string `"abd${X}def"` has three
pieces:

- the plain string `"abc"`
- the parameter expansion `${X}`
- the plain string `"def"`

But, a string like this could represent a single command we need to execute, so
we want to be able to treat it as a single entity.

I ended up with an interface called `ast.StrPiece`. The `ast.Str` type stores a
list of `ast.StrPiece`. This lets me store both plain string and parameter
expansion node types in an `ast.Str` node:

```
&ast.Str{
    Pieces: []StrPiece{
        &ast.RawStr{"abc"},
        &ast.ParameterExpansion{VarName: "X"},
        &ast.RawStr{"def"},
    },
}
```

This gets passed along to the interpreter which knows how to evaluate and
concatenate the different string pieces into a single string, which we can
then insert into an environment variable or wherever.
