#!/bin/bash

git submodule init
git submodule update
go build ./cmd/server