#!/bin/bash

#脚本要存放在项目根目录
readonly GONELIST_ROOT=$(cd $(dirname ${BASH_SOURCE:-$0})/../; pwd -P)
source "${GONELIST_ROOT}/build/lib/var.sh"


# BUILD_VERSION不为空则切换到指定tag，否则按照最近tag-dirty走
read TAG_NUM LDFLAGS < <(GONELIST::SetVersion)

case "$1" in
  "release")
    bash ${GONELIST_ROOT}/build/lib/all-release.sh
    ;;
  "var")
    echo go build -o ${GONELIST_ROOT}/gonelist -ldflags "${LDFLAGS}" ${GONELIST_ROOT}/main.go
    ;;
  "build")
    go build -o ${GONELIST_ROOT}/gonelist -ldflags "${LDFLAGS}" ${GONELIST_ROOT}/main.go
    ;;
  "docker-local")
    Dockerfile=Dockerfile.local
    go build -o ${GONELIST_ROOT}/gonelist -ldflags "${LDFLAGS}" ${GONELIST_ROOT}/main.go
    ;&
  "docker")
    [ -n "$TAG_NUM" ] && build_arg="--build-arg VERSION=$TAG_NUM"
    docker build -t zhangguanzhang/gonelist:$TAG_NUM $build_arg \
      --build-arg LDFLAGS="-ldflags ${LDFLAGS}" -f ${Dockerfile:=Dockerfile} .
    ;;
  "clean")
    rm -rf ${GONELIST_ROOT:=/tmp}/release/*
    rm -f ${GONELIST_ROOT:=/tmp}/gonelist
    ;;
  *)
    echo -e "\t\033[1;31m must choose one to run \033[0m"
    exit 1
    ;;
esac
