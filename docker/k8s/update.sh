#!/usr/bin/env bash
# k8s 更新执行
# 一次定时任务
# apt install at -y
at 02:01 << EOF
kubectl rollout restart user-srv-dep common-srv-dep api-gateway-dep
EOF
# 查看结果
echo -e "更新时间:"
at -l
#exit