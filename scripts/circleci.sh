#!/bin/bash

git submodule init
git submodule update

go build -v -o ./bin/server ./cmd/server