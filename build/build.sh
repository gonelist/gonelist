#!/bin/bash

[ -n "$DEBUG" ] && set -x

#脚本要存放在项目根目录
readonly GONELIST_ROOT=$(cd $(dirname ${BASH_SOURCE:-$0})/../; pwd -P)
source "${GONELIST_ROOT}/build/lib/var.sh"

# BUILD_VERSION不为空则切换到指定tag，否则按照最近tag-dirty走
read TAG_NUM LDFLAGS < <(GONELIST::SetVersion)

echo go build -o ${GONELIST_ROOT}/gonelist -ldflags "${LDFLAGS}" ${GONELIST_ROOT}/main.go

# 前端dist版本无值则取TAG版本号
[ -z "$DIST_VERSION" ] && DIST_VERSION=$TAG_NUM

case "$1" in
  "release") # checkout到tag构建完再checkout回来
    bash ${GONELIST_ROOT}/build/lib/all-release.sh
    ;;
  "build") #使用master构建测试版本
    if [ -z `command -v go ` ];then
      echo go is not in PATH
      exit 1
    fi
    go build -o ${GONELIST_ROOT}/gonelist -ldflags "${LDFLAGS}" ${GONELIST_ROOT}/main.go
    ;;
  "docker-local") #使用本地编译二进制文件打包docker和dist
    Dockerfile=Dockerfile.local
    go build -o ${GONELIST_ROOT}/gonelist -ldflags "${LDFLAGS}" ${GONELIST_ROOT}/main.go
    ;&
  "docker") #使用容器编译和打包dist
    [ -n "$TAG_NUM" ] && build_arg="--build-arg VERSION=$DIST_VERSION"
    docker build -t zhangguanzhang/gonelist:$TAG_NUM $build_arg \
      --build-arg LDFLAGS="${LDFLAGS}" -f ${Dockerfile:=Dockerfile} .
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

