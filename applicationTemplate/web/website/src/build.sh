#!/bin/sh
APP="website"
ADDR=":8080"
C_PWD=`pwd`/..
GO_PWD=${C_PWD}/../../..
echo "Building framework..."
export GOPATH=${GO_PWD}/add-on:${GO_PWD}/framework:${C_PWD}/src/add-on:${C_PWD}
echo "Building ${APP}"
cd ${C_PWD}/src

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
