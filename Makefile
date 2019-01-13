DOCKER_TAG   ?= latest

.PHONY: build
build:
	@./build.sh '${DOCKER_TAG}'
