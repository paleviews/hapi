#!/usr/bin/env sh

rm -rf annotations/*

find proto -name "*.proto" | xargs \
  protoc --proto_path proto \
    --go_out annotations \
    --go_opt module=github.com/paleviews/hapi/descriptor/annotations
