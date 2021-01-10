# Architecture

On a higher level, `cxgo` works like a regular C compiler. It operates on a "translation unit" (TU) level, meaning that
it considers a single file at a time with all included files concatenated.

As a regular C compiler, it runs a preprocessor and then parses the output to generate C code AST.
This part of the work is done by [cc](https://gitlab.com/cznic/cc), a C compiler frontend written in Go.

The C AST produced by `cc` is then converted by `cxgo` into a Go equivalent. `cxgo` uses a custom AST to be able to
represent both C concepts and Go concepts at the same time. Most of the decisions are taken when the translator
reaches a specific AST node.

Although `cc` type-checks the AST, `cxgo` does a separate type-check pass, adhering to Go rules this time. This allows
us to add missing casts, convert to/from `unsafe.Pointer`, etc. AST might be slightly changed at this stage, because
`cxgo` may need to insert helper calls to our C runtime, materialize function literals for expressions unsupported in Go, etc.

When the type check is done, a postprocessing step on the resulting AST is run. This step will make structural adjustments
to the AST, for example, it may adapt the `main` function to the Go standard, add implicit returns and fix `goto`.

After postprocessing completes, `cxgo` emits Go declarations for a specific C file. And the process repeats for the next TU.

Having this in mind, there are a few details missing in this explanation:

- How include files are found?
- How the mapping to the Go stdlib is done?
- How the common C pattern are rewritten to Go?
- How Go `string` is used?
- What about slices?

The following sections will provide more details on those. For even more details, please refer to [C quirks](quirks.md).

## Include files lookup

`cxgo` looks up include files similar to a regular C compiler, except that the include path are set in the config,
instead of an environment.

To make things easier, `cxgo` automatically adds `./include` and `./includes` to a list of lookup paths. This is helpful
if you need to override a specific header present on the host system.

Another interesting thick that `cxgo` uses is a virtual include filesystem. Thanks to FS hooks in `cc`, `cxgo` emulates
a directory at `/_cxgo_overrides`, which contain customized C stdlib headers bundled into `cxgo`.

It serves two purposes: provides a zero-config experience for common use cases and allows `cxgo` to implement C stdlib
differently and adapt it to the needs of Go.

## Stdlib mapping

`cxgo` implements a C stdlib mapping mechanism based on the VFS discussed above. Each virtual header file may also declare
a set of identifier-level overrides for types, functions and variables.

For example, we can take C `exit` from `stdlib.h` and define an override for this identifier to have Go name `os.Exit`.

This sounds trivial, but is a very powerful mechanism when combined with VFS.

For example, C has no notion of struct methods, right? But having control on the stdlib header content we can easily
emulate methods! Consider this example for `FILE` and `close` method:

```c
typedef struct {
	// ...
	int (*Close)(void);
} FILE;

#define close(f) ((FILE*)f)->Close()
```

Every time C preprocessor sees `close` it is replaced with a call to a function pointer field on the argument.
At the same time, the override is set on `FILE` to resolve `Close` to a method, but omit it from the struct definition.
So at the end the type checker sees a perfectly valid indirection to a `FILE` struct field, leading to a function call.
And in Go, the expression is converted to a method call that is implemented by our C runtime.

The interesting consequence of full control of stdlib and C's implicit type casts is that it's possible to define a
different function signature in the virtual stdlib header and let `cxgo` to adapt types in a best possible way.

## Rewriting C patterns to Go

There are different kinds of rewrite rules in `cxgo`.

In some cases, the stdlib override may define a Go function (e.g. `make`) and expose it to C with a different name (`_cxgo_go_make`).
Then, an AST translator has a hook that intercepts all function call AST nodes. If it matches a well-know pattern (`calloc(n, sizeof(T))`)
it will rewrite the AST node to a call to Go method with different arguments.

The second type of overrides is done on the statement level. Most statements are allowed to lead to more than one resulting
statement. When a pattern is recognized (e.g. `x = a ? 1 : 0`), the translator can emit multiple nodes that are semantically
the same, but are preferred in Go. Of course for each such case there is an ugly fallback.

## Support for Go string

All string literals are converted directly to Go string literals. However, C expects zero-terminated string literals.
Again, due to the fact that `cxgo` fully controls its own C stdlib headers, we can easily define a custom `_cxgo_go_string`
type and use it where `const char*` is expected. `cxgo` type checker will then use helper functions to convert to/from
zero-terminated strings automatically.

Of course, this approach is not ideal in terms of performance, but `cxgo` [goals](../CONTRIBUTING.md#project-goals-and-principles)
doesn't include performance as a guiding principle. We'd rather help the user read the code and let him rewrite the bottlenecks
in a more idiomatic Go, instead of providing fast but unreadable code.

## Support for slices

Although `cxgo` could in theory detect slice-like variables automatically, it doesn't implement this heuristic yet.

Instead, it allows user to [mark](config.md#identstype) specific struct fields, function arguments and variables as Go slices.

We admit that this is against our principle regarding "less human intervention", and fix is planned for the future.

For now, though, `cxgo` will check a list of user-defined overrides and will adjust all usages of the variable to use
slice-related features. Of course marking one variable as a slice may cause the user to be forced to mark dependant
variables as slices as well.
