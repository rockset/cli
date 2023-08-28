# Rockset CLI
[![CircleCI](https://circleci.com/gh/rockset/cli.svg?style=shield)](https://circleci.com/gh/rockset/cli)
[![Documentation](https://godoc.org/github.com/rockset/rockset-go-cli?status.svg)](http://godoc.org/github.com/rockset/rockset-go-cli)
[![License](https://img.shields.io/github/license/rockset/cli.svg?maxAge=2592000)](https://github.com/rockset/rockset-go-cli/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/rockset/cli.svg)](https://github.com/rockset/rockset-go-cli/issues)
[![Release](https://img.shields.io/github/release/rockset/cli.svg?label=Release)](https://github.com/rockset/rockset-go-cli/releases)

## Usage

### Configuration

Environment variables

* `ROCKSET_APIKEY`
* `ROCKSET_APISERVER`

Configuration file
`~/.config/rockset/cli.yaml`

```yaml
---
current: demo
configs:
  demo:
    apikey: ...
    apiserver: api.usw2a1.rockset.com
```

## Building

```
go build -o rockset
```

## Testing

```
go test ./...
```

### Integration testing

Requires the environment variable `ROCKSET_APIKEY` to be set

```
go test ./...
```
