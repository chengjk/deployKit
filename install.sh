#!/usr/bin/env bash

if [ ! -f "install.sh" ]; then
echo 'install must be run within its container folder' 1>&2
exit 1
fi

CUR_DIR=`pwd`
OLD_GOPATH="$GOPATH"
OLD_GOBIN="$GOBIN"

echo "set new go env"
export GOPATH="$CUR_DIR"
export GOBIN="$CUR_DIR"/bin

echo "format code"
gofmt -w src

echo "install"
go install src/ndp/main/dk.go

echo "rollback go env"
export GOPATH="$OLD_GOPATH"
export GOBIN="$OLD_GOBIN"

echo 'finished'