#!/bin/sh
set -e

apt-get update
apt-get install -y curl libpq-dev
apt-get clean

echo "\nDone installing stellar-core dependencies...\n"
