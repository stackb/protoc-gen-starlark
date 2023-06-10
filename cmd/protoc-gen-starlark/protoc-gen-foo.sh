#!/bin/bash

set -euox pipefail

./cmd/protoc-gen-starlark/protoc-gen-unittest \
    -file ./cmd/protoc-gen-starlark/protoc-gen-foo.star
