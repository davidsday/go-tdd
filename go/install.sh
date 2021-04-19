#! /bin/sh

# cd "${HOME}/.config/nvim/plugged/goTestParser/go/"
cd "${0%/*}"

mkdir -p ./bin

go build

mv goTestParser ./bin/
