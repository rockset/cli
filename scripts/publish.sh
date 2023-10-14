#!/bin/sh

set -e

# TODO replace with goreleaser

BINARY="rockset"
VERSION_FILE="version.json"

go vet ./...

# make sure we have working credentials
aws sts get-caller-identity --no-cli-pager > /dev/null

BUILD=""
if [ -z "${SENTRY_DSN}" ]; then
  BUILD="-ldflags '-X main.dsn=${SENTRY_DSN}'"
fi

echo "compiling binaries..."
for OS in Darwin linux; do
  for ARCH in arm64 amd64; do
    # GOOS is all lowercase while macOS will report Darwin from uname -s
    GOOS=$(echo ${OS} | tr '[:upper:]' '[:lower:]') GOARCH=${ARCH} go build ${BUILD} -o "${BINARY}"
    aws s3 cp "${BINARY}" "s3://rockset.sh/install/${OS}/${ARCH}/${BINARY}"
  done
done

rm "${BINARY}"

VERSION="$(grep -ohE 'v\d+\.\d+\.\d+' version.go)"
cat > "${VERSION_FILE}" <<EOT
{
  "stable": "${VERSION}"
}
EOT

aws s3 cp web/index.html "s3://rockset.sh/index.html"
aws s3 cp web/install.sh "s3://rockset.sh/install"

aws s3 cp version.json "s3://rockset.sh/install/${VERSION_FILE}"
rm "${VERSION_FILE}"
