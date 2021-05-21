#! /bin/sh

# cd "${HOME}/.config/nvim/plugged/go-tdd/go/"
# cd "${0%/*}"
cd $(dirname "$(readlink -f "$0")")

mkdir -p ./bin

go build

mv go-tdd ./bin/
