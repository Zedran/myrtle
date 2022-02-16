#!/usr/bin/env bash

src=./src
out=./build/myrtle

go test $src -test.v

if [ $? == 0 ]
then
    go build -o $out -trimpath -ldflags "-s -w" $src
fi
