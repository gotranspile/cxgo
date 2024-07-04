# C to Go translator

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/gotranspile/cxgo/master/LICENSE)
[![Gitter](https://badges.gitter.im/gotranspile/community.svg)](https://gitter.im/gotranspile/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)
[![Go Reference](https://pkg.go.dev/badge/github.com/gotranspile/cxgo.svg)](https://pkg.go.dev/github.com/gotranspile/cxgo)

CxGo is a tool for translating C source code to Go (aka transpiler, source-to-source compiler).

It uses [cc v3](https://modernc.org/cc/v3) for preprocessing and parsing C (no clang/gcc dependencies!) and
a custom type-checker and AST translation layer to make the best output possible.

The only requirement is: C code **must compile** with `cxgo`, including headers.

Having said that, `cxgo` uses a few tricks to make this process easier.

**TL;DR for the [project goals](CONTRIBUTING.md#project-goals-and-principles):**

1. Implement a practical C to Go translator ([no C++](https://github.com/gotranspile/cxgo/issues/1) for now).
2. Keep the output program code correct.
3. Make the generated code human-readable and as idiomatic as possible.
4. Make it easy to use and customize.

Check the [FAQ](FAQ.md) for more common question about the project.

## Status

The project is **experimental**! Do not rely on it in production and other sensitive environments!

Although it was successfully tested on multiple projects, it might _change the behavior_ of the code due to yet unknown bugs.

**Compiler test results:**

- TCC: 62/89 (70%)
- GCC: 783/1236 (63%)

**Transpiled projects:**

- [Potrace](./examples/potrace) (image vectorization library)
- [G722](https://github.com/gotranspile/g722) (audio codec)
- [PortableGL](https://github.com/TotallyGamerJet/pgl) (OpenGL 3.x implementation)
- [Physac](https://github.com/koteyur/physac-go) (2D physics engine)

## Installation

```bash
go install github.com/gotranspile/cxgo/cmd/cxgo@latest
```

or download the [latest release](https://github.com/gotranspile/cxgo/releases/latest) from Github.

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

- preserving comments from C code ([#2](https://github.com/gotranspile/cxgo/issues/2))
- `static` ([#4](https://github.com/gotranspile/cxgo/issues/4))
- `auto` ([#5](https://github.com/gotranspile/cxgo/issues/5))
- bitfields ([#6](https://github.com/gotranspile/cxgo/issues/6))
- `union` with C-identical data layout ([#7](https://github.com/gotranspile/cxgo/issues/7))
- `packed` structs ([#8](https://github.com/gotranspile/cxgo/issues/8))
- `asm`
- `case` in weird places ([#9](https://github.com/gotranspile/cxgo/issues/9))
- `goto` forbidden by Go (there is a [workaround](docs/config.md#identsflatten), though, see [#10](https://github.com/gotranspile/cxgo/issues/10))
- label variables ([#11](https://github.com/gotranspile/cxgo/issues/11))
- thread local storage ([#12](https://github.com/gotranspile/cxgo/issues/12))
- `setjmp` (will compile, but panics at runtime)
- some stdlib functions and types are missing ([good first issue!](CONTRIBUTING.md#adding-a-new-known-header))
- deep type inference (when converting to Go string/slices)
- considering multiple `#ifdef` paths for different OS/envs

## Community

Join our [community](COMMUNITY.md)! We'd like to hear back from you!

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md).

## License

[MIT](LICENSE)
