#! /bin/sh

#cd $(dirname "$(readlink -f "$0")")

mkdir -p ./bin

go build -o ./bin/

