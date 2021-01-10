# Contributing guideline

CxGo is an open-source project, so your contributions are always welcome!

Here is how you can help the project:

- Find a C file or project that `cxgo` fails to translate. Please file an issue with your config and link to the project.
- Or find a case when `cxgo` changes the program behavior after the transpilation.
  We consider these bugs critical, and they definitely deserve an issue.
- Found a missing C signature in `cxgo` stdlib? Feel free to PR a fix! It's usually a single line change.
- `cxgo` might be missing a specific C stdlib implementation in our Go runtime. Again, feel free to PR an implementation.
  Even the most naive implementation is better than nothing!
- If you found one of the edge cases and filed an issue, you may go one step further and narrow down that edge case to a
  test for `cxgo`. Having a test makes fixing the bug a lot easier.
- We will also accept a larger feature contributions to `cxgo`. But before starting the work, please do [contact us](COMMUNITY.md)
  We might give some helpful context about `cxgo` internals and help flush the design.
- When you successfully convert an open-source C project to Go, we'd like to [hear your feedback](COMMUNITY.md)! Also, consider adding
  a link to your project to our [examples](examples/README.md) section.
  
If you have any questions, you can always [ask for help](COMMUNITY.md).

Also, make sure to check out [project goals](#project-goals-and-principles).

## Running tests

Before submitting a contribution, make sure to run `cxgo` tests:

```bash
go test ./...
```

This will only run a "fast" test suite.

When developing larger features you may also want to run TCC tests to make sure there are no regressions:

```bash
CXGO_RUN_TESTS_TCC=true go test ./tcc_test.go
```

GCC tests are also available, but are too slow to check for each iteration. You can run them with:

```bash
CXGO_RUN_TESTS_GCC=true go test ./gcc_test.go
```

## Adding a missing stdlib function

_TODO: explain how to add function signatures to existing headers_

## Adding a new stdlib header

_TODO: explain how to add new bundled headers_

## Useful resources

- Project [architecture](docs/architecture.md)

_TODO: add useful resources related to C, stdlib, etc_

## Project goals and principles

### 1. Provide a (practical) C-to-Go transpiler

This project aims to provide a generic C-to-Go conversion tool. Ideally, everything that can be compiled with a regular
C compiler should be accepted and translated by `cxgo`. It means we always aim to support real-world C code and add
any workarounds that might be required. However, no tool is perfect, so the output is provided on the best-effort basis.

Note that we [don't plan](https://github.com/gotranspile/cxgo/issues/1) to support C++** yet! We do believe that full C
support will eventually help us bootstrap C++ support (by converting GCC), but that requires a ton of work.
So do not file issues about C++ code, unless you are ready to bootstrap C++ support yourself :D

So the goal number one can be summarized as: be practical in converting C code to Go. Ideally without human intervention,
but even minimizing the intervention is a huge win for us.

### 2. Keep the program code correct

Transpilers are also called "source-to-source compilers". Similar to how bugs in compiler may lead to unexpected
program behavior, the same is true for transpiler such as `cxgo`. In other words, if translation rules are wrong,
`cxgo` will not only corrupt the program code, but the resulting program may misbehave and damage anything it touches.

Thus, we take correctness very seriously.

This is goal number two: generate a correct code. Or at least dump something useful for a human, but which will definitely
fail to compile (see goal 1).

### 3. Make the generated code human-readable and idiomatic

The result of the transpilation must be kept as human-readable as possible, and even improved to make it as close to
idiomatic Go as possible. This is with alignment with goal 1, but has less priority than correctness. If there is only
an ugly way to make things work correctly - we will use it. But if there is a chance to add a special case to `cxgo` that
will make code more readable without sacrificing correctness, we should pursue it.

This also has a less obvious consequence: `cxgo` won't pursue any of the virtualization technique such as VMs and other
approaches to abstract away the C runtime. This, of course, will render the code mostly unreadable.

### 4. Make it easy to use and/or customize

This is the hard one, but the tool will try to rewrite the output code to look more idiomatic and Go-like as possible
out of the box. This has a significant effect on `cxgo` design: we would try to automate as much as possible, bundle
everything that could be bundled, avoid C dependencies as much as possible (they are usually hard to configure), etc.

But "easy to use" also means that we should allow users to quickly change the output according to their needs.
And allow more advanced users to extend `cxgo` to run source code analysis over C code.

### 5. Make it fast (as the last goal)

Of course, we would like the `cxgo` to be fast, and what is more important that our Go runtime is fast.

But as mentioned above, we would prefer to generate a more human-readable code and let the user to fix bottlenecks by
rewriting parts of the code to pure Go, instead of providing output that is fast, but is impossible to modify.

As for `cxgo` itself, the transpilation is a batch process, instead of a realtime process, thus we might decide to
sacrifice translation performance if it will lead to better generated code.