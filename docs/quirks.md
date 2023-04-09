# C quirks

This is an (incomplete) list of C quirks and current solutions for them in `cxgo`.

## Compilation and preprocessing

### `#define` instead of `const`

Due to C being C, `const` is not really a constant declaration. Because of this, most project define constants with
the preprocessor directives like `#define`.

Unfortunately, `cxgo` lets the preprocessor replace all those constant with underlying values.

This must be fixed at some point.

### Comments

Comments are usually discarded by the preprocessor.

Due to this, `cxgo` drops all comments for now.

This must be fixed at some point.

### `#include` concatenating files

In C each file is processed individually with all included files concatenated together.

Thanks to information exposed by `cc`, `cxgo` can track the source file of each declaration.
During translation, it tracks what file is currently active and emits only declarations from that file.

### `.c` and `.h` files

`cxgo` considers `.c` and `.h` files as a one unit and will automatically merge declarations from both.

### Different included file content with `#ifdef`

`cxgo` assumes that the user will configure all project-specific `#define` switches on his own.

However, there is a potential to manipulate a platform-related `#define` directives that are currently emulated by `cxgo`.
We could potentially translate the same file under different OS/arch combinations and compare the resulting code.
All differences can be placed into separate Go files with `+build` directives.

### Private fields with incorrect headers

Some projects define the same struct type differently in public and private header files.

For now `cxgo` ignores this and lets the user control what to do in this case.

We can be smarter here and guess a private-public struct pair based on names or function signatures. We can then properly
merge those structs and set private-public field names accordingly in Go.

### Different stdlib headers

Each platform may have different stdlib headers. They might also be located in different places on different distros.

To solve this, `cxgo` emulates a virtual include file system at `/_cxgo_overrides`, which contain customized C stdlib
headers bundled into `cxgo`.

It serves two purposes: provides a zero-config experience for common use cases and allows `cxgo` to implement C stdlib
differently and adapt it to the needs of Go.

### typedef void

Example from struct_FILE.h:

    typedef void _IO_lock_t;
 
    _IO_lock_t *_lock;

For now `cxgo` interprets it as int8_t



