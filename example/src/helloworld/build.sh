#!/bin/sh

GO_PWD=`pwd`/../../../

echo "export GOPATH for golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework

echo "Building example helloworld"
cd $GO_PWD/example/src/helloworld
rm helloworld
go build .

$GO_PWD/example/src/helloworld/helloworld -addr=:8080
