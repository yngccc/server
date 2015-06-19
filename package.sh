#!/bin/sh

cd $GOPATH/src/github.com/yngccc/server/

go install

rm -r server
mkdir server
cd server

cp -R $GOPATH/src/github.com/yngccc/server/assets ./
cp -R $GOPATH/src/github.com/yngccc/server/database ./
cp -R $GOPATH/bin/server ./
