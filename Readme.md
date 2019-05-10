# Go bindings for C++ Measurement Kit code

Measurement Kit bindings for Go. The code in this repository exposes
some APIs of [Measurement Kit's C++ implementation](
https://github.com/measurement-kit/measurement-kit) to Go apps.

Not to be confused with [measurement-kit/engine](
https://github.com/measurement-kit/engine), which contains a Measurement Kit
implementation in Go.

**Attention** this is work on progress and is highly unstable.

Do not use it for anything serious, for the moment.

## Getting started

Install MaxMind databases using:

```bash
./script/download-mmdb.sh
```

### macOS

Install Measurement Kit using brew:

```bash
brew tap measurement-kit/measurement-kit
brew install measurement-kit
```

If you've already installed `measurement-kit`, do:

```bash
brew upgrade
```

to make sure you're on the latest released version.

Then you're all set. Just `go get -v ./...` as usual.

### MingGW

Install Measurement Kit using brew:

```bash
brew tap measurement-kit/measurement-kit
brew install mingw-w64-measurement-kit
```

If you've already installed `mingw-w64-measurement-kit`, do:

```bash
brew upgrade
```

to make sure you're on the latest released version.

Then you're all set. Because you're cross compiling you need to provide
more environment variables to make the build work:

```bash
CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++                           \
  CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go get -v ./...
```

To test binaries, you can use wine:

```bash
wine $GOPATH/bin/windows_amd64/web_connectivity.exe 
```

It is recommended to _also_ test using a real Windows box.

### Linux

We have a Docker container. Build the container with:

```bash
docker build -t gomkbuild .
```

Enter into the container with:

```bash
docker run -it -v`pwd`:/gomkbuild -w/gomkbuild gomkbuild
```

Then you're all set. Just `go get -v ./...` as usual.

## Examples

See the `_examples/` directory.
