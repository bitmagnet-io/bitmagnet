#!/bin/sh

set -e

DIR="$(cd "$(dirname "$0")" && pwd)"

$DIR/gen_protoc.sh \
  --go-plugin_out="$DIR/fixtures" \
  -I="$DIR/fixtures" \
  "$DIR/fixtures/imported/imported.proto" \
  "$DIR/fixtures/api/api.proto" \
  "$DIR/fixtures/host/host.proto"
