#!/bin/sh

GO_PWD=`pwd`/../../../

echo "Building golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework
cd $GO_PWD/framework/src/golanger
go install .

echo "Building example guestbook"
cd $GO_PWD/example/src/guestbook
rm guestbook
go build .

$GO_PWD/example/src/guestbook/guestbook -addr=:8082
