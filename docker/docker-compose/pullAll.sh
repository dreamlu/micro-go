# 批量拉取镜像
#!/bin/bash
cat docker-compose.yaml  | grep common | awk '{print "docker pull "$2}' | sh