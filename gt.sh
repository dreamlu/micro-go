#!/usr/bin/env bash
./devMode.sh prod

cd docker
./pushAll.sh
cd ..
./devMode.sh dev