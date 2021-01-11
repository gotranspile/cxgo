# Frequently Asked Questions

## How to use it?

Check our [examples](./examples/README.md) section. It will guide you through basic usage patterns as well as a more advanced ones (on real-world projects).

## How to fix: "include file not found: xyz.h"?

This means that `cxgo` was unable to locate the included file.

Include lookup considers paths in the following order:

- Directory of the source file: `./xyz.h`
- Local include directory: `./include/xyz.h` (or `./includes/xyz.h`)
- Any user-defined include paths from the config ([`include`](docs/config.md#include) and [`sys_include`](docs/config.md#sys_include))
- Bundled headers from `cxgo`

Having this in mind you could either:

- Add a config file [directive](docs/config.md#sys_include): `sys_include: ['/your/path/here']`
  (or [`include`](docs/config.md#include) if the file is included as `"xyz.h"` and not `<xyz.h>`)
- Find and copy an included file into `./include`
- If this is a header in question is from a C stdlib, consider [contributing it](CONTRIBUTING.md#adding-a-new-known-header) to `cxgo`

## How to add support for a new header file?

See the corresponding [contribution guide section](CONTRIBUTING.md#adding-a-new-known-header).

## How to add support for a new function in existing library?

See the corresponding [contribution guide section](CONTRIBUTING.md#adding-c-symbol-definitions-to-the-library).