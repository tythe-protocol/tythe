[![Build Status](https://travis-ci.com/tythe-protocol/tythe.svg?branch=master)](https://travis-ci.com/tythe-protocol/tythe)

# Dependencies

* Go 1.11+
* Node 10+

# Get the Code

```
git clone https://github.com/tythe-protocol/tythe
git clone https://github.com/tythe-protocol/z_test1
git clone https://github.com/tythe-protocol/z_test2
cd tythe
```

# Important Note on Go Modules

Tythe uses Go modules. Either check out to some directory ***other than GOPATH*** (recommended), or else use `GO111MODULE=on` whenever building Tythe.

# Build the CLI

```
go build ./cmd/tythe
./tythe
```

# Build the Server

### Only needed the first time

```
cd cmd/webtythe/ui
npm install
cd -
```

### Needed every time

```
cd cmd/webtythe
go run build.go

# Run server on :8080
./webtythe
```

# Run Tests

Note that this will fail until the server has been built once. Stupid codegen.

```
cd tythe
go test ./...
```
