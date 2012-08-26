#!/bin/sh
APP="guestbook"
ADDR=":8081"
PWD=`pwd`/..
GO_PWD=${PWD}/../..
echo "Building framework..."
export GOPATH=${GO_PWD}/add-on:${GO_PWD}/framework:${PWD}/src/add-on:${PWD}
echo "Building ${APP}"
cd ${PWD}/src

if [ -f ${APP} ]; then
    echo "rm ${APP}"
    rm ${APP}
fi

go build .

if [ -f src ]; then     
    echo "mv src to ${APP}"
    mv ./src ${APP}
    echo "run ${APP}"
    ./$APP -addr=${ADDR}
fi