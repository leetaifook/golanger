#!/bin/sh

GO_PWD=`pwd`/../../../

echo "Building golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework
cd $GO_PWD/framework/src/golanger
go install .

echo "Building example todo-orm"
cd $GO_PWD/example/src/todo-orm
rm todo-orm
go build .

$GO_PWD/example/src/todo-orm/todo-orm -addr=:8083
