#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

readonly GONELIST_ROOT=$(cd "$(dirname ${BASH_SOURCE:-$0})/../"; pwd -P)
export CGO_ENABLED=0

source ${GONELIST_ROOT}/build/lib/var.sh

# BUILD_VERSION不为空则切换到指定tag，否则按照最近tag-dirty走
GONELIST::SetVersion

cd $GONELIST_ROOT

case "$1" in :
  "all-os":
    source ${GONELIST_ROOT}/build/lib/all-release.sh
    ;;
  "master":
    go build -o ${GONELIST_ROOT}/gonelist  ${LDFLAGS} main.go
    ;;
  "docker":
    docker build -t zhangguanzhang/gonelist:test \
      $DOCKER_Version \
      --build-arg LDFLAGS="${LDFLAGS}" .
    ;;
  *:
    exit 1
  esac
