#!/bin/bash

set -e

go build -o app

cat >app.json<<EOF
{
	"host": "127.0.0.1",
	"port": 8080
}
EOF

echo ===run root command
./app -h
echo

echo =case 1: read from file
./app
echo

echo =case 2: read from ENV variables
HOST=10.0.0.1 ./app
echo

echo =case 3: read from command line
HOST=10.0.0.1 ./app -H 168.0.0.1
echo

echo
echo ===run sub command
./app sub -h
echo

echo =case 1: read from file
./app sub -w xxx
echo

echo =case 2: read from ENV variables
HOST=10.0.0.1 ./app sub -w xxx
echo

echo =case 3: read from command line
HOST=10.0.0.1 ./app sub -H 168.0.0.1 -w xxx
echo

echo
rm app.json
rm app
