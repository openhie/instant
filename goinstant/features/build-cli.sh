#!/bin/bash

cp -r ../../goinstant ../../goinstant-tmp
cp test-conf/* ../../goinstant-tmp
cd ../../goinstant-tmp

GOOS=darwin GOARCH=amd64 go build -o ../goinstant/features/test-platform-macos
GOOS=linux GOARCH=amd64 go build -o ../goinstant/features/test-platform-linux
GOOS=windows GOARCH=amd64 go build -o ../goinstant/features/test-platform.exe
go clean

cd ..
rm -rf goinstant-tmp
