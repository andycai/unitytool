#!/bin/bash

# 优化编译大小
# go build -ldflags="-s -w" main.go

# 打包 release 版本的 go 程序，windows 平台
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o unitool_serve_windows.exe main.go

# 打包 release 版本的 go 程序，linux 平台
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o unitool_serve_linux main.go

# 打包 release 版本的 go 程序，mac 平台
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o unitool_serve_mac main.go
