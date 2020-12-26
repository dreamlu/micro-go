#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o main

./main --server_name api-gateway api --address=0.0.0.0:9000 --namespace=demo.api --handler=http

rm -f main