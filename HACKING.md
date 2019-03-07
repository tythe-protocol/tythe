[![Build Status](https://travis-ci.com/tythe-protocol/tythe.svg?branch=master)](https://travis-ci.com/tythe-protocol/tythe)

```
git clone https://github.com/tythe-protocol/tythe
cd tythe

# Development
GO111MODULE=on go build -tags 'dev' ./cmd/tythe
GO111MODULE=on go test ./...
./tythe

# UI Development
./tythe serve # serves API on :8080
cd cmd/tythe/ui
npm run start # serves UI on :3000

# Production
./build-prod.sh
./tythe serve # serves API and UI on 8080
```
