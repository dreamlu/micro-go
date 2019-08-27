#!/usr/bin/env bash
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags netgo -o main

./main --registry=consul --api_namespace micro-go.web api --handler web

rm -f main