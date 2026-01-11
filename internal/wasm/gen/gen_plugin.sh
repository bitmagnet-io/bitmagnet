#!/bin/sh

set -e

DIR="$(cd "$(dirname "$0")" && pwd)"

go run "$DIR/cmd" "$@"
