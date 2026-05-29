#!/usr/bin/env sh

set -e

AIR_BIN="${GOBIN:-${GOPATH:-$HOME/go}/bin}/air"

exec "$AIR_BIN" "$@"
