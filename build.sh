#!/bin/bash -e

ORG_PATH=github.com/nildev
REPO_PATH="${ORG_PATH}/artemis"
VER=$1
if [ -z "$1" ]; then
    VER=`git rev-parse --abbrev-ref HEAD`
fi

echo "Building artemisd"
CGO_ENABLED=0 godep go build -o bin/artemisd -a -installsuffix netgo -ldflags "-s -X main.Version=$VER -X main.GitHash=`git rev-parse HEAD` -X main.BuiltTimestamp=`date -u '+%Y-%m-%d_%I:%M:%S%p'`" ${REPO_PATH}
chmod +x bin/artemisd