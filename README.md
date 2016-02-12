# GCI

GCI (**G**o **C**ontinuous **I**ntegration) is a build and CI system build specifically for Go projects. By default, all build tasks are done inside of Docker containers, so everything is consistent across machines.

The only dependencies to use gci are a working Docker installation and the `gci` binary installed.

## Usage

GCI is both a build tool and Continuous Integration server.

As a build tool, it loosely resembles [`make`](https://www.gnu.org/software/make/), but of course it's designed specifically around Go projects. Like `make`, it operates on _build targets_. GCI ships with the following targets pre-configured, out of the box:

- `build` - build your Go code inside a container
- `test` - test your code inside a container
- `docker-build` - build an image with your Go program inside it
- `docker-push` - push your image to a Docker repository

### Configuration

Project structures and settings vary significantly, so each of those targets are configurable in a `gci.yaml` file (or other name if you choose). See [gci's own gci.yaml file](https://github.com/arschles/gci/blob/master/gci.yaml) for an example

#### Features Not Yet Implemented

The following high level features are planned, but not yet implemented:

- Target dependencies (https://github.com/arschles/gci/issues/18)
- Custom targets (https://github.com/arschles/gci/issues/17)
- CI Server (https://github.com/arschles/gci/issues/19)

## Development

Writing code for GCI is simple. Since it's a build tool, it can bootstrap itself by running `gci build` in the root of the repo.

Assuming you don't have a `gci` binary available, you can build it with the standard `go` toolchain as well. Here's what you'll need:

- Go 1.5+
- [Glide](https://github.com/Masterminds/glide) 0.8+

Once you ensure that you have both dependencies, simply:

```console
glide up
GO15VENDOREXPERIMENT=1 go install
```

Ensure that you have `$GOPATH/bin` somewhere on your `$PATH`, and now you can execute `gci` targets.

When you're ready to write code, simply fork and open a pull request.
