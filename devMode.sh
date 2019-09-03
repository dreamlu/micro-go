#!/usr/bin/env bash

# 开发模式 dev/prod
# 此处修改模式
# 执行该脚本
devMode=prod

# 配置文件地址
# 修改各个模块下app.yaml文件开发模式
confFiles=(user-srv/conf/app.yaml common-srv/conf/app.yaml)

# conf配置修改
for conf in ${confFiles[*]}
#也可以写成for element in ${array[@]}
do
echo "配置文件: ""${conf}"
# 替换源文件第三行内容
# 行首添加两个空格
sed -i '3c \  \devMode: '${devMode} "${conf}"
done

# prod 模式自动构建docker
# 可注释, 通过docker.sh执行构建
if [[ ${devMode} = 'prod' ]]; then
    echo "prod 开始构建docker"
    ./docker.sh
fi