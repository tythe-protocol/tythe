#!/bin/sh
GO111MODULE=on
cd ui
npm run build && go generate
result=$?
cd - > /dev/null
if [ $result -eq 0 ]; then
    go build -tags 'release'
fi
rm ui/*vfsdata.go
