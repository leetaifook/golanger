#!/bin/sh

GO_PWD=`pwd`/../../../

echo "export GOPATH for golanger framework..."
export GOPATH=$GO_PWD/add-on:$GO_PWD/framework

echo "Building example chatroom"
cd $GO_PWD/example/src/chatroom
rm chatroom
go build .

$GO_PWD/example/src/chatroom/chatroom -addr=:8086
