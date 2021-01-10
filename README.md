# C to Go translator

CxGo is a tool for translating C source code to Go (aka transpiler, source-to-source compiler).

It uses [cc](https://godoc.org/modernc.org/cc/v3) for preprocessing and parsing C (no clang/gcc dependencies!) and
a custom type-checker and AST translation layer to make the best output possible.

The only requirement is: C code **must compile** with `cxgo`, including headers.

Having said that, `cxgo` uses a few tricks to make this process easier.

**TL;DR for the [project goals](CONTRIBUTING.md#project-goals-and-principles):**

1. Implement a practical C to Go translator (no C++ for now).
2. Keep the output program code correct.
3. Make the generated code human-readable and as idiomatic as possible.
4. Make it easy to use and customize.

Check the [FAQ](FAQ.md) for more common question about the project.

## Status

The project is **experimental**! Do not rely on it in production and other sensitive environments!

Although it can successfully transpile multiple projects, it might _change the behavior_ of the transpiled code due to yet unknown bugs.

**Compiler test results:**
- TCC: 62/89 (70%)
- GCC: 783/1236 (63%)

## Installation

```bash
go get -u github.com/dennwc/cxgo/cmd/cxgo
```

## How to use

The fastest way to try it is:

```bash
cxgo file main.c
```

For more details, check our [examples](./examples/README.md) section.

It will guide you through basic usage patterns as well as a more advanced ones (on real-world projects).

You may also check [FAQ](FAQ.md) if you have any issues.

## Caveats

The following C features are currently accepted by `cxgo`, but may be implemented partially or not implemented at all:

- converting `#define` directives to Go constants
- preserving comments from C code
- `static`
- `auto`
- bitfields
- `union` with C-identical data layout
- `packed` structs
- `asm`
- `case` in weird places
- `goto` forbidden by Go (there is a workaround, though)
- label variables
- thread local storage
- `setjmp` (will compile, but panics at runtime)
- some stdlib functions and types are missing ([good first issue!](CONTRIBUTING.md))
- deep type inference (when converting to Go string/slices)
- considering multiple `#ifdef` paths for different OS/envs

## Community

Join our [community](COMMUNITY.md)! We'd like to hear back from you!

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md).

## License

[MIT](LICENSE.md)