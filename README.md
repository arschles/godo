# Godo

Godo ("Go do") is a build tool for executing consistent, reproducible builds inside Docker containers.

The only dependencies to use godo are a working Docker installation and the `godo` binary installed.

# Downloads

- [Mac OS X 64 Bit](https://github.com/arschles/godo/releases/download/v0.3.0/godo_darwin_amd64)
- [Linux 64 Bit](https://github.com/arschles/godo/releases/download/v0.3.0/godo_linux_amd64)
- [Windows 64 Bit](https://github.com/arschles/godo/releases/download/v0.3.0/godo_windows_amd64.exe)
- [Others](https://github.com/arschles/godo/releases/tag/v0.3.0)

## Usage

Godo loosely resembles [`make`](https://www.gnu.org/software/make/), except the commands and scripts that you tell it to execute must be inside a Docker conainer.

Like `make`, it operates on _build targets_ and ships with the following targets pre-configured, out of the box:

- `build` - build your Go code inside a container (soon to be deprecated)
- `docker-build` - build an image with your Go program inside it
- `docker-push` - push your image to a Docker repository

### Configuration

Project structures and settings vary significantly, so each of those targets are configurable in a `godo.yaml` file (or `godo.yml`). See [Godo's own godo.yaml file](https://github.com/arschles/godo/blob/master/godo.yaml) for an example

#### Features Not Yet Implemented

The following high level features are planned, but not yet implemented:

- Target dependencies (https://github.com/arschles/godo/issues/18)
- Custom targets (https://github.com/arschles/godo/issues/17)
- CI Server (https://github.com/arschles/godo/issues/19)

## Development

Writing code for Godo is simple. Since it's a build tool, it can bootstrap itself by running `godo build` in the root of the repo.

Assuming you don't have a `godo` binary available, you can build it with the standard `go` toolchain as well. Here's what you'll need:

- Go 1.6+
- [Glide](https://github.com/Masterminds/glide) 0.8+

Once you ensure that you have both dependencies, simply:

```console
glide install
go install
```

Ensure that you have `$GOPATH/bin` somewhere on your `$PATH`, and now you can execute `godo` targets.

When you're ready to write code, simply fork and open a pull request.
