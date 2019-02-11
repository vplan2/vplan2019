#!/bin/bash

##### CONFIGS ######
OUTPATH=./bin
BINNAME=server
OS=linux
ARCH=amd64
###################

if [ "$1" != "" ]; then
    OS=$1
fi

if [ "$2" != "" ]; then
    ARCH=$2
fi

if [ "$OS" == "windows" ]; then
    BINNAME=${BINNAME}.exe
fi

TAG=$(git describe --tags)
if [ "$TAG" == "" ]; then
    TAG="untagged"
fi

if [ ! -d $BUILDPATH ]; then
    mkdir -p $BUILDPATH
fi

COMMIT=$(git rev-parse HEAD)

echo "Getting dependencies..."
go get -v -t ./...

echo "Building..."
(
    env GOOS=$OS GOARCH=$ARCH \
    go build \
        -v \
        -o ${OUTPATH}/${BINNAME} \
        -ldflags " \
            -X github.com/zekroTJA/shinpuru/internal/util.AppVersion=$TAG \
            -X github.com/zekroTJA/shinpuru/internal/util.AppCommit=$COMMIT \
            -X github.com/zekroTJA/shinpuru/internal/util.Release=TRUE" \
        ./cmd/server
)

wait