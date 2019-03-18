#!/bin/bash

### STATICS ##########################################

NC='\033[0m'
RED="\033[0;31m"   
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
CYAN="\033[0;36m"

WORKTREE="github.com/zekroTJA/vplan2019"
TARGETDIR="./release"

OSES=(
    # "darwin" 
    # "linux" 
    "windows" 
)

ARCHS=(
    "386 amd64" 
    # "386 amd64 arm" 
    # "386 amd64" 
)

### FUNCTIONS ########################################

function logInfo {
    printf "${CYAN}[ INFO ] ${1}${NC}\n"
}

function logError {
    printf "${RED}[ ERROR ] ${1}${NC}\n"
}

function logWarn {
    printf "${YELLOW}[ WARN ] ${1}${NC}\n"
}

function checkCommand {
    which ${1} &> /dev/null || {
        logError "'${1}' command not found"
        exit 1
    }
}

function goBuild {
    BINNAME="vplan2_server_${1}_${2}"
    [ "${1}" == "windows" ] && BINNAME="${BINNAME}.exe"

    (
        env GOOS=${1} GOARCH=${2} \
        go build \
            -v \
            -o ${TARGETDIR}/${BINNAME} \
            -ldflags " \
                -X ${WORKTREE}/internal/ldflags.AppVersion=$TAG \
                -X ${WORKTREE}/internal/ldflags.AppCommit=$COMMIT \
                -X ${WORKTREE}/internal/ldflags.GoVersion=$GOVERS \
                -X ${WORKTREE}/internal/ldflags.Release=TRUE" \
            ./cmd/server
    )
}

### COMMAND CHECKS AND DIRS SETUP ####################

checkCommand git
checkCommand zola
checkCommand dep
checkCommand go
checkCommand tar
checkCommand sed
checkCommand sha256sum

[ -d ${TARGETDIR} ] && rm -r -f ${TARGETDIR}
mkdir -p ${TARGETDIR}

### SETTINGS VARS ####################################

if [ "${GOPATH}" == "" ]; then
    logError 'GOPATH is not set'
    exit 1
fi

TAG=$(git describe --tags)
if [ "$TAG" == "" ]; then
    TAG="untagged"
fi

GOVERS=$(go version | sed -e 's/ /_/g')
COMMIT=$(git rev-parse HEAD)

### GETTING DEPENDENCIES #############################

logInfo "Getting dependencies..."
dep ensure -v

### BUILDING FRONTEND ################################

logInfo "Building frontend dependencies..."
cp -f ./config/frontend.release.toml ./web/config.toml
cd ./web
zola build
cd ..
mv ./web/public ${TARGETDIR}/web

### BUILDING BINARIES ################################

logInfo "Building server binaries..."
for i in ${!OSES[*]}; do
    _archs=${ARCHS[i]}
    for arch in ${_archs[*]}; do
        goBuild ${OSES[i]} $arch
    done
done

### SUMMING UP, PACKING AND COMPRESSING ##############

logInfo "Creating checksums..."
cd ${TARGETDIR}
for file in ./*; do
    [ -f ${file} ] && {
        sha256sum ${file} >> sums.txt
    }
done

logInfo "Creating compressed tarball..."
tar -czvf release_${TAG}.tar.gz ./
sha256sum release_${TAG}.tar.gz >> sum.txt

logInfo "Cleaning up..."
for file in ./*; do
    [ "${file}" != "./release_${TAG}.tar.gz" ] && \
    [ "${file}" != "./sum.txt" ] && \
        rm -r -f ${file}
done

cd ..

logInfo "Finished."