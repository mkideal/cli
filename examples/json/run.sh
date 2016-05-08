#!/bin/bash

go build -o app

./app -h
./app --json='{"Int": 12, "String": "Hello"}'

rm app
