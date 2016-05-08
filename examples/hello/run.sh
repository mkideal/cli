#/bin/bash

go build

./hello
./hello -h
./hello --name Clipher
./hello --name Clipher -a 12
./hello --name Hiha --age 12
./hello --name Hiha --age=256

rm hello
