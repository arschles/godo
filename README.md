# canta

Canta is a build and CI system primarily for Go projects. It provides a full environment for building, testing, serving and doing other tasks during your development cycle, all within docker containers.

# Targets

Like [GNU Make](https://www.gnu.org/software/make/), canta's building blocks are _targets_. Each target defined one piece of work that accomplishes part of the build. A target executes one or more commands inside a Docker container running the image of your choice. Targets can be parameterized by environment variables, and they can depend on each other.

All targets are defined in a `canta.yml` or `canta.yaml` file in the current directory.

# Plugins

In addition to targets, Canta also provides _plugins_. A plugin is simply a collection of one or more targets that can be added to your project with the `canta plugin add` command. Plugins are stored in the `.canta/plugins` directory in the current working directory.
