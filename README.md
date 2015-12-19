# GCI

GCI (**G**o **C**ontinuous **I**ntegration) is a build and CI system primarily for Go projects. It does all of its work inside a Docker container, so builds stay consistent and you don't have to set anything up on your machine except docker.

# Pipelines

Like [wercker](http://wercker.com), you tell gci what to do by defining _pipelines_. Each pipeline is made up of 1 or more _steps_.

# Build File

You specify one or more pipelines in a YAML file called the build file. Each pipeline uses the steps that you've decided to import.

```yaml
version: 0.0.1
# a build file can define variables to be used later in the file.
# variables names must match the following regex: [a-zA-Z0-9_]+.
# all variables are available for use in the task definitions below
vars:
  # each variable has a default, but can be overridden by an environment variable
  - name: Version
    env: MYAPP_VERSION
    default: 0.0.2
  # variables without an 'env' are constant
  - name: AppName
    default: my-app
  - name: DockerHost
    env: DOCKER_HOST
    default: /var/run/docker.sock
# define what pipeline steps you'll need
steps:
  - name: glide-up
    version: 0.8.1
  - name: go-build
    version: 1.5.2
# set up a few pipelines
pipelines:
  # build the go program that is in the same working directory as this file
  - name: build
    description: build the program
    steps:
      # first run the glide-up step
      - name: glide-up
      # then run the build
      - name: go-build
        params:
          - name: Out
            value: {{.AppName}}.{{.Version}}
          - name: VendorExperiment
            value: 1
          - name: CGO
            value: 0
  # test the go package in the same directory as this file, and all sub packages from here too
  - name: test  
    description: test the program
    steps:
      - name: glide-up
      - name: go-test
        params:
          - name: Packages
            value: ./...
```

Assuming you saved this build file to my.yaml, you can run a build with the following command:

```console
gci -f my.yaml run build
```

Note that gci looks for gci.yaml in the current working directory by default.


# Building Your Own Steps

A pipeline step is basically a manifest file that tells gci how to run a docker container. It contains:

- A version (for gci forward compatibility)
- A Docker image name
- A command to run in the container
- A list of parameters that must be passed from a build file
- A list of volume mounts that gci should make when running the container
- A list of environment variables

Here's a sample step manifest file:

```yaml
version: 0.0.1
image: quay.io/deis/go-dev:0.3.0
command: glide up
volumes:
  - host: {{.PWD}}
    container: /pwd
```

# CI Server

TODO: this CI server will download some code and execute one or more pipelines on it.
