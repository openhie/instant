GOOS=darwin GOARCH=amd64 packr2 build && mv ./goinstant ./bin/goinstant-macos \
  && GOOS=linux GOARCH=amd64 packr2 build && mv ./goinstant ./bin/goinstant-linux \
  && GOOS=windows GOARCH=amd64 packr2 build && mv ./goinstant.exe ./bin/goinstant.exe \
  && packr2 clean