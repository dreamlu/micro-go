#!/usr/bin/env bash

# -tags netgo apline构建golang编译问题

# go mod 中的静态资源引入问题
#GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main
export CGO_ENABLED=0
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o main

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
docker build -f ./Dockerfile -t registry.cn-hangzhou.aliyuncs.com/dreamlu/common:"${version}" .
docker tag registry.cn-hangzhou.aliyuncs.com/dreamlu/common:"${version}" registry.cn-hangzhou.aliyuncs.com/dreamlu/common:latest

# remove build
rm -rf main