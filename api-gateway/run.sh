#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o main

./main --api_address=0.0.0.0:9000 --registry=consul --api_namespace micro-go.api api --handler web

rm -f main