#!/bin/sh

GO_PWD=`pwd`/../../../

echo "Building golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework
cd $GO_PWD/framework/src/golanger
go install .

echo "Building example todo"
cd $GO_PWD/example/src/todo
rm todo
go build .

$GO_PWD/example/src/todo/todo -addr=:8081
