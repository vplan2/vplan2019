#!/bin/bash

START_CMD=$1
CFG=$2
CFG_DEF=$3

if ! [ -f $CFG ]; then
    cp $CFG_DEF $CFG
    printf "\n\nDEFAULT CONFIG FILE WAS CREATED AT $CFG (container location).\n"
    printf "EDIT THIS CONFIG FILE AND RESTART.\n\n"
    exit
fi

$START_CMD
