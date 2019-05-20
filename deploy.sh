#!/bin/bash

echo "Deploying Astrograph-chronicler..."
echo "Building..."
GOOS=linux GOARCH=amd64 go build -o dist/astrologer -v

echo "Uploading..."
scp ./dist/astrologer $USER@astrograph.evilmartians.io:/home/astroband/astrologer

echo "Restarting..."
# cp /home/astroband/astrologer /home/deploy/
# ssh $USER@astrograph.evilmartians.io