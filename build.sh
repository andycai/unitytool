#!/bin/bash

# 设置通用的编译参数
VERSION=$(git describe --tags --always)
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')
LDFLAGS='-s -w -X "main.Version='${VERSION}'" -X "main.BuildTime='${BUILD_TIME}'"'

# 打包 release 版本的 go 程序，windows 平台
echo "Building for Windows..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -trimpath -o unitool_serve_windows.exe
if [ $? -eq 0 ]; then
    echo "Windows build successful"
    # 使用 UPX 进一步压缩（如果安装了 UPX）
    if command -v upx &> /dev/null; then
        upx --best --lzma unitool_serve_windows.exe
    fi
fi

# 打包 release 版本的 go 程序，linux 平台
echo "Building for Linux..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -trimpath -o unitool_serve_linux
if [ $? -eq 0 ]; then
    echo "Linux build successful"
    # 使用 UPX 进一步压缩（如果安装了 UPX）
    if command -v upx &> /dev/null; then
        upx --best --lzma unitool_serve_linux
    fi
fi

# 打包 release 版本的 go 程序，mac 平台
echo "Building for macOS..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -trimpath -o unitool_serve_mac
if [ $? -eq 0 ]; then
    echo "macOS build successful"
    # macOS 二进制文件不建议使用 UPX 压缩，可能会导致签名问题
fi

echo "Build process completed"

# 显示编译后的文件大小
echo -e "\nFile sizes:"
ls -lh unitool_serve_*
