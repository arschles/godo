# it's ironic right now that a build tool needs a makefile to build itself.
# TODO: bootstrap with canta
GO_DEV_IMG := quay.io/deis/go-dev:0.3.0
FULL_PATH_IMG := /go/src/github.com/arschles/canta
DOCKER_CMD := docker run -e GO15VENDOREXPERIMENT=1 -e CGO_ENABLED=0 --rm -v ${PWD}:${FULL_PATH_IMG} -w ${FULL_PATH_IMG} ${GO_DEV_IMG}
VERSION ?= 0.0.1
DOCKER_HOST ?= ${DOCKER_HOST}

bootstrap:
	${DOCKER_CMD} glide up

build:
	${DOCKER_CMD} go build -o canta

docker-build:
	docker build -t quay.io/arschles/gbs:${VERSION} .

docker-push:
	docker push quay.io/arschles/gbs:${VERSION}
