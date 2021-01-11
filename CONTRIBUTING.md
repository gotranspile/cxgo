# Contributing guideline

CxGo is an open-source project, so your contributions are always welcome!

Here is how you can help the project:

- Find a C file or project that `cxgo` fails to translate. Please file an issue with your config and link to the project.
- Or find a case when `cxgo` changes the program behavior after the transpilation.
  We consider these bugs critical, and they definitely deserve an issue.
- Found a missing C signature in `cxgo` stdlib? Feel free to PR a fix! It's usually a [single line](#adding-c-symbol-definitions-to-the-library) change.
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

## Adding a new known header

`cxgo` bundles well-known headers to simplify the transpilation and provide the mapping to native Go libraries.

The support for the new header can be added incrementally:

1. Define an empty header. This helps avoid "not found" errors. Only useful to get one step further.
2. Define (a subset) of declarations to the header. This will help avoid "unexpected identifier" errors, but Go compilation
   will still fail without an implementation. Might still be useful, since functions can still be defined manually.
3. Provide a mapping between C and Go. This will help `cxgo` automatically replace functions with Go equivalents.

Let's consider each step separately.

### Adding a header stub

For example, let's define a new header called `mylib/xyz.h`. All known headers are defined in a separate files in the
[libs](./libs) package of `cxgo`. We can take one of the smaller files in that package as the reference (e.g. [assert](libs/assert.go))
and add our own header in a new Go file (`mylib_xyz.go`):

```go
package libs

const xyzH = "mylib/xyz.h"

func init() {
	RegisterLibrary(xyzH, func(c *Env) *Library {
		return &Library{}
	})
}
```

This is a minimal possible definition that just registers the known header without defining anything in it.

Library is registered with a full path, as used in `#include`. In our case the library can be used with `#include <mylib/xyz.h>`.

The second argument to `RegisterLibrary` is a constructor for a `Library` instance. The reason why it's needed because a
library might define different symbols depending on the environment (e.g. OS), architecture, or symbols defined in other
libraries. We will consider this in the following sections.

### Adding C symbol definitions to the library

C header is controlled by `Library.Header` field. We can either define it as a constant string, or build it incrementally
depending on some environment variables.

```go
RegisterLibrary(xyzH, func(env *Env) *Library {
	l := &Library{
		Header: `
#define MY_CONST 1
`,
    }
    l.Header += fmt.Sprintf("#define MY_PTR_SIZE %d\n", env.PtrSize())
	return l
})
```

Note that you don't need to add header guards - `cxgo` takes care of that.

As the first step for supporting a new library, it makes sense to add only a few declarations that we care about.
At this stage aim for simplicity: add stub types where necessary, use builtin C types instead of copying the library 1:1.

For example, if your code fails to transpile when using function `foo` from `mylib/xyz.h`, then find the declaration in
the original header (or in online docs), simplify it as much as possible and add it to the `Header`.

```go
RegisterLibrary(xyzH, func(env *Env) *Library {
	l := &Library{
		Header: `
void foo(void* ptr, int n);
`,
    }
	return l
})
```

Eventually all necessary declarations will be added. We consider adding original headers in full a bad practice since
it will be necessary to rewrite it partially anyway to get a better Go mapping. It also allows dropping legacy declarations
or compatibility `#ifdef` conditions that might be unused in `cxgo`.

The following step will be to introduce Go mapping to the new library.

### Mapping to Go

Mapping to Go allows to solve the following issues:

- Different names in C and Go
- Automatically importing Go package(s) for symbols
- Using native Go types
- Adding methods to struct types

We will explain how those issues are resolved in `cxgo` on the following examples.

#### Mapping functions and variables

Let's start by defining our function `foo` to map to the `Foo` function from a `github.com/example/mylib` package in Go.
We need two things for this: define the import map and the function identifier with corresponding type and names.

Import map is as simple as specifying a short and a long name for all Go packages imported by this library:

```go
Imports: map[string]string{
	"mylib": "github.com/example/mylib",
},
```

Defining a function is more interesting. We need to define the signature as expected by Go using builtin `cxgo` helpers:

```go
Idents: map[string]*types.Ident{
    "foo": types.NewIdentGo("foo", "mylib.Foo", env.FuncTT(nil, env.PtrT(nil), env.C().Int())),
},
```

Here we specify that the symbol `foo` as found in the include header will be mapped to an identifier that has a C name
`foo` and Go name `mylib.Foo`. We also specify an exact signature of a function: retuning no type (`nil`, mapped to `void`)
and accepting a `void*` and a C `int`. 

To build those types we use our `Env` object that is passed to the constructor. It allows getting builtin types for C, Go,
as well as using types from other mapped C libraries.

If you are mapping one of the funtions defined in Go stdlib or `cxgo` runtime, it's strongly advised to use the following
helper instead:

```go
Idents: map[string]*types.Ident{
    "foo": env.NewIdent("foo", "mylib.Foo", mylib.Foo, env.FuncTT(nil, env.PtrT(nil), env.C().Int())),
},
```

Notice that it uses a reference to the real `mylib.Foo` function in the `cxgo` code. It allows the helper to verify that
the signature matches the actual Go function that you want to map. This, however, must not be used for external libraries
because it introduces a new dependency to `cxgo`.

Two approaches described above allow the maximal level of control: you can declare C header separately and the mapped
symbol separately. There is an easier way to achieve the same and to let `cxgo` define the C header for this function:

```go
RegisterLibrary(xyzH, func(env *Env) *Library {
    l := &Library{
        Header: `
// no declarations here
`,
        Imports: map[string]string{
            "mylib": "github.com/example/mylib",
        },
    }
    l.Declare(
        types.NewIdentGo("foo", "mylib.Foo", env.FuncTT(nil, env.PtrT(nil), env.C().Int())),
    )
    return l
})
```

This way we can omit the declaration in the header and avoid repetition of the C symbol name.

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