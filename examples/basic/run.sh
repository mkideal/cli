#!/bin/bash

go build
./basic
./basic --required=1
./basic --required=1 -s=not-a-bool
./basic --required=2 --long-flag=4
./basic --required=2 --long-flag=-4
./basic --required=2 --long-flag -4
./basic --required=2 --long-flag=1001
./basic --required=2 --slice 1 --slice 12 -D8 -Mx=1 -My=2 -M xx=3 -M yy=4 --map zz=33
rm basic
