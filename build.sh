#!/bin/bash
set -e

TAG=${1:-latest}
BUILDKIT_ENABLED=0

DOCKER_SERVER_VERSION=$(docker version --format '{{ .Server.Version }}')
DOCKER_CLIENT_VERSION=$(docker version --format '{{ .Client.Version }}')
if [[ $DOCKER_SERVER_VERSION < 18.09 || $DOCKER_CLIENT_VERSION < 18.09 ]]; then
  BUILDKIT_ENABLED=0
fi

DOCKER_BUILDKIT=$BUILDKIT_ENABLED docker image build --target test .
DOCKER_BUILDKIT=$BUILDKIT_ENABLED docker image build -t bjornmagnusson/pi-led:${TAG} .
