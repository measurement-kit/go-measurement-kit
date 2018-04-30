# Go Measurement Kit

Measurement Kit bindings for go.

**Attention** this is work on progress and is highly unstable.

Do not use it for anything serious, for the moment.

## Examples

See the `_examples/` directory.

## Windows

It is currently not possible to cross compile from macOS. Even
though we can get rid of most undefined symbols, sadly the symbols
for GNU libstdc++ threads are missing and cannot be generated because
Homebrew does not compile the mingw-w64 package with the POSIX
thread model. Hence, for now, you really need a
[MSYS2](https://www.msys2.org/) system. The following instructions
assume that you are inside the `x86_64` MSYS2 shell (different from the
normal shell in that by default you're using the `x86_64` toolchain).

```
go build -x .
(cd _examples/ndt && go build -x .)
(cd _examples/web_connectivity && go build -x .)
```
