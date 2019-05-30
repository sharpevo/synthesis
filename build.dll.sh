#!/bin/bash
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc GOOS=windows GOARCH=amd64 go build -v -buildmode=c-shared -o build/interpreter.dll #interpreter.c interpreter.go
# no interpreter.cpp
# or CXX=x86_64-w64-mingw32-g++
# https://github.com/kubeapps/kubeapps/issues/63
