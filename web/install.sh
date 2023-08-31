#!/bin/sh

set -e

OS=$(uname -s)
ARCH=$(uname -m)

# TODO enable CloudFront and change to https once we are done testing
URL="http://rockset.sh/install/${OS}/${ARCH}/rockset"
BINARY="rockset"

curl -o "${BINARY}" "${URL}"
chmod 755 "${BINARY}"

echo "Rockset CLI is installed as: ${BINARY}"
