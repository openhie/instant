# goinstant

This is a Go app and can be built as a native binary for any operating system. Static assets and templates are built using packr2.

## Dev prerequisites

* Install go, [see here](https://golang.org/doc/install). For Ubuntu you might want to use the go snap package, [see here](https://snapcraft.io/install/go/ubuntu).
* Install packr2: **Outside** of the goinstant folder (so that it doesn't get installed as a module) run: `go get -u github.com/gobuffalo/packr/v2/packr2`
* Add go binaries to you system $PATH, on ubuntu: Add `export PATH=$PATH:$HOME/go/bin` to the end of your ~/.bashrc file. To use this change immediately source it: `source ~/.bashrc`
* install dependencies, run this from the goinstant folder: `go get`

## Running and building

Run the app using `go run goinstant.go` or build the binary using `go build`. To build releases, create a tag and upload the binaries built. A convenience bash script is included to build binaries.

Note: this script won't work if you have a `goistant/data/` folder that gets created when upping the docker-compose files through the go app. Delete this first: `sudo rm -r goisntant/data`

```sh
bash ./buildreleases.sh
git tag 0.0.1
git push origin 0.0.1
# then upload the binaries.
```
