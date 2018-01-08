# Mzcoin

[![GoDoc](https://godoc.org/github.com/ShanghaiKuaibei/mzcoin?status.svg)](https://godoc.org/github.com/ShanghaiKuaibei/mzcoin)
[![Go Report Card](https://goreportcard.com/badge/github.com/ShanghaiKuaibei/mzcoin)](https://goreportcard.com/report/github.com/ShanghaiKuaibei/mzcoin)

Mzcoin is a next-generation cryptocurrency.

Mzcoin improves on Bitcoin in too many ways to be addressed here.

Mzcoin is a small part of OP Redecentralize and OP Darknet Plan.

## Table of Contents

<!-- MarkdownTOC depth="2" autolink="true" bracket="round" -->

- [Installation](#installation)
    - [Go 1.9+ Installation and Setup](#go-19-installation-and-setup)
    - [Go get mzcoin](#go-get-mzcoin)
    - [Run Mzcoin from the command line](#run-mzcoin-from-the-command-line)
    - [Show Mzcoin node options](#show-mzcoin-node-options)
    - [Run Mzcoin with options](#run-mzcoin-with-options)
    - [Docker image](#docker-image)
- [API Documentation](#api-documentation)
    - [Wallet REST API](#wallet-rest-api)
    - [JSON-RPC 2.0 API](#json-rpc-20-api)
    - [Mzcoin command line interface](#mzcoin-command-line-interface)
- [Development](#development)
    - [Modules](#modules)
    - [Running Tests](#running-tests)
    - [Formatting](#formatting)
    - [Code Linting](#code-linting)
    - [Dependency Management](#dependency-management)
    - [Wallet GUI Development](#wallet-gui-development)

<!-- /MarkdownTOC -->

## Installation

### Go 1.9+ Installation and Setup

[Golang 1.9+ Installation/Setup](./Installation.md)

### Go get mzcoin

```sh
go get github.com/ShanghaiKuaibei/mzcoin/...
```

This will download `github.com/ShanghaiKuaibei/mzcoin` to `$GOPATH/src/github.com/ShanghaiKuaibei/mzcoin`.

You can also clone the repo directly with `git clone https://github.com/ShanghaiKuaibei/mzcoin`,
but it must be cloned to this path: `$GOPATH/src/github.com/ShanghaiKuaibei/mzcoin`.

### Run Mzcoin from the command line

```sh
cd $GOPATH/src/github.com/ShanghaiKuaibei/mzcoin
make run
```

### Show Mzcoin node options

```sh
cd $GOPATH/src/github.com/ShanghaiKuaibei/mzcoin
make run-help
```

### Run Mzcoin with options

```sh
cd $GOPATH/src/github.com/ShanghaiKuaibei/mzcoin
make ARGS="--launch-browser=false" run
```

### Docker image

A Dockerfile is available at https://github.com/ShanghaiKuaibei/docker-img

## API Documentation

### Wallet REST API

[Wallet REST API](src/gui/README.md).

### JSON-RPC 2.0 API

[JSON-RPC 2.0 README](src/api/webrpc/README.md).

### Mzcoin command line interface

[CLI command API](cmd/cli/README.md).

### Modules

* `/src/cipher` - cryptography library
* `/src/coin` - the blockchain
* `/src/daemon` - networking and wire protocol
* `/src/visor` - the top level, client
* `/src/gui` - the web wallet and json client interface
* `/src/wallet` - the private key storage library
* `/src/api/webrpc` - JSON-RPC 2.0 API
* `/src/api/cli` - CLI library

### Running Tests

```sh
make test
```

### Formatting

All `.go` source files should be formatted `goimports`.  You can do this with:

```sh
make format
```

### Code Linting

Install prerequisites:

```sh
make install-linters
```

Run linters:

```sh
make lint
```

### Dependency Management

Dependencies are managed with [dep](https://github.com/golang/dep).

To install `dep`:

```sh
go get -u github.com/golang/dep
```

`dep` vendors all dependencies into the repo.

If you change the dependencies, you should update them as needed with `dep ensure`.

Use `dep help` for instructions on vendoring a specific version of a dependency, or updating them.

After adding a new dependency (with `dep ensure`), run `dep prune` to remove any unnecessary subpackages from the dependency.

When updating or initializing, `dep` will find the latest version of a dependency that will compile.

Examples:

Initialize all dependencies:

```sh
dep init
dep prune
```

Update all dependencies:

```sh
dep ensure -update -v
dep prune
```

Add a single dependency (latest version):

```sh
dep ensure github.com/foo/bar
dep prune
```

Add a single dependency (more specific version), or downgrade an existing dependency:

```sh
dep ensure github.com/foo/bar@tag
dep prune
```

### Wallet GUI Development

The compiled wallet source should be checked in to the repo, so that others do not need to install node to run the software.

Instructions for doing this:

[Wallet GUI Development README](src/gui/static/README.md)

#### Creating release builds

[Create Release builds](electron/README.md).

## Changelog

[CHANGELOG.md](CHANGELOG.md)
