#!/bin/sh

output="./build"

rm -rf $output
mkdir $output

osarches=(
    "android/arm64"
    "darwin/arm64"
    "darwin/amd64"
    "linux/386"
    "linux/arm"
    "linux/arm64"
    "linux/amd64"
    "linux/mips"
    "linux/mips64"
    "linux/mips64le"
    "linux/mipsle"
    "linux/riscv64"
    "windows/386"
    "windows/arm64"
    "windows/amd64"
)

for osarch in "${osarches[@]}"; do
    os=$(echo $osarch | cut -d'/' -f1)
    arch=$(echo $osarch | cut -d'/' -f2)
    echo "Building for $os/$arch..."
    if [ $arch == *"64" ]; then
        CGO_ENABLED=0
    else
        CGO_ENABLED=1
    fi
    if [ "$os" == "windows" ]; then
        GOOS=$os GOARCH=$arch go build -o $output/jwc-$REL-$os-$arch.exe
    else
        GOOS=$os GOARCH=$arch go build -o $output/jwc-$REL-$os-$arch
    fi
done
