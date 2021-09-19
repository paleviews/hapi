#!/usr/bin/env sh

rm -rf golang/*

find proto -name "*.proto" | xargs \
  protoc --proto_path proto --proto_path ../../../descriptor/proto \
    --hapi_opt paths=source_relative \
    --hapi_opt doc_file=doc/doc.yaml \
    --hapi_out golang
