#!/bin/sh

GO_PWD=`pwd`/../../../

echo "Building golanger framework..."
export GOPATH=$GOPATH:$GO_PWD/add-on:$GO_PWD/framework
cd $GO_PWD/framework/src/golanger
go install .

echo "Building example chatroom"
cd $GO_PWD/example/src/chatroom
rm chatroom
go build .

$GO_PWD/example/src/chatroom/chatroom -addr=:8086
