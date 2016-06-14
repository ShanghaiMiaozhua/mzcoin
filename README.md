
Mzcoin
=======

Mzcoin is children's education coin. Based upon mzcoin.

Installation
------------

*For detailed installation instructions, see [Installing mzcoin](../../wiki/Installation)*

For linux:
sudo apt-get install curl git mercurial make binutils gcc bzr bison libgmp3-dev screen -y

OSX:
brew install mercurial bzr

```
./run.sh -h
```

*Running Wallet

```
./run.sh
Goto http://127.0.0.1:6402

OR

go run ./cmd/mzcoin/mzcoin.go
```

Golang environment setup with gvm
---

The chinese firewall may block golang installation with gvm

```
sudo apt-get install bison curl git mercurial make binutils bison gcc build-essential
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source $HOME/.gvm/scripts/gvm

gvm install go1.4
gvm use go1.4
gvm install go1.6
gvm use go1.6 --default
```

If you open up new terminal and the go command is not found then add this to .bashrc . GVM should add this automatically

```
[[ -s "$HOME/.gvm/scripts/gvm" ]] && source "$HOME/.gvm/scripts/gvm"
gvm use go1.6 >/dev/null
```

---

The mzcoin repo must be in $GOPATH, under "/src/github.com/mzcoin". Otherwise golang programs cannot import the libraries.

```
#pull mzcoin repo into the gopath
#note: puts the mzcoin folder in $GOPATH/src/github.com/mzcoin/mzcoin
go get github.com/mzcoin/mzcoin

#create symlink of the repo
cd $HOME
ln -s $GOPATH/src/github.com/mzcoin/mzcoin mzcoin
```

Dependencies
---

```
go get github.com/robfig/glock
glock sync github.com/mzcoin/mzcoin
go get ./cmd/mzcoin
```

To update dependencies
```
glock save github.com/mzcoin/mzcoin/cmd/mzcoin
```

Running Node
---

```
cd mzcoin
screen
go run ./cmd/mzcoin/mzcoin.go 
#then ctrl+A then D to exit screen
#screen -x to reattach screen
```

Todo
---

Use gvm package set, so repo does not need to be symlinked. Does this have a default option?

```
gvm pkgset create mzcoin
gvm pkgset use mzcoin
git clone https://github.com/mzcoin/mzcoin
cd mzcoin
go install
```

Running
---

cd mzcoin
go run ./cmd/mzcoin/mzcoin.go

Cross Compilation
---

Install Gox:
```
go get github.com/mitchellh/gox
```

Compile:
```
cd compile
./build-dist-all.sh
```

Local Server API
----

Run the mzcoin client then
```
http://127.0.0.1:6420/wallets
http://127.0.0.1:6420/outputs
http://127.0.0.1:6420/blockchain/blocks?start=0&end=500
http://127.0.0.1:6420/blockchain
http://127.0.0.1:6420/connections
```

```
http://127.0.0.1:6420/wallets to get your wallet seed. Write this down

http://127.0.0.1:6420/outputs to see outputs (address balances)

http://127.0.0.1:6420/blockchain/blocks?start=0&end=5000 to see all blocks and transactions.

http://127.0.0.1:6420/connections to check network connections

http://127.0.0.1:6420/blockchain to check blockchain head
```

Public API
----

This is a public server. You can use these urls on local host too, with the mzcoin client running.
```
http://mzcoin-chompyz.c9.io/outputs
http://mzcoin-chompyz.c9.io/blockchain/blocks?start=0&end=500
http://mzcoin-chompyz.c9.io/blockchain
http://mzcoin-chompyz.c9.io/connections
```

Modules
-----

```
/src/cipher - cryptography library
/src/coin - the blockchain
/src/daemon - networking and wire protocol
/src/visor - the top level, client
/src/gui - the web wallet and json client interface
/src/wallet - the private key storage library
```

Meshnet
------

```
go run ./cmd/mesh/*.go -config=cmd/mesh/sample/config_a.json

go run ./cmd/mesh/*.go -config=cmd/mesh/sample/config_b.json
```
