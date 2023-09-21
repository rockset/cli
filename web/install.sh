#!/bin/sh

set -e

OS=$(uname -s)
ARCH=$(uname -m)

# TODO enable CloudFront and change to https once we are done testing
BINARY="rockset"
URL="http://rockset.sh/install/${OS}/${ARCH}/${BINARY}"

echo "Downloading Rockset CLI from: ${URL}"

curl -s -o "${BINARY}" "${URL}"
chmod 755 "${BINARY}"

QUICKSTART_URL="https://docs.rockset.com/documentation/docs/quickstart"

echo "Rockset CLI is installed as: ${BINARY}"
echo ""
echo "Quickstart guide: ${QUICKSTART_URL}"

if [ "${OS}" = "Darwin" ]; then
  open "${QUICKSTART_URL}"
fi
