#!/bin/bash

set -e

go build -o app

# read from stdin
printf "hello" | ./app

# read from file
printf hello > 1.txt && ./app -r 1.txt

# read from http
./app -r http://httpbin.org/ip

rm app
rm 1.txt
