Mzcoin [![GoDoc](https://godoc.org/github.com/ShanghaiKuaibei/mzcoin?status.svg)](https://godoc.org/github.com/ShanghaiKuaibei/mzcoin) [![Go Report Card](https://goreportcard.com/badge/github.com/ShanghaiKuaibei/mzcoin)](https://goreportcard.com/report/github.com/ShanghaiKuaibei/mzcoin) 

Mzcoin is children's education coin created by mzworld.

## Installation

For detailed installation instructions, see [Installing mzcoin](../../wiki/Installation)

## For OSX

Install [homebrew](brew.sh), if you don't have it yet.

```sh
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

Install the latest version of golang

```sh
brew install go
```

Setup $GOPATH variable, add it to ~/.bash_profile (or bashrc). After editing, open a new tab
Add to `bashrc` or `bash_profile`

```sh
export GOPATH=/Users/<username>/go
export PATH=$PATH:$GOPATH/bin

```

Install Mercurial and Bazaar

```sh
brew install mercurial bzr
```

Fetch the latest code of mzcoin from the github repository

```sh
go get github.com/ShanghaiKuaibei/mzcoin
```

Change your current directory to $GOPATH/src/github.com/ShanghaiKuaibei/mzcoin

```sh
cd $GOPATH/src/github.com/ShanghaiKuaibei/mzcoin
```

Run Wallet

```sh
./run.sh

OR
go run ./cmd/mzcoin/mzcoin.go

For Options
go run ./cmd/mzcoin/mzcoin.go --help
```

## For linux

```sh
sudo apt-get install curl git mercurial make binutils gcc bzr bison libgmp3-dev screen -y
```

## Setup Golang

use gvm or download binary and follow instructions.

### Golang ENV setup with gvm

In China, use `--source=https://github.com/golang/go` to bypass firewall when fetching golang source.

```sh
sudo apt-get install bison curl git mercurial make binutils bison gcc build-essential
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source $HOME/.gvm/scripts/gvm

gvm install go1.4 --source=https://github.com/golang/go
gvm use go1.4
gvm install go1.8
gvm use go1.8 --default
```

If you open up new terminal and the go command is not found then add this to .bashrc . GVM should add this automatically.

```sh
[[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
gvm use go1.8 >/dev/null
```

The mzcoin repo must be in $GOPATH, under "/src/github.com/ShanghaiKuaibei". Otherwise golang programs cannot import the libraries.

Pull mzcoin repo into the gopath, note: puts the mzcoin folder in $GOPATH/src/github.com/ShanghaiKuaibei/mzcoin

```sh
go get -v github.com/ShanghaiKuaibei/mzcoin/...

# create symlink of the repo
cd $HOME
ln -s $GOPATH/src/github.com/ShanghaiKuaibei/mzcoin mzcoin


## Dependencies

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

## Run A Mzcoin Node

```sh
cd mzcoin
screen
go run ./cmd/mzcoin/mzcoin.go
```

then ctrl+A then D to exit screen
screen -x to reattach screen

### Todo

Use gvm package set, so repo does not need to be symlinked. Does this have a default option?

```sh
gvm pkgset create mzcoin
gvm pkgset use mzcoin
git clone https://github.com/ShanghaiKuaibei/mzcoin
cd mzcoin
go install
```

### Cross Compilation

Install Gox:

```sh
go get github.com/mitchellh/gox
```

Compile:

```sh
gox --help
gox [options] cmd/mzcoin/
```

Modules
-------

## Modules

* /src/cipher - cryptography library
* /src/coin - the blockchain
* /src/daemon - networking and wire protocol
* /src/visor - the top level, client
* /src/gui - the web wallet and json client interface
* /src/wallet - the private key storage library

## Meshnet

```sh
go run ./cmd/mesh/*.go -config=cmd/mesh/sample/config_a.json
go run ./cmd/mesh/*.go -config=cmd/mesh/sample/config_b.json
```

## Meshnet reminders

* one way latency
* two way latency (append), latency between packet and ack
* service handler (ability to append services to meshnet)
* uploading bandwidth, latency measurements over time
* end-to-end route instrumentation

## Rebuilding Wallet HTML

```sh
cd src/gui/static
npm install
gulp build
```

## Release Builds

```sh
cd /src/gui/static
npm install
gulp dist
```

## Skycoin command line interface

See the doc of command line interface [here](cmd/cli/README.md).

## WebRPC

See the doc of webrpc [here](src/api/webrpc/README.md).