GOOS=darwin GOARCH=amd64 go build && mv ./goinstant ./bin/goinstant-macos \
&& GOOS=linux GOARCH=amd64 go build && mv ./goinstant ./bin/goinstant-linux \
&& GOOS=windows GOARCH=amd64 go build && mv ./goinstant.exe ./bin/goinstant.exe \
&& go clean