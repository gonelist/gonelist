#!/usr/bin/env bash

#GONELIST_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
readonly git=(git --work-tree "${GONELIST_ROOT}")
readonly BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
readonly HEAD=$("${git[@]}" rev-parse "HEAD^{commit}")
export GONELIST_ROOT BUILD_DATE HEAD


GONELIST::getVersion(){
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


  "${git[@]}" checkout $TAG 2>/dev/null
  BUILD_VERSION=${TAG}

  if [ -z "$("${git[@]}" status --porcelain 2>/dev/null)" ];then
    GIT_TREE_STATE='clean'
  else
    GIT_TREE_STATE='dirty'
  fi

  #master切到tag版本则置为dirty

  if [ "${HEAD}" != "${TAG_COMMITID}" ];then
    #tag的基础上改动，所以tag版本号-dirty
    BUILD_VERSION+="-dirty"
    COMMIT_ID=${HEAD}
  else
    COMMIT_ID=${TAG_COMMITID}
  fi

  "${git[@]}" checkout $HEAD 2>/dev/null
}

GONELIST::SetVersion(){
  GONELIST::getVersion &>/dev/null

  local -a ldflags
  function add_ldflag() {
    local key=${1}
    local val=${2}
    ldflags+=(
      "-X 'main.${key}=${val}'"
    )
  }

  add_ldflag 'Version' ${BUILD_VERSION}
  add_ldflag 'buildDate' ${BUILD_DATE}
  add_ldflag 'gitCommit' ${COMMIT_ID}
  add_ldflag 'gitTreeState' ${GIT_TREE_STATE}

  unset TAG_COMMITID BUILD_VERSION COMMIT_ID GIT_TREE_STATE

  # The -ldflags parameter takes a single string, so join the output.
  echo $TAG "${ldflags[*]-}"

}