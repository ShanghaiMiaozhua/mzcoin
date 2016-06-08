#!/usr/bin/env bash

GOCMD="go test"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

declare -a libs=(./src/lib/secp256k1-go)
declare -a pkgs=(./src/cipher ./src/gui ./src/util ./src/coin 
                ./src/daemon ./src/mzcoin ./src/visor)
declare -a cmds=(./cmd/mzcoin
                 ./cmd/blockchain ./cmd/address ./cmd/blocksigs
                 ./cmd/genesis ./cmd/cert ./cmd/wallet)

pushd "$DIR" >/dev/null

for i in "${pkgs[@]}" 
do
    $GOCMD "$i"
done

for i in "${libs[@]}" 
do
    $GOCMD "$i"
done

for i in "${cmds[@]}"
do
    $GOCMD "$i"
done

popd >/dev/null
