#!/bin/bash

go build -o app

./app daemon --echo "Hello, daemon"

rm app
