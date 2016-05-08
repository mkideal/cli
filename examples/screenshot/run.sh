#!/bin/bash

go build -o app

./app -h
./app

rm app
