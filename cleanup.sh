#!/bin/sh

find . -type f -name StdErr.txt -exec rm {} \;
find . -type f -name out.json -exec rm {} \;
find . -type f -name out.txt -exec rm {} \;
find . -type f -name out -exec rm {} \;
find . -type f -name go-tdd_log.json -exec rm {} \;
