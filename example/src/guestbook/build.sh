#!/bin/sh

GO_PWD=`pwd`/../../../

echo "export GOPATH for golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework

echo "Building example guestbook"
cd $GO_PWD/example/src/guestbook
rm guestbook
go build .

$GO_PWD/example/src/guestbook/guestbook -addr=:8082
