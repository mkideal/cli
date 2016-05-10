#!/bin/bash

go build -o app

./app -h
./app --json='{"Int": 12, "String": "Hello"}'

echo '{"Int": 222, "String": "Hello"}' > 2.json
./app --json='{"Int": 12, "String": "Hello"}' --jsonfile 2.json
rm 2.json

rm app
