# Go Measurement Kit

Measurement Kit bindings for go.

**Attention** this is work on progress and is highly unstable.

Do not use it for anything serious, for the moment.

## Getting started

Run `./download-libs.sh` to download the prebuilt libraries for all platforms.

You can also specify just a single plaform with:

```
./download-libs.sh macos
```

Supported platforms are: `macos`, `mingw`


## Examples

See the `_examples/` directory.

## Windows

You can cross compile from macOS. To this end, please install the
mingw-w64-cxx11 toolchain formula from our [homebrew tap](
https://github.com/measurement-kit/homebrew-measurement-kit).

Once you have such toolchain, you should be able to get going by
running the following commands:

```
CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -x .

CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -x _examples/ndt/ndt.go
wine ndt.exe

CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -x _examples/web_connectivity/web_connectivity.go
wine web_connectivity.exe
```

It is anyway recommended to _also_ test using a real Windows system.
