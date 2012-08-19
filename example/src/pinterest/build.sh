#!/bin/sh

GO_PWD=`pwd`/../../../

echo "export GOPATH for golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework

echo "Building example pinterest"
cd $GO_PWD/example/src/pinterest
rm pinterest
go build .

$GO_PWD/example/src/pinterest/pinterest -addr=:8085
