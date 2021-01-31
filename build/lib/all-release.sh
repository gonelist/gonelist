#!/bin/bash
[ -n "$RELEASE_DEBUG" ] && set -x
#脚本要存放在项目根目录

cd ${GONELIST_ROOT}
readonly git=(git --work-tree "${GONELIST_ROOT}")
: ${OUTPUT:=${GONELIST_ROOT}/release} ${PROJECT_NAME:=gonelist}

if [ -z `command -v go` ];then
  echo go is not in PATH
  exit 1
fi


# https://golang.org/doc/install/source#environment
OS_LIST=(
    darwin
    linux
    windows
    #freebsd
    #netbsd
    #openbsd
)
darwin=(
    386
    amd64
)
freebsd=(
    386
    amd64
    arm
)
linux=(
    386
    amd64
    arm
    arm64
    #mips
    #mips64
    #mpis64le
    #mipsle
)
netbsd=(
    386
    amd64
    arm
)
openbsd=(
    386
    amd64
)
windows=(
    386
    amd64
    arm
)

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


mkdir -p ${GONELIST_ROOT}/release/
[ ! -d "${GONELIST_ROOT}/release/dist" ] && {
  cd ${GONELIST_ROOT}/release/
  [ -z "$DIST_VERSION" ] && DIST_VERSION=$TAG
  curl -sL https://github.com/Sillywa/gonelist-web/releases/download/${DIST_VERSION}/dist.tar.gz | tar -zvxf -
  cd $GONELIST_ROOT
}

FILE_LIST=(
    ${GONELIST_ROOT}/release/dist
    ${GONELIST_ROOT}/config.json
)

for file in ${FILE_LIST[@]};do
    ls $file 1>/dev/null
    if [ "$?" -ne 0 ];then
        "${git[@]}" checkout master
        exit
    fi
done

echo go build -o ${GONELIST_ROOT}/gonelist -ldflags " -X main.Version=${BUILD_VERSION}
            -X main.buildDate=${BUILD_DATE}
            -X main.gitCommit=${COMMIT_ID}
            -X main.gitTreeState=${GIT_TREE_STATE}"  ${GONELIST_ROOT}/main.go


for os in ${OS_LIST[@]};do
    arch_array="${os}[@]"             # 间接引用数组
    for arch in "${!arch_array}";do   #
        bin_file=${PROJECT_NAME}_${os}_${arch}
        dir_name=${bin_file} #压缩包的根文件夹名和bin_file一样
        if [ "$os" == "windows" ];then
            bin_file=${bin_file}.exe
        fi
        printf "building %-30s" ${bin_file}
        save_dir=${OUTPUT}/${dir_name}
        mkdir -p $save_dir
        cd $GONELIST_ROOT
        GOOS=$os GOARCH=$arch go build -o ${save_dir}/${bin_file} -ldflags " -X main.Version=${BUILD_VERSION}
            -X main.buildDate=${BUILD_DATE}
            -X main.gitCommit=${COMMIT_ID}
            -X main.gitTreeState=${GIT_TREE_STATE}"  ${GONELIST_ROOT}/main.go 2>/dev/null
        if [ "$?" -ne 0 ];then
            echo -e "\t\033[1;31m[failed]\033[0m"
            rm -rf $save_dir
            continue
        fi

        cp -a ${FILE_LIST[@]} ${save_dir}/
        # 进入目录打包
        cd ${OUTPUT}
        tar -zcf $dir_name.tar.gz ${dir_name}/
        cd $GONELIST_ROOT
        rm -rf $save_dir
        echo -e "\t\033[1;32m[success]\033[0m"
    done
done
