# Config file reference

`cxgo` uses a YAML config file called `cxgo.yml`. For a usage example, see [this section](../examples/README.md#using-a-config-file).

## `root`

Specifies the root directory of C source that will be transpiled. All file names in [`files`](#files) are relative to this directory.

Directory can be specified as a relative path, in which case it will be resolved based on the config file path.

In case [`vcs`](#vcs) is specified, relative path will be resolved relative to a VCS root (e.g. Git repository root).

## `vcs`

Instead of specifying a local directory with C files, it's also possible to add a Git repository URL.
`cxgo` will automatically clone the project before transpiling.

See also: [`branch`](#branch), [`root`](#root).

## `branch`

If [`vcs`](#vcs) key is used, specifies the branch which will be cloned.

## `out`

Specifies the output path for Go source files.

## `package`

Specifies the Go package name to use in generated files.

## `include`

A list of include paths used for local header lookups (as in `#include "file.h"`).

Example:

```yaml
include:
  - /custom/include/path
```

## `sys_include`

A list of include paths used for system header lookups (as in `#include <file.h>`).

Example:

```yaml
sys_include:
  - /custom/include/path
```

## `define`

A list of `#define` directives added to all transpiled C files.

Example:

```yaml
define:
  - name: CXGO
  - name: CXGO_USED
    value: 1
  - name: VERSION
    value: '"dev"'
  - name: MYFUNC(x)
    value: my_func(&x)
```

## `int_size`

A size of the C `int` type in bytes. Defaults to a corresponding value for the current `GOARCH` value.

## `ptr_size`

A size of C pointer types in bytes. Defaults to a corresponding value for the current `GOARCH` value.

## `wchar_size`

A size of the `wchar_t` type in bytes. Defaults to `2`.

## `use_go_int`

Use Go `int` and `uint` in place of C `int` and `unsigned int`.

Defaults to explicit int sizes set by [`int_size`](#int_size) (`int32`, `uint32`).

## `skip`

Specifies a list of names of declarations to skip in all files. It allows removing specific functions/types/variables
from all generated Go files.

Example:

```yaml
skip:
  - some_func
  - some_type
  - some_var
```

See also: [`files.skip`](#filesskip).

## `replace`

Specifies a list of replacements applied to all files. See [`files.replace`](#filesreplace).

## `implicit_returns`

Automatically generates implicit returns, which are valid in C.

Defaults to `false`. This is done to cause a compilation error in Go to let the user decide if he wants to fix C code,
or add this workaround.

## `files`

A list of files to be processed by `cxgo`.

### `files.disabled`

Allows to disable this file, even though it's present in the config.

### `files.name`

Specifies a name of C file to be processed. 

If `.c` extension is specified, a corresponding `.h` file is also considered.

Can contain wildcards: `*.c`, `**/*.c`.

When [`files.content`](#filescontent) is specified, a file will be created in [`out`](#out) instead.
In this case it can have any extension.

### `files.go`

Specifies a corresponding Go file name for this C file. Defaults to `file.go`.

Example:

```yaml
files:
  - name: some_file.c
    go: custom.go
```

### `files.content`

If this key is specified, the files is not read from [`root`](#root) and instead is created with a given content in [`out`](#out).

Useful to automatically create Go files complementary to the generated code.

Example:

```yaml
files:
  - name: methods.go
    content: |
      package lib
      
      func (c *SomeCType) GetX() int {
        // ...
      }
```

### `files.max_decl`

Specifies a maximal number of declarations in a generated Go file. If there are more declarations, Go file will be split
to multiple files (`file_1.go`, `file_2.go`, etc).

Useful when converting giant C files.

### `files.skip`

Specifies a list of names of declarations to skip in a particular file. It allows removing specific functions/types/variables
from a specific Go file.

Example:

```yaml
files:
  - name: file.c
    skip:
      - some_func
      - some_type
      - some_var
```

See also: [`skip`](#skip).

### `files.replace`

Specifies a list of replacements applied to a particular file. Replacements are done after Go file is generated and formatted.

Example:

```yaml
files:
  - name: file.c
    replace:
      - old: 'asm("x")'
        new: 'panic("asm")'
      - regexp: 'a\s*b'
        new: 'a b'
      - old: |
          const FOO int32 = 1
        new: |
          const FOO = 1
```

See also [`replace`](#replace).

### `files.idents`

A list of configurations for translating identifiers (functions/types/variables) for a particular file.
See [`idents`](#idents).

## `idents`

A list of configurations for translating identifiers (functions/types/variables), applied to all files.

### `idents.name`

A name of the identifier in C.

### `idents.rename`

Sets a name for this identifier in Go

Example:

```yaml
idents:
  - name: some_func
    rename: SomeFunc
```

### `idents.alias`

Use the underlying type in place of this named type. Useful to remove unnecessary types generated from `typedef`.

It's not the same as Go alias types (`type A = B`). It will not generate the Go alias declaration and will just use
underlying type in all places.

Example:

```yaml
idents:
  - name: mytype_t
    alias: true
```

### `idents.type`

Allows overriding Go type for this identifier.

Valid values are:
- `bool` - uses Go `bool` instead of C `int`
- `slice` - uses Go `[]T` instead of C `*T`
- `iface` - uses Go `interface{}`

Example:

```yaml
idents:
  - name: mytype_t
    fields:
      - name: ptr
        type: slice
  - name: myfunc
    fields:
      - name: arg1
        type: slice
```

### `idents.flatten`

Flattens function control flow to workaround invalid gotos.

Example:

```yaml
idents:
  - name: myfunc
    flatten: true
```

### `idents.fields`

Allows controlling transpilation of struct fields or function arguments.

Example:

```yaml
idents:
  - name: mytype_t
    fields:
      - name: ptr
        type: slice
  - name: myfunc
    fields:
      - name: arg1
        type: slice
```

## `exec_before`

A command to execute before transpiling. Executed in [`root`](#root).

Example:

```yaml
exec_before: ['./configure']
```

## `exec_after`

A command to execute after transpiling. Executed in [`out`](#out).

Example:

```yaml
exec_after: ['go', 'generate']
```