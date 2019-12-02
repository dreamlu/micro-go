# 批量推向私有仓库/这里是公有仓库
#!/bin/bash
docker images | grep registry.gitlab.com/dreamlu/micro-go | awk '{print "docker push "$1":"$2}' | sh

# 删除空镜像
docker images | grep none | awk '{print $3 }'| xargs docker rmi

# 删除停止的容器
#docker rm `docker ps -a|grep Exited|awk '{print $1}'`

# ssh 登录
# 执行更新脚本, 取消ssh命令后面注释
# 进入在线部署的目录, 执行更新脚本, 退出
# [一键进行推送更新部署]
#ssh root@ip #"cd micro-go/docker/k8s;./update.sh;exit"