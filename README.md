# canta

Canta is a build and CI system primarily for Go projects. It provides a full environment for building, testing, serving and doing other tasks during your development cycle, all within docker containers.

# Targets

Like [GNU Make](https://www.gnu.org/software/make/), canta's building blocks are _targets_. Each target defined one piece of work that accomplishes part of the build. A target executes one or more commands inside a Docker container running the image of your choice. All targets are defined in a `canta.yml` or `canta.yaml` file in the current directory.

Targets can depend on each other (like in Make) and can be parameterized by environment variables.

# Plugins

In addition to targets, Canta also provides _plugins_. A plugin is simply a collection of one or more targets that can be added to your project with the `canta plugin add` command. Plugins are stored in the `.canta/plugins` directory in the current working directory.

# CI Server

Canta can be run locally, but since all targets run exclusively in docker containers, they're portable and consistent. Canta can also run as a server that can fetch your code, run the targets that you want it to, and report on the output of each.
