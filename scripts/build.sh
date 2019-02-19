#!/bin/bash

##### CONFIGS ######
OUTPATH=./bin
BINNAME=vplan2019_server
OS=linux
ARCH=amd64
WORKTREE=github.com/zekroTJA/vplan2019
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

GOVERS=$(go version | sed -e 's/ /_/g')

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
            -X ${WORKTREE}/internal/ldflags.AppVersion=$TAG \
            -X ${WORKTREE}/internal/ldflags.AppCommit=$COMMIT \
            -X ${WORKTREE}/internal/ldflags.GoVersion=$GOVERS \
            -X ${WORKTREE}/internal/ldflags.Release=TRUE" \
        ./cmd/server
)

wait