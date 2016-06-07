#!/bin/bash

MY_GOPATH=`pwd`
MY_GOPATH=${MY_GOPATH%src*}
export GOPATH=$MY_GOPATH
go run cmd/sche_app.go --master="10.8.12.174:25050" --address="10.8.12.174"
#go get cmd/sche_app.go --master="10.8.12.174:25050" --address="10.8.12.174" --docker-image="hello-world"
