#!/bin/sh

GO_PWD=`pwd`/../../../

echo "Building golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework
cd $GO_PWD/framework/src/golanger
go install .

echo "Building example website-admin"
cd $GO_PWD/example/src/website-admin
rm website-admin
go build .

$GO_PWD/example/src/website-admin/website-admin -addr=:8084
