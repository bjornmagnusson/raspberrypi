DOCKER_TAG   ?= latest

.PHONY: build
build:
	@./build.sh '${DOCKER_TAG}'

test:
	@docker-compose -f docker-compose.dev.yml run api -demo=true -num=3