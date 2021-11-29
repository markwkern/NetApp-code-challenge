#!/bin/bash
if [[ ! -d ./build ]]; then
    mkdir -p build
fi
go build  -gcflags="all=-N -l" -o build/imagesrv imagesrv.go
