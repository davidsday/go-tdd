#! /bin/sh

#cd $(dirname "$(readlink -f "$0")")
CURDIR="$(pwd)"

mkdir -p ${CURDIR}/bin

go build -o ${CURDIR}/bin/

