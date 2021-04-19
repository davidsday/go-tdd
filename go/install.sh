#! /bin/sh

cd "${HOME}/.config/nvim/plugged/goTestParser/go/"
mkdir -p ./bin

go build

mv goTestParser ./bin/
