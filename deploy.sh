#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o dist/stellar-core-export -v
scp ./dist/stellar-core-export $USER@astrograph.evilmartians.io:/home/gzigzigzeo/stellar-core-export