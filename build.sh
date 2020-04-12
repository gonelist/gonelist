#!/bin/bash

#脚本要存放在项目根目录
readonly CUR_DIR=$(cd $(dirname ${BASH_SOURCE:-$0}); pwd)

: ${OUTPUT:=${CUR_DIR}/release} ${PROJECT_NAME:=gonelist}



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
    amd64
    386
    arm
    arm64
)
linux=(
    386
    amd64
    arm
    arm64
    ppc64
    ppc64le
    mips
    mipsle
    mips64
    mips64les
    s390x
)
windows=(
    386
    amd64
)
freebsd=(
    386
    amd64
    arm
)
netbsd=(
    386
    amd64
    arm
)
openbsd=(
    386
    amd64
    arm
    arm64
)


FILE_LIST=(
    release/dist
    config.json
)

for file in ${FILE_LIST[@]};do
    ls $file 1>/dev/null
    if [ "$?" -ne 0 ];then
        exit
    fi
done


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
        GOOS=$os GOARCH=$arch go build -o ${save_dir}/${bin_file} main.go 2>/dev/null
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
