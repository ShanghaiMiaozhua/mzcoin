#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"


echo "DIR:" "$DIR"

pushd "$DIR" >/dev/null
#./scripts/clean-static-libs.sh >/dev/null 2>&1
go run cmd/mzcoin/mzcoin.go --gui-dir="${DIR}/src/gui/static/" $@

popd >/dev/null
