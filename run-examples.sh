#!/bin/bash

set -e

CWD=`pwd`
EXAMPLES="./examples"

# example basic
APP=basic
cd $EXAMPLES/$APP
go build -o $APP
./$APP
rm $APP
cd $CWD

# examples hello
APP=hello
cd $EXAMPLES/$APP
go build -o $APP
./$APP
rm $APP
cd $CWD

# examples http
APP=http
cd $EXAMPLES/$APP
go build -o $APP
./$APP
rm $APP
cd $CWD

# examples multi-command
APP=multi-command
cd $EXAMPLES/$APP
go build -o $APP
./$APP
rm $APP
cd $CWD

# examples screenshot
APP=screenshot
cd $EXAMPLES/$APP
go build -o $APP
rm $APP
cd $CWD

# examples simplecmd
APP=simple-command
cd $EXAMPLES/$APP
go build -o $APP
rm $APP
cd $CWD

# examples tree
APP=tree
cd $EXAMPLES/$APP
go build -o $APP
rm $APP
cd $CWD
