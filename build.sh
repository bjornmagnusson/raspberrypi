#!/bin/bash
set -e

TAG=${1:-latest}
DOCKER_BUILDKIT=1 docker image build --target test .
DOCKER_BUILDKIT=1 docker image build -t bjornmagnusson/pi-led .
