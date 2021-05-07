# goinstant

This is a Go CLI app and is provided as a native binary for the AMD64 architecture on Windows, macOS, and Linux.

> Warning: This app is not meant to be used for container and cluster management in production or with sensitive data. It is meant for demos, training sessions, and by developers. In production and with sensitive data, administrators should use the purpose-built tools like the Docker and Kubernetes CLIs to manage resources directly and according to best practices which are outside the scope of this app.

## Usage

Download the latest release for your OS.

On Unix-like operating systems, you must add execute permissions, ie. `chmod +x goinstant-linux`.

Without arguments, the CLI defaults to interactive mode. The CLI can also be used non-interactively as so:
```
Commands: 
	help 		this menu
	docker		usage: docker <package> <state> e.g. docker core init
	kubernetes	usage: k8s/kubernetes <package> <state>, e.g. k8s core init
	utils		usage: utils ig load <url> <fhirserver>, ig examples <url> <fhirserver>
```


## FHIR NPM protobuf

If changing the proto, it must be regenerated. Otherwise, regenerating the proto is not necessary.
```sh
protoc --go_out=fhirnpmproto --go_opt=paths=source_relative --proto_path=fhirnpmproto indexjson.proto
```



## Security

This desktop app is meant as a prototype and may change. This app resides in userspace but it invokes the command line for containers and clusters. The apps it invokes, Docker and Kubernetes CLI, launch and manage containers and may have admin/root privileges.

Therefore, this app is not meant to be used for container and cluster management in production or with sensitive data. It is meant for demos, training sessions, and by developers. In production and with sensitive data, administrators should use the purpose-built tools like the Docker and Kubernetes CLIs to manage resources directly and according to best practices which are outside the scope of this app.

## Developers

### Dev prerequisites

* Install go, [see here](https://golang.org/doc/install). For Ubuntu you might want to use the go snap package, [see here](https://snapcraft.io/install/go/ubuntu).
* Add go binaries to you system $PATH, on ubuntu: Add `export PATH=$PATH:$HOME/go/bin` to the end of your ~/.bashrc file. To use this change immediately source it: `source ~/.bashrc`
* Install dependencies, run this from the goinstant folder: `go mod tidy`

### Running

For development, run the app using `go run *.go`.

### Building

To build releases, create a tag and upload the binaries. A convenience bash script is included to build binaries.

> Note that ARM builds are not yet supported in Go 1.15, but such builds can be supported in future for Linux and macOS.

```sh
bash ./buildreleases.sh
git tag 0.0.1
git push origin 0.0.1
# then upload the binaries to GitHub
```

