#/bin/bash

set -euox pipefail

file_exists_or_die() {
  if [ ! -e "$1" ]; then
    echo "File '$1' does not exist. Exiting..."
    exit 1
  fi
}

# Arrange
mkdir gen

# Action
./cmd/protoc-gen-starlark/protoc.exe \
    --descriptor_set_in=./cmd/protoc-gen-starlark/unittest_descriptor.pb \
    --foo_out=gen \
    --plugin=protoc-gen-foo=./cmd/protoc-gen-starlark/protoc-gen-foo.sh \
    google/protobuf/unittest.proto

find .

# Assert
file_exists_or_die './gen/google/protobuf/unittest_import.proto.foo.txt'
file_exists_or_die './gen/google/protobuf/unittest.proto.foo.txt'
file_exists_or_die './gen/google/protobuf/unittest_import_public.proto.foo.txt'
