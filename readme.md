[![BuildStatus](https://travis-ci.org/e8vm/e8vm.png?branch=master)](https://travis-ci.org/h8liu/e8vm)

```
go get e8vm.io/e8vm/...
```

# E8VM

Project goal: a self-contained simulated computer system, including:

- `arch8`: A simulated dead simple instruction set that is barely
  enough for writing an OS (done).
- `asm8`: An assembler for `arch8` (done). 
  [Try it live!](http://lonnie.io/asmplay/)
- `g8`: A programming language that looks like Go but actually works
  like C (working on it).
- `os8`: A dead simple operating system that is written in `g8` (not
  started).
- `build8`: A simple building system for building and testing `g8` and 
  `asm8` code.
- Since `arch8` is dead simple, it can be easily ported to Javascript
  so that everything can run (slowly) in a browser.
- For self hosting, I could either rewrite everything in `g8`, or port
  golang to `os8`. Not sure which one is more practical.

If you would like to contribute, please contact me via email for
copyright/license related details.

## Readability

This is not a project of just about rebuilding everything from scratch.
This is an experiment of trying to write large software systems while
keeping the architecture easily comprhensible by new code readers.
This is how I plan to achieve it:

- **Use a simple language.** Written in golang.
- **Write in small files.** Each file has no more than 300 lines, and
  each line contains no more than 80 chars.
- **Keep no circular dependency.** With no circular dependency among
  files, the project can be plotted as a [DAG](http://8k.lonnie.io),
  where you can read the code. 

The DAG visualization gives the project an auto-generated "Table of
Contents", where a reader can read the entire project from left to
right in the graph. While the code might be still hard to read, I hope
that now a reader can provide detailed feedback without the need to
dive super deep first.  For example, to read and provide feedback to
the left-most block in a package, you now do not need to understand
the details of any other blocks in the package, because the left-most
block depends on nothing.

If you are just interested with the public interface, but does not
care about the internal implementation, GoWalker can provide the docs:
[here](https://gowalker.org/e8vm.io/e8vm/).

## About the Makefile

The `makefile` have several shortcuts for building tasks. To use the
`makefile`, you also need to install some tools:

```
go get e8vm.io/tools
go get github.com/golang/lint/golint
go get github.com/jstemmer/gotags
```

## TODOs

- stack trace on panic (p1)
- global variables init (p1)
- uninit var space (p1)
- debug info, func symbol table, codegen log
- code formatter (p1)
- refactor the `g8.ref` implementation (p1)
- faster memcpy when word aligned
- interface and varargs
- big numbers consts
- floating point numbers
- IR optimization

More small stuff:

- string compare
- code size/stack size checking
- unused variable check
- unreachable code check
- missing return check
- labeled break and continue

Extra:

- `e8doc`: variable width e8doc blocks
- `e8doc`: code search
- `asm8`: consts
- `link8`: symbol linking refactor
- `ir8`: ir generate refactor
- `e8doc`: online file system and editing
- `g8` `asm8`: doc genenerator