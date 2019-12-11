#!/usr/bin/env bash
docker run --restart=unless-stopped -p 80:80 -p 443:443 --name d-rancher rancher/rancher:stable

#docker rm d-rancher