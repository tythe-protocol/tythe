[![Build Status](https://travis-ci.com/tythe-protocol/tythe.svg?branch=master)](https://travis-ci.com/tythe-protocol/tythe)

# Get the code

```
git clone https://github.com/tythe-protocol/tythe
cd tythe
```

# Important Note on Go Modules

Tythe uses Go modules. Either check out to some directory ***other than GOPATH*** (recommended), or else use `GO111MODULE=on` whenever building Tythe.

# Build the CLI

```
go test ./...
go build -tags 'dev' ./cmd/tythe
./tythe
```

# Build the Server

### Development

```
cd cmd/webtythe
go build

# Runs the API server on :8080
./webtythe &

cd ui

# Runs the UI on :3030
npm run start
```

### Production

```
cd cmd/webtythe
./buildprod.sh

# Runs the API and UI on :8080
./webtythe
```
