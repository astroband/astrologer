#!/bin/bash

echo "Deploying Astrograph-chronicler..."
echo "Building..."
GOOS=linux GOARCH=amd64 go build -o dist/stellar-core-export -v

echo "Uploading..."
scp ./dist/stellar-core-export $USER@astrograph.evilmartians.io:/home/gzigzigzeo/stellar-core-export

echo "Restarting..."
# cp /home/gzigzigzeo/stellar-core-export /home/deploy/
# ssh $USER@astrograph.evilmartians.io