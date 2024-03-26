#! /bin/sh

cd $(dirname "$(readlink -f "$0")")

CURDIR=$(pwd)

mkdir -p ./bin

go build -o ${CURDIR}/bin/

