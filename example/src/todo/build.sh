#!/bin/sh

GO_PWD=`pwd`/../../../

echo "export GOPATH for golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework

echo "Building example todo"
cd $GO_PWD/example/src/todo
rm todo
go build .

$GO_PWD/example/src/todo/todo -addr=:8081
