#!/bin/bash

PKGNAME="github.com/zekroTJA/vplan2019"
BUILDPATH="./bin"
BUILDNAME="vplan2019_server"

TAG=$(git describe --tags)
if [ "$TAG" == "" ]; then
    TAG="untagged"
fi

COMMIT=$(git rev-parse HEAD)

git submodule init
git submodule update

BUILDS=( \
    'linux;arm' \
    # 'linux;amd64' \  <- actually no needs to be tested atm
    'windows;amd64' \
    # 'darwin;amd64' \ <- actually no needs to be tested atm
)

for BUILD in ${BUILDS[*]}; do

    IFS=';' read -ra SPLIT <<< "$BUILD"
    OS=${SPLIT[0]}
    ARCH=${SPLIT[1]}

    echo "Building ${OS}_$ARCH..."
    (env GOOS=$OS GOARCH=$ARCH \
        go build \
            -o ${BUILDPATH}/${BUILDNAME}_${OS}_$ARCH \
            -ldflags " \
                -X ${PKGNAME}/ldflags.AppVersion=$TAG \
                -X ${PKGNAME}/ldflags.AppCommit=$COMMIT \
                -X ${PKGNAME}/ldflags.Release=TRUE" \
                ./cmd/server
    )
            

    if [ "$OS" = "windows" ]; then
        mv ${BUILDPATH}/${BUILDNAME}_windows_$ARCH $BUILDPATH/${BUILDNAME}_windows_${ARCH}.exe
    fi

done