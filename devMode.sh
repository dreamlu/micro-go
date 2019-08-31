#!/usr/bin/env bash

# 开发模式 dev/prod
# 此处修改模式
# 执行该脚本
# dev模式
devMode=dev

# prod模式
#devMode=prod

# 后端配置文件地址
# 修改各个模块下app.conf文件开发模式
confFiles=(user-srv/conf/app.yaml common-srv/conf/app.yaml)

# 后端conf配置修改
for conf in ${confFiles[*]}
#也可以写成for element in ${array[@]}
do
echo "配置文件: ""${conf}"
# 替换源文件第三行内容
# 行首添加两个空格
sed -i '3c \  \devMode: '${devMode} "${conf}"
done