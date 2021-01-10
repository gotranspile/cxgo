# Examples

- [Converting a single main file](#converting-a-single-main-file)
- [Converting a pair of library files (.c and .h)](#converting-a-pair-of-library-files)
- [Using a config for larger projects](#using-a-config-file)
- [Real-world examples](#real-world-examples)
   - [Potrace](potrace/README.md)

## Converting a single main file

The simplest possible use case for `cxgo` is to convert a single file. For example, create a file named `main.c`:

```c
#include <stdio.h>

void main() {
    printf("Hello world!\n");
}
```

Then convert it to Go with:

```bash
cxgo file main.c
```

It should generate a file `main.go` with content similar to this:

```go
package main

import "github.com/gotranspile/cxgo/runtime/stdio"

func main() {
	stdio.Printf("Hello world!\n")
}
```

This single-file mode is useful when testing `cxgo` behavior in comparison with `gcc` and other compilers.

Note that the file imports `github.com/gotranspile/cxgo/runtime` - this is the C standard library that `cxgo` uses.
Although we try to use Go stdlib as much as possible, in some cases it's still necessary to use wrappers to preserve
different C semantics.

The same is true for `stdio.Printf`, as opposed to `fmt.Printf` - C version will accept zero-terminated strings, while
native Go version doesn't.

Also, have you noted that you were not asked to provide a path for a `stdio.h` include (`-I` in `gcc`)?
This is because `cxgo` bundles most of the stdlib headers and automatically uses them when needed.

## Converting a pair of library files

Converting a library is similar to converting a main file. For example, given a header `lib.h`:

```c
#ifndef MYLIB_H
#define MYLIB_H

enum MyEnum {
    ENUM_ONE,
    ENUM_TWO
};

typedef struct {
    int* ptr;
    MyEnum val;
} my_struct;

my_struct* new_struct(int v);

void free_struct(my_struct* p);

#endif // MYLIB_H
```

And a library file `lib.c`:

```c
#include <stdlib.h>
#include "lib.h"

my_struct* new_struct(int v) {
    my_struct* p = malloc(sizeof(my_struct));
    p->ptr = 0;
    p->val = v;
    return p;
}

void free_struct(my_struct* p) {
    free(p);
}
```

It can then be converted with:

```
cxgo file --pkg mylib lib.c
```

Note that you need to specify the `lib.c` file, not `lib.h`. cxgo will use the header automatically. 

The process should generating an output file `lib.go` similar to this:

```go
package mylib

type MyEnum int

const (
	ENUM_ONE = MyEnum(iota)
	ENUM_TWO
)

type my_struct struct {
	Ptr *int
	Val MyEnum
}

func new_struct(v int) *my_struct {
	var p *my_struct = new(my_struct)
	p.Ptr = nil
	p.Val = MyEnum(v)
	return p
}
func free_struct(p *my_struct) {
	p = nil
}
```

This mode is useful when generating types, constants and functions from a self-contained C libraries.
The downside is that it gives no control regarding variable/function names, Go types used, etc.
For more control, see [this example](#using-a-config-file).

Note that the Go code has no references to `github.com/gotranspile/cxgo/runtime` this time, although the original C code
did use the `stdlib.h`. cxgo recognizes some common patterns like struct and array allocation and automatically replaces
them with Go equivalents.

Also, note that there are no calls to `free` in the Go version. Go has GC, so it's usually safe to just set the
corresponding variable to `nil`. But, as in this case, `free_struct` will only set its own argument to `nil`,
instead of setting the caller's variable. This may cause the memory to be kept longer than in the C version.

## Using a config file

`cxgo` also offers a [config file](../docs/config.md) that can help to customize the translation process.

We can use our previous example (`lib.c` and `lib.h`) and write a `cxgo.yml` config file for it:

```yaml
# Specifies the root path for C files
root: ./
# Specifies the output path for Go files
out: ./
# Package name for Go files
package: mylib

# Replace C int with Go int (the default is to use int32/int64)
use_go_int: true

# Allows to control identifier (type/variable/function) names and types
idents:
  - name: my_struct
    rename: MyStruct

# List of files to convert. Supports wildcards (*.c).
files:
  - name: lib.c
    # Allows to skip generating declarations by name
    skip:
      - free_struct
    # The same as a top level 'idents', but only affects a specific file
    idents:
      - name: new_struct
        rename: NewStruct
    # Allows to replace code in the output Go file (in case cxgo is not smart enough)
    replace:
      # simplify declaration, just for an example
      - old: 'var p *MyStruct ='
        new: 'p :='
```

You can now generate the Go files with:

```
cxgo
```

It should create a file `lib.go` as well as `go.mod`. The `lib.go` will now look like this:

```
package mylib

type MyEnum int

const (
	ENUM_ONE = MyEnum(iota)
	ENUM_TWO
)

type MyStruct struct {
	Ptr *int
	Val MyEnum
}

func NewStruct(v int) *MyStruct {
	p := new(MyStruct)
	p.Ptr = nil
	p.Val = MyEnum(v)
	return p
}
```

Note that names are different now, `free_struct` is now gone and the declaration in `NewStruct` is more idiomatic.
Let's break down the config file to explain those changes.

```yaml
# Specifies the root path for C files
root: ./
# Specifies the output path for Go files
out: ./
# Package name for Go files
package: mylib
```

The first part is pretty self-explanatory: it sets the input, output paths and the Go package name.

```yaml
# Replace C int with Go int (the default is to use int32/int64)
use_go_int: true
```

This option controls how C int types are interpreted. In this case we don't really care about int size, so we can
replace all C `int` types with Go `int` types, which is platform-dependant. Sometimes it may be important to keep
the original int size. In that case you can omit `use_go_int` and set `int_size` and `ptr_size` manually.

```yaml
# Allows to control identifier (type/variable/function) names and types
idents:
  - name: my_struct
    rename: MyStruct
```

This section controls different aspects of the type and name conversion. As in the example, you can rename identifiers
however you like. It's also possible to [alias](../docs/config.md#identsalias) the type, [promote it](../docs/config.md#identstype) to Go bool or slice.

```yaml
# List of files to convert. Supports wildcards (*.c).
files:
  - name: lib.c
```

This section control which [files](../docs/config.md#files) are converted. Since C compilers operate on a "translation unit"
(a single C file with all included files inlined), `cxgo` does the same and allows to selectively translate individual files.
It will only translate declarations found in the specified file (or corresponding header). To convert other included files,
you need to explicitly add them here. You can also use a regex to convert all C files.

```yaml
    # Allows to skip generating declarations by name
    skip:
      - free_struct
```

This section controls which declarations are [skipped](../docs/config.md#skip) when generating a specific Go file.
You may decide that you don't need certain C functions or type declarations, for example. You may also plan to use
manually written Go functions to replace generated ones, and [`skip`](../docs/config.md#skip) might help you avoid the needless declaration.

```yaml
    # The same as a top level 'idents', but only affects a specific file
    idents:
      - name: new_struct
        rename: NewStruct
```

There may be cases when two C files declare conflicting functions withs the same name. You can solve this by using
per-file [`idents`](../docs/config.md#filesidents) section and rename those problematic identifiers and avoid the name conflict.


```yaml
    # Allows to replace code in the output Go file (in case cxgo is not smart enough)
    replace:
      # simplify declaration, just for an example
      - old: 'var p *MyStruct ='
        new: 'p :='
```

`cxgo` tries to be smart, but sometimes it may fail to generate a good Go code. In this case you may decide to override
some parts of it. [`replace`](../docs/config.md#replace) offers a simple find-and-replace mechanism that works as a post-processing
step on the output file. As in the example, we replace `var p *MyStruct =` with a more idiomatic `p :=`. Note that the type of
variable is written as `NewStruct` - the replacement is run after all other conversions are applied. You may also use
`regexp` key instead of `old` to use regular expressions instead of an exact match.

## Real-world examples

Short examples might be good to understand the basics, but there might be a lot of different edge-cases in the wild.

To make your journey easier, we provide a set of examples that convert real-world C projects to Go.

- [Potrace](http://potrace.sourceforge.net/) is a vectorization tool written by Peter Selinger.
  You can convert it to a Go library using [our example](./potrace/README.md).