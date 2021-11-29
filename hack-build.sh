#!/bin/bash
go build  -gcflags="all=-N -l" -o build/imagesrv imagesrv.go
