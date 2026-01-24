#!/bin/sh

set -e

DIR="$(cd "$(dirname "$0")" && pwd)"
SCRIPT_DIR=$DIR/../internal/wasm/gen

$SCRIPT_DIR/gen_protoc.sh \
  --go-plugin_out="$DIR" \
  -I="$DIR" \
  "$DIR/common/http/http.proto" \
  "$DIR/common/model/model.proto" \
  "$DIR/common/search/search.proto" \
  "$DIR/api/api.proto" \
  "$DIR/host/configurator/configurator.proto" \
  "$DIR/host/http_client/http_client.proto" \
  "$DIR/host/search/search.proto"
