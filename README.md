MZCoin
=======

MZCoin is a next-generation cryptocurrency.

MZCoin improves on Bitcoin in too many ways to be addressed here.

MZCoin is small part of OP Redecentralize and OP Darknet Plan.

Installation
------------

*For detailed installation instructions, see [Installing Skycoin](../../wiki/Installation)*

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

go run ./cmd/skycoin/skycoin.go
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

The skycoin repo must be in $GOPATH, under "/src/github.com/skycoin". Otherwise golang programs cannot import the libraries.

```
#pull skycoin repo into the gopath
#note: puts the skycoin folder in $GOPATH/src/github.com/skycoin/skycoin
go get github.com/skycoin/skycoin

#create symlink of the repo
cd $HOME
ln -s $GOPATH/src/github.com/skycoin/skycoin skycoin
```

Dependencies
---

```
go get github.com/robfig/glock
glock sync github.com/skycoin/skycoin
go get ./cmd/skycoin
```

To update dependencies
```
glock save github.com/skycoin/skycoin/cmd/skycoin
```

Running Node
---

```
cd skycoin
screen
go run ./cmd/skycoin/skycoin.go 
#then ctrl+A then D to exit screen
#screen -x to reattach screen
```

Todo
---

Use gvm package set, so repo does not need to be symlinked. Does this have a default option?

```
gvm pkgset create skycoin
gvm pkgset use skycoin
git clone https://github.com/skycoin/skycoin
cd skycoin
go install
```

Running
---

cd skycoin
go run ./cmd/skycoin/skycoin.go

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

Run the skycoin client then
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

This is a public server. You can use these urls on local host too, with the skycoin client running.
```
http://skycoin-chompyz.c9.io/outputs
http://skycoin-chompyz.c9.io/blockchain/blocks?start=0&end=500
http://skycoin-chompyz.c9.io/blockchain
http://skycoin-chompyz.c9.io/connections
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
