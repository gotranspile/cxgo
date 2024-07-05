# Potrace

This example will guide you trough the process of converting [Potrace](http://potrace.sourceforge.net/) from C to Go.

## TL;DR

```bash
cd examples/potrace
cxgo
cd ../../.examples/potrace-go
go test ./...
xdg-open stanford.pdf
```

If the last command fails, open `stanford.pdf` or `stanford.svg` manually.

This sequence uses our example [config](cxgo.yml) to fetch Potrace source from the Github,
convert it with `cxgo` as a library (Potrace can also be build as a binary) and traces
an [image](http://potrace.sourceforge.net/img/stanford.pbm) from the official website.
You should see a vectorized output image as a result.

You just converted the C codebase to Go!

To learn more about this example, read a guide below.

## Prerequisites

This guide assumes that you have `cxgo` and `git` installed. You don't need anything else.

## Getting the C source

To make things easier, we will pull Potrace code from a [Github mirror](https://github.com/skyrpex/potrace).
This can be done by setting the corresponding config options:

```yaml
vcs: https://github.com/skyrpex/potrace.git
branch: '1.15'
```

We also need to specify the root for C files, which is `src`:

```yaml
root: ./src
```

## Minimal config

To get started with the translation, we need to set a few basic options.

First, set the output directory and package name for Go files:

```yaml
out: ../../.examples/potrace-go
package: gotrace
```

It's always a good idea to force a specific int/pointer type size, so you get the same result regardless of the host.
Let's set int and pointer size to 8 bytes (64 bit):

```yaml
int_size: 8
ptr_size: 8
```

Now, we need to decide what files to convert. If we look in the source folder, you will see files `potracelib.h` and `potracelib.c`.
Those sounds like a good starting point for a library. To convert both files (cxgo will pick `.h` file automatically),
add the following directives:

```yaml
files:
  - name: potracelib.c
```

We are ready to test it out. You should have the following file:

```yaml
vcs: https://github.com/skyrpex/potrace.git
branch: '1.15'
root: ./src
out: ../../.examples/potrace-go
package: gotrace

int_size: 8
ptr_size: 8

files:
  - name: potracelib.c
```

Let's give it a try!

```bash
cxgo
```

Unfortunately, things are not that simple. You would get an error similar to this:

```
potracelib.c: parsing failed: /tmp/potrace/src/potracelib.c:113:24: `VERSION`: expected ;
```

From here, an iterative process of fixing the compilation starts.
Don't be too scared though - it might not be as hard as you think.

## Make it compile

In the previous step we were left with the following error:

```
potracelib.c: parsing failed: /tmp/potrace/src/potracelib.c:113:24: `VERSION`: expected ;
```

If you look for the usages of `VERSION` is becomes clear that it's set to a specific value by the build system.
Since we are just reading the source, that variable is not defined, causing an error.

We can fix it by adding a custom `#define` directive in the config:

```yaml
define:
- name: VERSION
  value: '"dev"'
```

Running `cxgo` again will succeed this time. Nice! Now we have a first Go file named `potracelib.go`.

## Adding more files

If you check the `potracelib.go` file or try to build it, you will see a lot of Go error complaining about undefined
functions. That's not something scary, we just need to convert more files.

For example, there are a lot of functions named `progress_*`, so we will add `progress.h`.

Same for other files:
- `bm_to_pathlist` is defined in `decompose.c`
- `process_path` is defined in `trace.c`
- `path_t` type is defined in `curve.c`
- etc

If you keep adding new files until there are no undefined identifiers, you may get the following list:

```yaml
files:
  - name: potracelib.c
  - name: progress.h
  - name: trace.c
  - name: decompose.c
  - name: curve.c
  - name: bitmap.h
  - name: bbox.c
  - name: auxiliary.h
```

You may have noticed that there are other kinds of errors, so let's examine them.

### Duplicate declarations

Starting from `auxiliary.go`, you may see that `potrace_dpoint_t` and `dpoint_t` are already defined in a few other places.
This is one of the mistakes cxgo might make, and we'll need to remove the duplicates.

For `potrace_dpoint_t`, it's defined in `auxiliary.go`, `bbox.go` and `potracelib.go` and the definition always looks the same:

```go
type potrace_dpoint_t potrace_dpoint_s
```

Since this looks like a redefinition of the type, let's just ignore it and ask cxgo to use
the underlying type (`potrace_dpoint_s`) whenever it sees `potrace_dpoint_t`.

This can be done by adding a `idents` section to the top level of the config:

```yaml
idents:
  - name: potrace_dpoint_t
    alias: true
```

Adding to the top level instead of a specific file ensures that it will be used consistently in every file that has this declaration.

If you check the result by running `cxgo` again, you will notice that we now have `potrace_dpoint_s` defined in all 3 files.
That is the same issue as the last time, but on an underlying type this time:

```go
type potrace_dpoint_s struct {
	X float64
	Y float64
}
```

The declaration looks good, so let's just remove it from 2 files. Let's say that we want to keep it in `potracelib.go`,
but remove in other files. We can use `skip` to achieve this:

```yaml
idents:
  - name: potrace_dpoint_t
    alias: true

files:
  - name: potracelib.c
  - name: progress.h
  - name: trace.c
  - name: decompose.c
  - name: curve.c
  - name: bitmap.h
  - name: bbox.c
    skip:
      - potrace_dpoint_s
  - name: auxiliary.h
    skip:
      - potrace_dpoint_s
```

We have similar issues with other types (`dpoint_t`, `potrace_path_t`), so let's address them as well:

```yaml
idents:
  - name: potrace_dpoint_t
    alias: true
  - name: potrace_path_t
    alias: true
  - name: dpoint_t
    alias: true

files:
  - name: potracelib.c
  - name: progress.h
  - name: trace.c
  - name: decompose.c
  - name: curve.c
    skip:
      - potrace_path_s
  - name: bitmap.h
  - name: bbox.c
    skip:
      - potrace_dpoint_s
  - name: auxiliary.h
    skip:
      - potrace_dpoint_s
```

This should have solved the duplicate type declarations. But we still have some work to do.

### Function name collisions

Examining error further, you will notice weird issues with functions named `interval`, `iprod` and `bezier`.

Specifically, if we look at the declarations of `iprod`, we will notice that those are in fact different functions with the same name:

```go
// bbox.go
func iprod(a potrace_dpoint_s, b potrace_dpoint_s) float64 {
    return a.X*b.X + a.Y*b.Y
}

// trace.go
func iprod(p0 potrace_dpoint_s, p1 potrace_dpoint_s, p2 potrace_dpoint_s) float64 {
    var (
        x1 float64
        y1 float64
        x2 float64
        y2 float64
    )
    x1 = p1.X - p0.X
    y1 = p1.Y - p0.Y
    x2 = p2.X - p0.X
    y2 = p2.Y - p0.Y
    return x1*x2 + y1*y2
}
```

To fix this name collision, we can use `idents` configs for a specific file (and keep the other name intact):

```yaml
  - name: trace.c
    idents:
      - name: iprod
        rename: trace_iprod
```

This will rename the `iprod` in `trace.go` to `trace_iprod`, while version in `bbox.go` will still be named `iprod`.

Let's try to do the same for `interval` (in `bbox.go` and `auxiliary.go`):

```yaml
  - name: auxiliary.h
    skip:
      - potrace_dpoint_s
    idents:
      - name: interval
        rename: aux_interval
```

The declarations are now fixed, but there is another issue now: `trace.go` uses `interval` from `bbox.go`, but expects
a signature from `auxiliary.go` (now named `aux_interval`).

To fix this we can rename the usage of this function in `trace.go` as well:

```yaml
  - name: trace.c
    idents:
      - name: iprod
        rename: trace_iprod
      - name: interval
        rename: aux_interval
```

This will now work. Lastly, do the same for `bezier` (defined in `bbox.go` and `trace.go`):

```yaml
- name: trace.c
  idents:
    # ...
    - name: bezier
      rename: trace_bezier
```

We are mostly there! But a few smaller issue remain.

### Type conversion issues

Go and C type systems are slightly different, so you may encounter edge cases not yet recognized by cxgo.

In those cases, a builtin code replacement function may come in handy.

In our example, you may find the following function declaration:

```go
func bm_free(bm *potrace_bitmap_t) {
	if bm != nil && bm.Map != nil {
		bm_base(bm) = nil
	}
	bm = nil
}
```

The problematic line is: `bm_base(bm) = nil`. The original code reads as `free(bm_base(bm))`, but cxgo replaces `free`
with a `nil` pointer assignment, which is invalid when applied to a result of a function call.

Those issues will eventually be solved in cxgo itself, but for now we can add a workaround:

```yaml
  - name: bitmap.h
    replace:
      - old: 'bm_base(bm) = nil'
        new: 'bm.Map = nil'
```

This will replace all occurrences of `bm_base(bm) = nil` to `bm.Map = nil` in `bitmap.go`.

Those workarounds are not ideal, but may help to get the project converted without waiting for an upsteam fix.

### Undefined types

Sometimes you may find unused fields or undefined types after the conversion. This is usually due to arcane ways how
private fields can be implemented in C. 

In our case, you may see that `potrace_privstate_s` is not defined anywhere, and the only usage of the field with this
type is `st.Priv = nil`. We have two options here.

1. Use the code replacement discussed in the previous section.

2. Define the missing struct type in a separate Go file.

The second approach is cleaner, so we will proceed with it. We can either create the file manually, or let cxgo do it:

```yaml
files:
  # ...
  - name: hacks.go
    content: |
      package gotrace

      type potrace_privstate_s struct{}
```

## Testing attempt

We now have a first version that should compile. To be on the same page, the config should now contain the following:

```yaml
vcs: https://github.com/skyrpex/potrace.git
branch: '1.15'
root: ./src
out: ../../.examples/potrace-go
package: gotrace
int_size: 8
ptr_size: 8
define:
  - name: VERSION
    value: '"dev"'

idents:
  - name: potrace_dpoint_t
    alias: true
  - name: potrace_path_t
    alias: true
  - name: dpoint_t
    alias: true

files:
  - name: potracelib.c
  - name: progress.h
  - name: trace.c
    idents:
      - name: iprod
        rename: trace_iprod
      - name: interval
        rename: aux_interval
      - name: bezier
        rename: trace_bezier
  - name: decompose.c
  - name: curve.c
    skip:
      - potrace_path_s
  - name: bitmap.h
    replace:
      - old: 'bm_base(bm) = nil'
        new: 'bm.Map = nil'
  - name: bbox.c
    skip:
      - potrace_dpoint_s
  - name: auxiliary.h
    skip:
      - potrace_dpoint_s
    idents:
      - name: interval
        rename: aux_interval
  - name: hacks.go
    content: |
      package gotrace

      type potrace_privstate_s struct{}
```

We can add a dummy test file to quickly check if we can build the project or not:

```yaml
files:
  # ...
  - name: potrace_test.go
    content: |
      package gotrace

      import "testing"

      func TestBuild(t *testing.T) {

      }
```

Now `go test -v` in the target directory (`../../.examples/potrace-go`) should pass with no errors.
But that's not the real test, of course. Let's write something better:

```yaml
files:
  # ...
  - name: potrace_test.go
    content: |
      package gotrace

      import (
      	"math"
      	"testing"

      	"github.com/gotranspile/cxgo/runtime/libc"
      )

      func TestPotrace(t *testing.T) {
      	bm := bm_new(64, 64)
      	if bm == nil {
      		t.Fatal(libc.Error())
      	}
      	*bm.Map = math.MaxUint32

      	p := potrace_param_default()
      	st := potrace_trace(p, bm)
      	if st == nil {
      		t.Fatal(libc.Error())
      	} else if st.Status != 0 {
      		t.Fatal(st.Status)
      	}
      	t.Logf("%+v", st)
      }
```

Here we create a new Potrace bitmap with a size `64x64`, set some values into it, and then trace the bitmap with default
parameters.

You probably noticed that test checks values for `nil` and prints `libc.Error()`. Since C code uses `errno` to pass errors,
we should call `libc.Error()` to get the error value in case a function fails by returning `nil`.

Also, the way we set values into a bitmap is weird. The `bm.Map` is a `*potrace_word`, because `cxgo` currently won't
auto-detect slices. Because of this it's hard to set the whole bitmap for testing, so we only set the first word
(where each bit is a pixel value). We will address an issues with a slice a bit later in this guide.

Now, running this test (`go test -v`) should print something like this: 

```
=== RUN   TestPotrace
    potrace_test.go:24: &{Status:0 Plist:0xc0000d40f0 Priv:<nil>}
--- PASS: TestPotrace (0.00s)
PASS
ok      gotrace 0.002s
```

So it looks like we have some traced paths in `st.Plist`. It means that it actually works!

But the problem is that we can't quickly check the output, since we need to know how to traverse this path list.
So let's convert more files and address issues described above. We can also make the code a bit nicer :)

## Handling IO

Now we need to locate the files related to reading bitmaps from files and writing paths to output files.

### Reading bitmaps

Looks like functions related to reading files are located in `bitmap_io.c`, so let's add it to the config as well.
We would also need `bitops.h` since it contains dependencies for the previous file.

```yaml
files:
  # ...
  - name: bitmap_io.c
  - name: bitops.h
```

Running the test should succeed, so we will now be able to read bitmaps from files using `bm_read`. This function only
accepts a `stdio.File`, so we need some adapter for it. `cxgo` runtime provides a function for this purpose: `stdio.OpenFrom`.

Let's extend our test case to read a real bitmap:

```go
package gotrace

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
)

func TestPotrace(t *testing.T) {
	resp, err := http.Get("http://potrace.sourceforge.net/img/stanford.pbm")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal(resp.Status)
	}

	tfile, err := os.CreateTemp("", "potrace_")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tfile.Close()
		_ = os.Remove(tfile.Name())
	}()

	_, err = io.Copy(tfile, resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	tfile.Seek(0, io.SeekStart)

	var bm *potrace_bitmap_t
	e := bm_read(stdio.OpenFrom(tfile), 0.5, &bm)
	if e != 0 {
		t.Fatal(e)
	}

	p := potrace_param_default()
	st := potrace_trace(p, bm)
	if st == nil {
		t.Fatal(libc.Error())
	} else if st.Status != 0 {
		t.Fatal(st.Status)
	}
	t.Logf("%+v", st)
}
```

The test is very naive: it downloads an image from the Potrace page, saves it in a temporary file and feeds it into
`potrace_trace`. Of course, we could uses pipes to avoid on-disk files, or cache the file instead downloading it each time,
but for purposes of this guide we will use the simplest option for illustration purposes.

### Writing paths

Now we need to have a way to save paths generated by `potrace_trace`. Potrace has different backends, so let's pick a few
ones that should be the most useful ones: SVG and PDF.

#### SVG

First, let's try to add `backend_svg.c` to our config. We'll get a familiar:

```
backend_svg.c: parsing failed: /tmp/potrace/src/backend_svg.c:307:31: `POTRACE`: expected , or )
```

This is the same issues as we faced with `VERSION`, so let's add another `define`:

```yaml
define:
  - name: VERSION
    value: '"dev"'
  - name: POTRACE
    value: '"potrace"'
```

Now the transpilation succeeds, but we see that an `info` variable is referenced, but is not defined anywhere.

Looking at the files, it seems like there are two conflicting declarations: one from `main.h` and the second one from `mkbitmap.c`.
It makes sense, because those globals store Potrace configuration for `potrace` and `mkbitmap` binaries.

This is really unfortunate, since it would be nice to pass this configuration explicitly to avoid confusion for potential
users of our library.

We can solve it by using `replace`. First, add one of the files (`main.h`), suppress the declaration of a global `info` variable.
Then, use `replace` to rewrite Go code and add one more argument to SVG-related functions with the same name as a global.
This is a bit involved, but should make the API surface much better.

So, let's add the `main.h` (without `main.c`):

```yaml
files:
  # ...
  - name: main.h
    go: backend.go
```

Actually, that's all we need here: the header references to `info` as `extern`, and the variable is declared in `main.c`.
We are almost done here, but we now need to fix a few more undefined types:

```yaml
files:
  # ...
  - name: trans.c
  - name: progress_bar.c
```

Next, let's start adding replacement directives, starting from the top-level API functions:

```yaml
files:
  # ...
  - name: backend_svg.c
    replace:
      - old: 'func page_svg('
        new: 'func page_svg(info *info_s, '
      - old: 'func page_gimp('
        new: 'func page_gimp(info *info_s, '
```

And fix the function calls to them as well:

```yaml
    replace:
      # ...
      - old: 'page_svg(fout'
        new: 'page_svg(info, fout'
```

Now, let's repeat this for all other functions in this file that depend on `info`. The result should look like this:

```yaml
  - name: backend_svg.c
    replace:
      - old: 'func unit('
        new: 'func unit(info *info_s, '
      - old: 'func svg_moveto('
        new: 'func svg_moveto(info *info_s, '
      - old: 'func svg_rmoveto('
        new: 'func svg_rmoveto(info *info_s, '
      - old: 'func svg_lineto('
        new: 'func svg_lineto(info *info_s, '
      - old: 'func svg_curveto('
        new: 'func svg_curveto(info *info_s, '
      - old: 'func svg_path('
        new: 'func svg_path(info *info_s, '
      - old: 'func svg_jaggy_path('
        new: 'func svg_jaggy_path(info *info_s, '
      - old: 'func write_paths_opaque('
        new: 'func write_paths_opaque(info *info_s, '
      - old: 'func write_paths_transparent_rec('
        new: 'func write_paths_transparent_rec(info *info_s, '
      - old: 'func write_paths_transparent('
        new: 'func write_paths_transparent(info *info_s, '
      - old: 'func page_svg('
        new: 'func page_svg(info *info_s, '
      - old: 'func page_gimp('
        new: 'func page_gimp(info *info_s, '
      - old: 'unit(p)'
        new: 'unit(info, p)'
      - old: 'unit(p1)'
        new: 'unit(info, p1)'
      - old: 'unit(p2)'
        new: 'unit(info, p2)'
      - old: 'unit(p3)'
        new: 'unit(info, p3)'
      - old: 'svg_moveto(fout'
        new: 'svg_moveto(info, fout'
      - old: 'svg_rmoveto(fout'
        new: 'svg_rmoveto(info, fout'
      - old: 'svg_lineto(fout'
        new: 'svg_lineto(info, fout'
      - old: 'svg_curveto(fout'
        new: 'svg_curveto(info, fout'
      - old: 'svg_jaggy_path(fout'
        new: 'svg_jaggy_path(info, fout'
      - old: 'svg_path(fout'
        new: 'svg_path(info, fout'
      - old: 'write_paths_opaque(fout'
        new: 'write_paths_opaque(info, fout'
      - old: 'write_paths_transparent_rec(fout'
        new: 'write_paths_transparent_rec(info, fout'
      - old: 'write_paths_transparent(fout'
        new: 'write_paths_transparent(info, fout'
      - old: 'page_svg(fout'
        new: 'page_svg(info, fout'
```

Transpiling it will only show that `backend_s` is undefined. Since it's unused in the code, let's add a stub for it:

```yaml
  - name: hacks.go
    content: |
      package gotrace

      type backend_s struct{}
      type potrace_privstate_s struct{}
```

It should now compile again, and we can use new functions in our test:

```go
package gotrace

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
)

func TestPotrace(t *testing.T) {
	resp, err := http.Get("http://potrace.sourceforge.net/img/stanford.pbm")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatal(resp.Status)
	}

	tfile, err := os.CreateTemp("", "potrace_")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = tfile.Close()
		_ = os.Remove(tfile.Name())
	}()

	_, err = io.Copy(tfile, resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	tfile.Seek(0, io.SeekStart)

	var bm *potrace_bitmap_t
	e := bm_read(stdio.OpenFrom(tfile), 0.5, &bm)
	if e != 0 {
		t.Fatal(e)
	}

	p := potrace_param_default()
	st := potrace_trace(p, bm)
	if st == nil {
		t.Fatal(libc.Error())
	} else if st.Status != 0 {
		t.Fatal(st.Status)
	}
	var tr trans_t
	trans_from_rect(&tr, float64(bm.W), float64(bm.H))
	bi := &info_s{
		Unit: 1,
	}
	iinfo := &imginfo_t{
		Pixwidth: bm.W, Pixheight: bm.H,
		Width: float64(bm.W), Height: float64(bm.H),
		Trans: tr,
	}

	svgOut, err := os.Create("stanford.svg")
	if err != nil {
		t.Fatal(err)
	}
	defer svgOut.Close()
	page_svg(bi, stdio.OpenFrom(svgOut), st.Plist, iinfo)
}
```

Test should again succeed and generate a `stanford.svg` files that can be viewed in the browser.

#### PDF

We promised to have PDF support as well, so let's try to add `backend_pdf.c` as well.

Of course, we will see similar issues with `info` here as well. 

But, before addressing these, we can also see that some
functions have conflicts in those two backends files. We already know how to solve it:

```yaml
files:
  # ...
  - name: backend_pdf.c
    idents:
      - name: unit
        rename: pdf_unit
      - name: ship
        rename: pdf_ship
```

Now, let's again add replacements for `info`:

```yaml
  - name: backend_pdf.c
    # ...
    replace:
      - old: 'func render0('
        new: 'func render0(info *info_s, '
      - old: 'func render0_opaque('
        new: 'func render0_opaque(info *info_s, '
      - old: 'func pdf_render('
        new: 'func pdf_render(info *info_s, '
      - old: 'func pdf_callbacks('
        new: 'func pdf_callbacks(info *info_s, '
      - old: 'func pdf_unit('
        new: 'func pdf_unit(info *info_s, '
      - old: 'func pdf_coords('
        new: 'func pdf_coords(info *info_s, '
      - old: 'func pdf_moveto('
        new: 'func pdf_moveto(info *info_s, '
      - old: 'func pdf_lineto('
        new: 'func pdf_lineto(info *info_s, '
      - old: 'func pdf_curveto('
        new: 'func pdf_curveto(info *info_s, '
      - old: 'func pdf_path('
        new: 'func pdf_path(info *info_s, '
      - old: 'func pdf_pageinit('
        new: 'func pdf_pageinit(info *info_s, '
      - old: 'func page_pdfpage('
        new: 'func page_pdfpage(info *info_s, '
      - old: 'func page_pdf('
        new: 'func page_pdf(info *info_s, '
      - old: 'func init_pdf('
        new: 'func init_pdf(info *info_s, '
      - old: 'func term_pdf('
        new: 'func term_pdf(info *info_s, '
      - old: 'render0_opaque(plist)'
        new: 'render0_opaque(info, plist)'
      - old: 'render0(plist)'
        new: 'render0(info, plist)'
      - old: 'pdf_unit(p'
        new: 'pdf_unit(info, p'
      - old: 'pdf_coords(p'
        new: 'pdf_coords(info, p'
      - old: 'pdf_moveto(*'
        new: 'pdf_moveto(info, *'
      - old: 'pdf_lineto(*'
        new: 'pdf_lineto(info, *'
      - old: 'pdf_curveto(*'
        new: 'pdf_curveto(info, *'
      - old: 'pdf_callbacks(fout'
        new: 'pdf_callbacks(info, fout'
      - old: 'pdf_pageinit(imginfo'
        new: 'pdf_pageinit(info, imginfo'
      - old: 'pdf_render(plist'
        new: 'pdf_render(info, plist'
      - old: 'pdf_path(&'
        new: 'pdf_path(info, &'
```

We notice that `dummy_xship` and `pdf_xship`. They are defined in `flate.c`, so we add it too.

Since we don't really need anything from that file except those two functions so let's skip `lzw_xship` that causes more errors.

```yaml
files:
  # ...
  - name: flate.c
    skip:
      - lzw_xship
```

And finally, let's add a few more lines to our test to generate a PDF as well:


```go
pdfOut, err := os.Create("stanford.pdf")
if err != nil {
    t.Fatal(err)
}
defer pdfOut.Close()
init_pdf(bi, stdio.OpenFrom(pdfOut))
page_pdf(bi, stdio.OpenFrom(pdfOut), st.Plist, iinfo)
term_pdf(bi, stdio.OpenFrom(pdfOut))
```

Test should succeed and generate a `stanford.pdf` files that should be more convenient to check.

## Improving the code

At this stage we have a working prototype, but the code might not be as nice as it could be.

Here are some issues we have:

- Fixed int types (`int64`) instead of Go `int`
- C names for types and functions
- Int types where we want bool types
- Pointer types where we want to have slices

We will now show how to address these issues.

To be on the same page, your config should look like this:

```yaml
vcs: https://github.com/skyrpex/potrace.git
branch: '1.15'
root: ./src
out: ../../.examples/potrace-go
package: gotrace
int_size: 8
ptr_size: 8
define:
  - name: VERSION
    value: '"dev"'
  - name: POTRACE
    value: '"potrace"'

idents:
  - name: potrace_dpoint_t
    alias: true
  - name: potrace_path_t
    alias: true
  - name: dpoint_t
    alias: true

files:
  - name: potracelib.c
  - name: progress.h
  - name: trace.c
    idents:
      - name: iprod
        rename: trace_iprod
      - name: interval
        rename: aux_interval
      - name: bezier
        rename: trace_bezier
  - name: decompose.c
  - name: curve.c
    skip:
      - potrace_path_s
  - name: bitmap.h
    replace:
      - old: 'bm_base(bm) = nil'
        new: 'bm.Map = nil'
  - name: bbox.c
    skip:
      - potrace_dpoint_s
  - name: auxiliary.h
    skip:
      - potrace_dpoint_s
    idents:
      - name: interval
        rename: aux_interval
  - name: bitmap_io.c
  - name: bitops.h
  - name: backend_svg.c
    replace:
      - old: 'func unit('
        new: 'func unit(info *info_s, '
      - old: 'func svg_moveto('
        new: 'func svg_moveto(info *info_s, '
      - old: 'func svg_rmoveto('
        new: 'func svg_rmoveto(info *info_s, '
      - old: 'func svg_lineto('
        new: 'func svg_lineto(info *info_s, '
      - old: 'func svg_curveto('
        new: 'func svg_curveto(info *info_s, '
      - old: 'func svg_path('
        new: 'func svg_path(info *info_s, '
      - old: 'func svg_jaggy_path('
        new: 'func svg_jaggy_path(info *info_s, '
      - old: 'func write_paths_opaque('
        new: 'func write_paths_opaque(info *info_s, '
      - old: 'func write_paths_transparent_rec('
        new: 'func write_paths_transparent_rec(info *info_s, '
      - old: 'func write_paths_transparent('
        new: 'func write_paths_transparent(info *info_s, '
      - old: 'func page_svg('
        new: 'func page_svg(info *info_s, '
      - old: 'func page_gimp('
        new: 'func page_gimp(info *info_s, '
      - old: 'unit(p)'
        new: 'unit(info, p)'
      - old: 'unit(p1)'
        new: 'unit(info, p1)'
      - old: 'unit(p2)'
        new: 'unit(info, p2)'
      - old: 'unit(p3)'
        new: 'unit(info, p3)'
      - old: 'svg_moveto(fout'
        new: 'svg_moveto(info, fout'
      - old: 'svg_rmoveto(fout'
        new: 'svg_rmoveto(info, fout'
      - old: 'svg_lineto(fout'
        new: 'svg_lineto(info, fout'
      - old: 'svg_curveto(fout'
        new: 'svg_curveto(info, fout'
      - old: 'svg_jaggy_path(fout'
        new: 'svg_jaggy_path(info, fout'
      - old: 'svg_path(fout'
        new: 'svg_path(info, fout'
      - old: 'write_paths_opaque(fout'
        new: 'write_paths_opaque(info, fout'
      - old: 'write_paths_transparent_rec(fout'
        new: 'write_paths_transparent_rec(info, fout'
      - old: 'write_paths_transparent(fout'
        new: 'write_paths_transparent(info, fout'
      - old: 'page_svg(fout'
        new: 'page_svg(info, fout'
  - name: flate.c
    skip:
      - lzw_xship
  - name: backend_pdf.c
    idents:
      - name: unit
        rename: pdf_unit
      - name: ship
        rename: pdf_ship
    replace:
      - old: 'func render0('
        new: 'func render0(info *info_s, '
      - old: 'func render0_opaque('
        new: 'func render0_opaque(info *info_s, '
      - old: 'func pdf_render('
        new: 'func pdf_render(info *info_s, '
      - old: 'func pdf_callbacks('
        new: 'func pdf_callbacks(info *info_s, '
      - old: 'func pdf_unit('
        new: 'func pdf_unit(info *info_s, '
      - old: 'func pdf_coords('
        new: 'func pdf_coords(info *info_s, '
      - old: 'func pdf_moveto('
        new: 'func pdf_moveto(info *info_s, '
      - old: 'func pdf_lineto('
        new: 'func pdf_lineto(info *info_s, '
      - old: 'func pdf_curveto('
        new: 'func pdf_curveto(info *info_s, '
      - old: 'func pdf_path('
        new: 'func pdf_path(info *info_s, '
      - old: 'func pdf_pageinit('
        new: 'func pdf_pageinit(info *info_s, '
      - old: 'func page_pdfpage('
        new: 'func page_pdfpage(info *info_s, '
      - old: 'func page_pdf('
        new: 'func page_pdf(info *info_s, '
      - old: 'func init_pdf('
        new: 'func init_pdf(info *info_s, '
      - old: 'func term_pdf('
        new: 'func term_pdf(info *info_s, '
      - old: 'render0_opaque(plist)'
        new: 'render0_opaque(info, plist)'
      - old: 'render0(plist)'
        new: 'render0(info, plist)'
      - old: 'pdf_unit(p'
        new: 'pdf_unit(info, p'
      - old: 'pdf_coords(p'
        new: 'pdf_coords(info, p'
      - old: 'pdf_moveto(*'
        new: 'pdf_moveto(info, *'
      - old: 'pdf_lineto(*'
        new: 'pdf_lineto(info, *'
      - old: 'pdf_curveto(*'
        new: 'pdf_curveto(info, *'
      - old: 'pdf_callbacks(fout'
        new: 'pdf_callbacks(info, fout'
      - old: 'pdf_pageinit(imginfo'
        new: 'pdf_pageinit(info, imginfo'
      - old: 'pdf_render(plist'
        new: 'pdf_render(info, plist'
      - old: 'pdf_path(&'
        new: 'pdf_path(info, &'
  - name: main.h
    go: backend.go
  - name: trans.c
  - name: progress_bar.c
  - name: hacks.go
    content: |
      package gotrace

      type backend_s struct{}
      type potrace_privstate_s struct{}
  - name: potrace_test.go
    content: |
      package gotrace

      import (
      	"io"
      	"io/ioutil"
      	"net/http"
      	"os"
      	"testing"

      	"github.com/gotranspile/cxgo/runtime/libc"
      	"github.com/gotranspile/cxgo/runtime/stdio"
      )

      func TestPotrace(t *testing.T) {
      	resp, err := http.Get("http://potrace.sourceforge.net/img/stanford.pbm")
      	if err != nil {
      		t.Fatal(err)
      	}
      	defer resp.Body.Close()

      	if resp.StatusCode != 200 {
      		t.Fatal(resp.Status)
      	}

      	tfile, err := ioutil.TempFile("", "potrace_")
      	if err != nil {
      		t.Fatal(err)
      	}
      	defer func() {
      		_ = tfile.Close()
      		_ = os.Remove(tfile.Name())
      	}()

      	_, err = io.Copy(tfile, resp.Body)
      	if err != nil {
      		t.Fatal(err)
      	}
      	tfile.Seek(0, io.SeekStart)

      	var bm *potrace_bitmap_t
      	e := bm_read(stdio.OpenFrom(tfile), 0.5, &bm)
      	if e != 0 {
      		t.Fatal(e)
      	}

      	p := potrace_param_default()
      	st := potrace_trace(p, bm)
      	if st == nil {
      		t.Fatal(libc.Error())
      	} else if st.Status != 0 {
      		t.Fatal(st.Status)
      	}
      	var tr trans_t
      	trans_from_rect(&tr, float64(bm.W), float64(bm.H))
      	bi := &info_s{
      		Unit: 1,
      	}
      	iinfo := &imginfo_t{
      		Pixwidth: bm.W, Pixheight: bm.H,
      		Width: float64(bm.W), Height: float64(bm.H),
      		Trans: tr,
      	}

      	svgOut, err := os.Create("stanford.svg")
      	if err != nil {
      		t.Fatal(err)
      	}
      	defer svgOut.Close()
      	page_svg(bi, stdio.OpenFrom(svgOut), st.Plist, iinfo)

      	pdfOut, err := os.Create("stanford.pdf")
      	if err != nil {
      		t.Fatal(err)
      	}
      	defer pdfOut.Close()
      	init_pdf(bi, stdio.OpenFrom(pdfOut))
      	page_pdf(bi, stdio.OpenFrom(pdfOut), st.Plist, iinfo)
      	term_pdf(bi, stdio.OpenFrom(pdfOut))
      }
```

### Using Go ints

Since we can now ensure that the code works correctly, we can experiment with the types. `cxgo` allows using Go `int`
instead of fixed-size ints generated by default. Add the following line to the config:

```yaml
use_go_int: true
```

Test should lead to exactly the same result. For some projects it might be false, so use this option with a caution.

### Better names

`cxgo` won't rename functions and types by default, but we can instruct it to do so, if desired.

For example, let's rename the functions use in tests and types in `potracelib.go`. We will also remove duplicates
like `type_t` and `type_s` to make the code cleaner.

```yaml
idents:
  # rename identifiers
  - name: potrace_progress_s
    rename: Progress
  - name: potrace_param_s
    rename: Param
  - name: potrace_word
    rename: Word
  - name: potrace_bitmap_s
    rename: Bitmap
  - name: potrace_dpoint_s
    rename: DPoint
  - name: potrace_curve_s
    rename: Curve
  - name: potrace_path_s
    rename: Path
  - name: potrace_state_s
    rename: State
  - name: trans_s
    rename: Trans
  - name: imginfo_s
    rename: ImgInfo
  - name: info_s
    rename: BackendInfo
  - name: dim_s
    rename: Dim
  - name: potrace_param_default
    rename: ParamDefault
  - name: potrace_trace
    rename: Trace
  # cleanup duplicate types
  - name: potrace_dpoint_t
    alias: true
  - name: potrace_path_t
    alias: true
  - name: dpoint_t
    alias: true
  - name: potrace_progress_t
    alias: true
  - name: potrace_param_t
    alias: true
  - name: potrace_bitmap_t
    alias: true
  - name: potrace_curve_t
    alias: true
  - name: potrace_state_t
    alias: true
  - name: trans_t
    alias: true
  - name: imginfo_t
    alias: true
  - name: info_t
    alias: true
  - name: dim_t
    alias: true
```

Also, don't forget to adjust test file and the `replace` directives. We can continue improving names in other files,
but there are other issues to address as well.

### Using Go bools

C code doesn't always use `bool` where it should. `cxgo` has a way to give a hint about those fields/arguments.

For example, in our `BackendInfo` struct type (former `info_s`) we have a name called `Debug`. If we check usages,
it's clear that it's used as a boolean value.

We can add the following to the `BackendInfo` definition to force `bool` type here:

```yaml
idents:
  # ...
  - name: info_s
    rename: BackendInfo
    fields:
      - name: debug
        type: bool
```

Note that `debug` is a C name, not a name that you see in Go (`Debug`).

We can do the same for other fields like `Compress`, `Opaque`, etc.

### Using slices

C doesn't have a notion of slices and unfortunately `cxgo` is not smart enough to detect those automatically.
But again, we have a way to give it a hint for a specific struct field or function argument.

We can start with promoting `Bitmap.Map` (`potrace_bitmap_t.map`) to a slice:

```yaml
idents:
  # ...
  - name: potrace_bitmap_s
    rename: Bitmap
    fields:
      - name: map
        type: slice
```

Unfortunately this won't compile: there are some unusual slice usages in the code.

First, this one:

```go
bm.Map = ([]Word)((*Word)(libc.Calloc(1, int(size))))
```

In fact, `cxgo` supports converting `calloc` to `make`, but it gots confused by `size` variable instead of `sizeof(T)`.
We can help it with `replace` (on `bitmap.h`):

```yaml
files:
  # ...
  - name: bitmap.h
    replace:
      - old: 'bm_base(bm) = nil'
        new: 'bm.Map = nil'
      - old: 'bm.Map = ([]Word)((*Word)(libc.Calloc(1, int(size))))'
        new: 'bm.Map = make([]Word, uintptr(size)/unsafe.Sizeof(Word(0)))'
```

The second usage is the following (in `bm_flip`):

```yaml
bm.Map = ([]Word)(&bm.Map[int64(bm.H-1)*int64(bm.Dy)])
bm.Dy = -dy
```

Potrace author doe a trick to flip the bitmap: he sets a pointer to the last element instead of the first and makes the
stride (row offset multiplier) negative. This way the indexing will automatically go backward.

Although it's a nice trick, Go won't allow this trickery on slices. So we have to disable this function and reimplement
it later in a different way (e.g. checking `Dy` sign and change how we index the slice).

Disabling can be done with `skip` on `bitmap.h`, and we will add a stub for into `hacks.go` since it's actually use in the code:

```yaml
files:
  # ...
  - name: bitmap.h
    skip:
      - bm_flip
  # ...
  - name: hacks.go
    content: |
      package gotrace

      // C flips bitmaps by using negative bitmap strides, which we cannot represent in Go with slices

      func bm_flip(bm *Bitmap) {
          // TODO: implement
      }

      type backend_s struct{}
      type potrace_privstate_s struct{}
```

And finally, the last issue is:

```go
newmap = (*Word)(libc.Realloc(unsafe.Pointer(&bm.Map[0]), int(newsize)))
if newmap == nil {
    goto error
}
bm.Map = ([]Word)(newmap)
```

`cxgo` doesn't translate `realloc` yet for slices, and even if it did, the local variable `newmap` still has a pointer type.
So we need to make at least two replacements here (plus one cosmetic):

```yaml
    replace:
      # ...
      - old: 'newmap = (*Word)(libc.Realloc(unsafe.Pointer(&bm.Map[0]), int(newsize)))'
        new: 'newmap = make([]Word, uintptr(newsize)/unsafe.Sizeof(Word(0))); copy(newmap, bm.Map)'
      - old: 'newmap  *Word'
        new: 'newmap []Word'
      - old: 'bm.Map = ([]Word)(newmap)'
        new: 'bm.Map = newmap'
```

This time it should compile, and the tests should still pass.