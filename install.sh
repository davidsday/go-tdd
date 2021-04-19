#! /bin/sh

# cd "${HOME}/.config/nvim/plugged/goTestParser/go/"
# cd "${0%/*}"
cd $(dirname "$(readlink -f "$0")")

mkdir -p ./bin

go build

mv goTestParser ./bin/
