#!/usr/bin/env bash

set -o errexit
set -o pipefail
export GONELIST_ROOT TAG_NUM LDFLAGS CGO_ENABLED=0
readonly GONELIST_ROOT=$(cd "$(dirname ${BASH_SOURCE:-$0})"; pwd -P)


source ${GONELIST_ROOT}/build/lib/var.sh

# BUILD_VERSION不为空则切换到指定tag，否则按照最近tag-dirty走

read TAG_NUM LDFLAGS < <(GONELIST::SetVersion)

cd $GONELIST_ROOT

case "$1" in
  "release")
    bash ${GONELIST_ROOT}/build/lib/all-release.sh
    ;;
  "build")
    bash -c "go build -o ${GONELIST_ROOT}/gonelist  ${LDFLAGS} ${GONELIST_ROOT}/main.go"
    ;;
  "docker-local")
    Dockerfile=Dockerfile.local
    ;&
  "docker")
    [ -n "$TAG_NUM" ] && build_arg="--build-arg VERSION=$TAG_NUM"
    docker build -t zhangguanzhang/gonelist:test $build_arg \
      --build-arg LDFLAGS="${LDFLAGS}" -f ${Dockerfile:=Dockerfile} .
    ;;
  "clean")
    rm -rf ${GONELIST_ROOT:=/tmp}/release/*
    rm -f ${GONELIST_ROOT:=/tmp}/gonelist
    ;;
  *)
    echo -e "\t\033[1;31m must choose one to run \033[0m"
    exit 1
  esac
