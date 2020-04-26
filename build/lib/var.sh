#!/usr/bin/env bash

#GONELIST_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
git=(git --work-tree "${GONELIST_ROOT}")


GONELIST::SetVersion(){
  # 不为空则切换到指定tag版本编译
  [ -n $"BUILD_VERSION" ] && {
    "${git[@]}" checkout ${BUILD_VERSION}
  } || { # 默认取最新的tag作为版本号
    # repo上最新tag的commitID
    TAG_COMMITID=$("${git[@]}" rev-list --tags --max-count=1)
    # 取最新tag名
    BUILD_VERSION=$("${git[@]}"  describe --tags ${TAG_COMMITID})
  }

  DOCKER_Version="--build-arg VERSION=${BUILD_VERSION}"

  COMMIT_ID=$("${git[@]}" rev-parse HEAD)
  BUILD_DATE=$(date +'%Y-%m-%dT%H:%M:%S')
  GIT_TREE_STATE=$("${git[@]}" status --porcelain 2>/dev/null)

  [ -n "${GIT_TREE_STATE}" ] && {
    # 不为空则为master分支修改代码没提交,把版本号设置为最新tag-dirty，通过版本号+commitID追溯
    GIT_TREE_STATE='dirty'
    BUILD_VERSION=${BUILD_VERSION}-'dirty'
  } || GIT_TREE_STATE='clean'

  LDFLAGS="-ldflags -X main.Version=${BUILD_VERSION}
                               -X main.buildDate=${BUILD_DATE}
                               -X main.gitCommit=${COMMIT_ID}
                               -X main.gitTreeState=${GIT_TREE_STATE}"

}