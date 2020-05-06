# Rockset CLI
[![CircleCI](https://circleci.com/gh/rockset/cli/tree/go-client-v0.8.svg?style=shield)](https://circleci.com/gh/rockset/cli/tree/go-client-v0.8)
[![Documentation](https://godoc.org/github.com/rockset/cli?status.svg)](http://godoc.org/github.com/rockset/cli)
[![License](https://img.shields.io/github/license/rockset/cli.svg?maxAge=2592000)](https://github.com/rockset/cli/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/rockset/cli.svg)](https://github.com/rockset/cli/issues)
[![Release](https://img.shields.io/github/release/rockset/cli.svg?label=Release)](https://github.com/rockset/cli/releases)

## Building

```
go build -o rock
```

## Testing

```
go test ./...
```

### Integration testing

Requires the environment variable `ROCKSET_APIKEY` to be set

```
go test -tags=integration ./...
```
