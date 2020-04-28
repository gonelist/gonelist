#!/usr/bin/env bash

#GONELIST_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
readonly git=(git --work-tree "${GONELIST_ROOT}")
readonly BUILD_DATE=$(date +'%Y-%m-%dT%H:%M:%S')
readonly HEAD=$("${git[@]}" rev-parse "HEAD^{commit}")

GONELIST::SetVersion(){

  local -a ldflags
  function add_ldflag() {
    local key=${1}
    local val=${2}
    # If you update these, also update the list component-base/version/def.bzl.
    ldflags+=(
      "-X main.${key}=${val}"
    )
  }

  #指定tag版本时候判断存在否
  if [ -n "$TAG" ]; then
    TAG_COMMITID=$("${git[@]}" rev-parse $TAG 2>/dev/null)
    if [ "$?" -ne 0 ];then
        echo no such tag: $TAG
        exit 1
    fi
  else #默认取最新的tag
    TAG_COMMITID=$("${git[@]}" rev-list --tags --max-count=1)
    TAG=$("${git[@]}" describe --tags ${TAG_COMMITID})
  fi


  "${git[@]}" checkout $TAG 1>/dev/null
  BUILD_VERSION=${TAG}

  GIT_TREE_STATE=$("${git[@]}" status --porcelain 2>/dev/null)

  if [ -z "${GIT_TREE_STATE}" ];then
    GIT_TREE_STATE='clean'
  else
    GIT_TREE_STATE='dirty'
  fi

  #在tag的版本上更改了代码则置为dirty
  HEAD_COMMIT=$("${git[@]}" rev-parse HEAD)
  if [ "${HEAD_COMMIT}" != "${TAG_COMMITID}" ];then
    #tag的基础上改动，所以tag版本号-dirty
    BUILD_VERSION=${BUILD_VERSION}-dirty
    COMMIT_ID=${HEAD_COMMIT}
  else
    COMMIT_ID=${TAG_COMMITID}
  fi

  add_ldflag 'Version' ${BUILD_VERSION}
  add_ldflag 'buildDate' ${BUILD_DATE}
  add_ldflag 'gitCommit' ${COMMIT_ID}
  add_ldflag 'gitTreeState' ${GIT_TREE_STATE}

  # The -ldflags parameter takes a single string, so join the output.
  echo $TAG_NUM -ldflags \'${ldflags[*]-}\'
}