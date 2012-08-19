#!/bin/sh

GO_PWD=`pwd`/../../../

echo "Building golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework
cd $GO_PWD/framework/src/golanger
go install .

echo "Building example pinterest"
cd $GO_PWD/example/src/pinterest
rm pinterest
go build .

$GO_PWD/example/src/pinterest/pinterest -addr=:8085
