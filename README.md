# Godo

Godo ("Go do") is a build tool for executing consistent, reproducible builds inside Docker containers.

The only dependencies to use godo are a working Docker installation and the `godo` binary installed.

# Downloads

- [Mac OS X 64 Bit](https://github.com/arschles/godo/releases/download/v0.3.0/godo_darwin_amd64)
- [Linux 64 Bit](https://github.com/arschles/godo/releases/download/v0.3.0/godo_linux_amd64)
- [Windows 64 Bit](https://github.com/arschles/godo/releases/download/v0.3.0/godo_windows_amd64.exe)
- [Others](https://github.com/arschles/godo/releases/tag/v0.3.0)

# Motivation

Over the last year, there have been a [few](https://blog.docker.com/2016/09/docker-golang/) [posts](https://www.iron.io/the-easiest-way-to-develop-with-go%E2%80%8A-%E2%80%8Aintroducing-a-docker-based-go-tool/) explaining how to build Go code using Docker containers. The key features of this workflow are as follows:

1. You don't need Go installed on your machine
2. You can compile with any version of Go, without changing your host's setup
3. You can move your build environment anywhere, as long as you have Docker

I started writing a tool to make it dead-simple to compile Go programs with this method, but quickly realized that Docker containers can be used to provide _consistent and portable builds_, everywhere Docker is installed.

In addition to creating a tool to make the Docker-based workflow dead-simple, I also created a yaml-based file format that tells the `godo` CLI what builds it can run.

## Usage

Godo loosely resembles [`make`](https://www.gnu.org/software/make/), except the commands and scripts that you tell it to execute must be inside a Docker conainer.

Like `make`, it operates on _build targets_ and ships with the following targets pre-configured, out of the box:

- `build` - build your Go code inside a container (soon to be deprecated, in favor of a custom target)
- `docker-build` - build an image with your Go program inside it
- `docker-push` - push your image to a Docker repository
- `custom` - execute an arbitrary script inside a Docker container

### Configuration

Project structures and settings vary significantly, so each of those targets are configurable in a `godo.yaml` file (or `godo.yml`). See [Godo's own godo.yaml file](https://github.com/arschles/godo/blob/master/godo.yaml) for an example

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
