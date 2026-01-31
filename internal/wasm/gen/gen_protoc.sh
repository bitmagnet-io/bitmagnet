#!/bin/sh

set -e

DIR="$(cd "$(dirname "$0")" && pwd)"

protoc \
  --plugin=protoc-gen-go-plugin=$DIR/gen_plugin.sh \
  --go-plugin_opt=paths=source_relative \
  "$@"
