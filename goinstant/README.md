# goinstant

This is a Go desktop app and is provided as a native binary for the AMD64 architecture on Windows, macOS, and Linux.

> Warning: This app is not meant to be used for container and cluster management in production or with sensitive data. It is meant for demos, training sessions, and by developers. In production and with sensitive data, administrators should use the purpose-built tools like the Docker and Kubernetes CLIs to manage resources directly and according to best practices which are outside the scope of this app.

## How it Works

> The prototype of goinstant was CLI only and that is now disabled. If a CLI is desired, it it suggested to use the yarn commands provided for running from the command line. The yarn CLI will be maintained.

The desktop is rendered using web pages.

* Static assets are bundled into the executable.
* Web pages are served by Go HTTP server.
* Calls from JS go to the Go API, which in turn calls Go functions. While Go apps can easily render and serve pages directly, this separation of concerns exists so that a contributor may add a different Web framework like React or Vue, and build a new frontend while the API remains in Go.
* Go functions generally invoke the shell on the user's OS. There are a few exceptions. 
* On init and after accepting the disclaimer, this repo is cloned to the system. (This is done with a git-compatible library, git doesn't need to be installed.)
* The console uses server-side events. These are one-way events sent to the client, not bidirectional like websockets. There is a function which can be called `consoleSender` to send events to the console.

## Security

This desktop app is meant as a prototype and may change. This app resides in userspace but it invokes the command line for containers and clusters. The apps it invokes, Docker and Kubernetes CLI, launch and manage containers and may have admin/root privileges.

Therefore, this app is not meant to be used for container and cluster management in production or with sensitive data. It is meant for demos, training sessions, and by developers. In production and with sensitive data, administrators should use the purpose-built tools like the Docker and Kubernetes CLIs to manage resources directly and according to best practices which are outside the scope of this app.

## Developers

### Dev prerequisites

* Install go, [see here](https://golang.org/doc/install). For Ubuntu you might want to use the go snap package, [see here](https://snapcraft.io/install/go/ubuntu).
* Install packr2: **Outside** of the goinstant folder (so that it doesn't get installed as a module) run: `go get -u github.com/gobuffalo/packr/v2/packr2`
* Add go binaries to you system $PATH, on ubuntu: Add `export PATH=$PATH:$HOME/go/bin` to the end of your ~/.bashrc file. To use this change immediately source it: `source ~/.bashrc`
* Install dependencies, run this from the goinstant folder: `go get`

### Running

Static assets and templates are built using pkger. For development, run the app using `pkger && go run *.go`. `pkger` must be run if you change the static assets.

### Building

To build releases, create a tag and upload the binaries. A convenience bash script is included to build binaries.

> Note that ARM builds are not yet supported in Go 1.15, but such builds can be supported in future for Linux and macOS.

Note: this script won't work if you have a `goinstant/data/` folder that gets created when when starting up the docker-compose files through the go app. Delete this first: `sudo rm -r goinstant/data`

```sh
bash ./buildreleases.sh
git tag 0.0.1
git push origin 0.0.1
# then upload the binaries to GitHub
```

