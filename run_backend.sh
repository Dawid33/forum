#!/bin/sh
cd main
go build -o /tmp/backend .
cd ..
/tmp/backend