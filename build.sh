#!/bin/bash

#脚本要存放在项目根目录
readonly GONELIST_ROOT=$(cd $(dirname ${BASH_SOURCE:-$0}); pwd)

: ${OUTPUT:=${CUR_DIR}/release} ${PROJECT_NAME:=gonelist}
git=(git --work-tree "${GONELIST_ROOT}")


# https://golang.org/doc/install/source#environment
OS_LIST=(
    darwin
    linux
    windows
    freebsd
    netbsd
    openbsd
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
    mips
    mips64
    mpis64le
    mipsle
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


FILE_LIST=(
    ${GONELIST_ROOT}/release/dist
    ${GONELIST_ROOT}/config.json
)

for file in ${FILE_LIST[@]};do
    ls $file 1>/dev/null
    if [ "$?" -ne 0 ];then
        exit
    fi
done



TAG_COMMITID=$("${git[@]}" rev-list --tags --max-count=1)
TAG_NUM=$("${git[@]}" describe --tags $(TAG_COMMITID))
rm -f gonelist main

git checkout $TAG_NUM
BUILD_VERSION=$TAG_NUM


HEAD_COMMIT=$("${git[@]}" rev-parse HEAD)
BUILD_DATE=$(date +'%Y-%m-%dT%H:%M:%S')
GIT_TREE_STATE=$("${git[@]}" status --porcelain 2>/dev/null)

[ -z "{$GIT_TREE_STATE}" ] && {
  GIT_TREE_STATE='clean'
  COMMIT_ID=$HEAD_COMMIT
} || {
  GIT_TREE_STATE='dirty'
	BUILD_VERSION=${BUILD_VERSION}-dirty
}


[ "${HEAD_COMMIT}" != "${TAG_COMMITID}" ] && {
  BUILD_VERSION=${BUILD_VERSION}-dirty
  COMMIT_ID=$HEAD_COMMIT
}


LDFLAGS="-ldflags 'main.Version=${BUILD_VERSION}
                             -X main.buildDate=${BUILD_DATE}
                             -X main.gitCommit=${COMMIT_ID}
                             -X main.gitTreeState=${GIT_TREE_STATE}'"

mkdir -p ${GONELIST_ROOT}/release/
[ ! -d "${GONELIST_ROOT}/release/dist" ] && {
  cd ${GONELIST_ROOT}/release/
  curl -sL https://github.com/Sillywa/gonelist-web/releases/download/${TAG_NUM}/dist.tar.gz | tar -zxf -
  cd $GONELIST_ROOT
}

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
        cd $CUR_DIR
        GOOS=$os GOARCH=$arch go build -o ${save_dir}/${bin_file} ${LDFLAGS} \
            ${GONELIST_ROOT}/main.go 2>/dev/null
        if [ "$?" -ne 0 ];then
            echo -e "\t\033[1;31m[failed]\033[0m"
            rm -rf $save_dir
            continue
        fi

        cp -a ${FILE_LIST[@]} ${save_dir}/
        # 进入目录打包
        cd ${OUTPUT}
        tar -zcf $dir_name.tar.gz ${dir_name}/
        cd $CUR_DIR
        rm -rf $save_dir
        echo -e "\t\033[1;32m[success]\033[0m"
    done
done