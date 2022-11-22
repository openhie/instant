#!/bin/bash

if [ "$1" == "linux" ]; then
    curl -L https://github.com/openhie/package-starter-kit/releases/download/latest/gocli-linux -o goinstant
    chmod +x goinstant

elif [ "$1" == "macos" ]; then
    curl -L https://github.com/openhie/package-starter-kit/releases/download/latest/gocli-macos -o goinstant
    chmod +x goinstant

elif [ "$1" == "windows" ]; then
    curl -L https://github.com/openhie/package-starter-kit/releases/download/latest/gocli.exe -o goinstant.exe
    chmod +x ./goinstant.exe

else
    echo 'Usage: ./get-cli.sh "linux|macos|windows"'
fi
