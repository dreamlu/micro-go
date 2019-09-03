#!/usr/bin/env bash

# 统一docker 构建

# 工作目录
workDir=`pwd`
echo "工作目录: ${workDir}"
# execute cmd
modules=(api-gateway common-srv user-srv)
for module in ${modules[*]}
#也可以写成for element in ${array[@]}
do
echo -e "\n模块: ""${module} 开始构建docker"
cd ${workDir}/${module} && ./docker.sh
done