# goinstant

This is a Go app and can be built as a native binary for any operating system. Static assets and templates are built using packr2. 

Run it using `go run goinstant.go` or build the binary using `go build`. To build releases, create a tag and upload the binaries built. A convenience bash script is included to build binaries. 

```sh
bash ./buildreleases.sh
git tag 0.0.1
git push origin 0.0.1
# then upload the binaries.
```