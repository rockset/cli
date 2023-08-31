# Rockset CLI
[![CircleCI](https://circleci.com/gh/rockset/cli.svg?style=shield)](https://circleci.com/gh/rockset/cli)
[![Documentation](https://godoc.org/github.com/rockset/rockset-go-cli?status.svg)](http://godoc.org/github.com/rockset/rockset-go-cli)
[![License](https://img.shields.io/github/license/rockset/cli.svg?maxAge=2592000)](https://github.com/rockset/rockset-go-cli/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/rockset/cli.svg)](https://github.com/rockset/rockset-go-cli/issues)
[![Release](https://img.shields.io/github/release/rockset/cli.svg?label=Release)](https://github.com/rockset/rockset-go-cli/releases)

## Usage

![screen recording](vhs/demo.gif)

## Configuration

The Rockset CLI requires having access to an API key and an API server, which can be configured using either
environment variables or a configuration file.

### Environment variables

* `ROCKSET_APIKEY`
* `ROCKSET_APISERVER`

### Configuration File

`~/.config/rockset/cli.yaml`

```yaml
---
current: demo
configs:
  demo:
    apikey: ...
    apiserver: api.usw2a1.rockset.com
  demo2:
    apikey: ...
    apiserver: api.use1a1.rockset.com
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

### Create recordings

We use [vhs](https://github.com/charmbracelet/vhs) to record terminal sessions

```
vhs vhs/demo.tape
```

