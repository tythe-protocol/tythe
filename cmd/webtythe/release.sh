#!/bin/sh
go run build.go --prod
scp webtythe_linux_amd64 aa@aaronboodman.com:webtythe
ssh aa@aaronboodman.com tythe.dev/deploy.sh
