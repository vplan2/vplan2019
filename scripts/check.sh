#!/bin/bash

# --------------------------------- COLORS ---------------------------------
NC='\033[0m'
BACK="\033[0;30m"
RED="\033[0;31m"   
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
PURPLE="\033[0;35m"
CYAN="\033[0;36m"
WHITE="\033[0;37m"
# --------------------------------------------------------------------------

# -------------------------------- STATICS ---------------------------------
CHECK="✔️ "
CROSS="❌ "
# --------------------------------------------------------------------------

# ------------------------------- FUNCTIONS --------------------------------
function logclr {
    printf "${1}${2}${NC}\n"
}

function check_tool_req {
    which ${1} &> /dev/null && {
        logclr $GREEN "${CHECK} Found '${1}' at '$(which ${1})'"
        return 1
    } || {
        logclr $RED   "${CROSS} Not found '${1}'"
        logclr $WHITE "Download/Installation: ${CYAN}${2}"
    }
}

function check_envvar {
    [ "${2}" == "" ] && {
        logclr $RED   "${CROSS} '${1}' not set"
        logclr $WHITE "${3}"
    } || {
        logclr $GREEN "${CHECK} '${1}' is set to ${2//\\/\/}"
        return 1
    }
}

function check_go_tool {
    which ${1} &> /dev/null && {
        logclr $GREEN "${CHECK} '${1}' is set up"
    } || {
        logclr $RED   "${CROSS} '${1}' not set up. Setting up now..."
        ${2} && {
            logclr $GREEN "  ${CHECK} set up succeed"
        } || {
            logclr $RED   "  ${CROSS} setup failed"
        }
    }
}

# --------------------------------------------------------------------------
# CHECK FOR REQUIRED TOOLS

_failed_req=0

logclr $PURPLE "Checking for required tools..."

check_tool_req "go"   "https://golang.org/dl/" \
    && _failed_req=$((_failed_req + 1))

check_tool_req "dep"  "https://github.com/golang/dep" \
    && _failed_req=$((_failed_req + 1))

check_tool_req "gcc"  "https://gcc.gnu.org/install/" \
    && _failed_req=$((_failed_req + 1))

check_tool_req "make" "https://www.gnu.org/software/make/" \
    && _failed_req=$((_failed_req + 1))

if ((_failed_req > 0)); then
    exit 1
fi

# --------------------------------------------------------------------------
# CHECK FOR ENV VARS

_failed_envvars=0

logclr $PURPLE "\nChecking envoirement variables..."

check_envvar "GOPATH" "${GOPATH}" "Read her how to set up goapth: https://golang.org/cmd/go/#hdr-GOPATH_environment_variable" \
    && _failed_envvars=$((_failed_envvars + 1))

check_envvar "GOROOT" "${GOROOT}" "Set this to the location where your go installation is located" \
    && _failed_envvars=$((_failed_envvars + 1))

if ((_failed_envvars > 0)); then
    exit 1
fi

# --------------------------------------------------------------------------
# CHECK FOR GO TOOLS

logclr $PURPLE "\nChecking for go tools..."

check_go_tool "golint" "go get -u golang.org/x/lint/golint"
check_go_tool "godoc"  "go get -u golang.org/x/tools/cmd/godoc"

logclr $PURPLE "\nChecks completed. Everything set up correctly. 👌"