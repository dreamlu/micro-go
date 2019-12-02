#!/usr/bin/env bash

# 统一docker 构建

# 工作目录
workDir=`pwd`
# docker build
# default version
version=0.1
# 参数处理
# :需要参数
while getopts ":v:h" opt
do
    case ${opt} in
        v)
        version=$OPTARG
        echo "版本号version的值${version}"
        ;;
        h)
        echo -e "-v 版本号id\n-h 帮助\n"
        exit 1
        ;;
        ?)
        echo "未知参数"
        exit 1;;
    esac
done
echo "工作目录: ${workDir}"
# execute cmd
modules=(api-gateway common-srv user-srv)
for module in ${modules[*]}
#也可以写成for element in ${array[@]}
do
echo -e "\n模块: ""${module} 开始构建docker"
cd ${workDir}/${module} && ./docker.sh -v ${version}
done