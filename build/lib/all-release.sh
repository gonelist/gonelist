#!/bin/bash

: ${OUTPUT:=${GONELIST_ROOT}/release} ${PROJECT_NAME:=gonelist}


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


mkdir -p ${GONELIST_ROOT}/release/
[ ! -d "${GONELIST_ROOT}/release/dist" ] && {
  cd ${GONELIST_ROOT}/release/
  curl -sL https://github.com/Sillywa/gonelist-web/releases/download/${TAG_NUM}/dist.tar.gz | tar -zxf -
  cd $GONELIST_ROOT
}

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
        bash -c "GOOS=$os GOARCH=$arch go build -o ${save_dir}/${bin_file}  ${LDFLAGS} ${GONELIST_ROOT}/main.go 2>/dev/null"
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