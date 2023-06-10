[![CI](https://github.com/stackb/protoc-gen-starlark/actions/workflows/ci.yaml/badge.svg)](https://github.com/stackb/protoc-gen-starlark/actions/workflows/ci.yaml)

# protoc-gen-starlark

<table border="0">
  <tr>
    <td><img src="https://user-images.githubusercontent.com/50580/141900696-bfb2d42d-5d2c-46f8-bd9f-06515969f6a2.png" height="120"/></td>
    <td><img src="https://static.vecteezy.com/system/resources/previews/007/038/145/non_2x/nightingale-singing-tune-song-bird-musical-notes-music-concept-icon-in-circle-round-black-color-illustration-flat-style-image-vector.jpg" height="120"/></td>
    <!-- image credit: https://www.vecteezy.com/vector-art/7038145-nightingale-singing-tune-song-bird-musical-notes-music-concept-icon-in-circle-round-black-color-vector-illustration-flat-style-image -->
  </tr>
  <tr>
    <td>protobuf</td>
    <td>starlark</td>
  </tr>
</table>

`protoc-gen-starlark` is a scriptable protocol buffer plugin.  It's arguably the easiest way to write a protoc plugin.

## Installation

Download a binary from the [releases
page](https://github.com/stackb/protoc-gen-starlark/releases), or install from source:

```sh
go install github.com/stackb/protoc-gen-starlark/cmd/protoc-gen-starlark@latest
```

## Usage

`protoc-gen-starlark` works like any other typical protoc plugin: it reads an
encoded `CodeGeneratorRequest` from stdin and writes an encoded
`CodeGeneratorResponse` to stdout.

The logic of generating a response is performed within a starlark script that
you must write.  The simplest such script looks something like:

```py
pb = proto.package("google.protobuf.compiler")

def generate(request):
    """generate prepares the response.

    Args:
      request: the pb.CodeGeneratorRequest that was read from stdin.
    Returns:
      a pb.CodeGeneratorRequest
    """
    return pb.CodeGeneratorResponse(
        error = "not implemented",
    )

def main(ctx):
    """main is the entrypoint function.

    Args:
      ctx: the script context.  It has a struct member named
      'vars' which is a StringDict of variables injected into
      the entrypoint.  vars will contain an entry named "request"
      that is the pb.CodeGeneratorRequest read from stdin.
    Returns:
      A single pb.CodeGeneratorResponse.  The return value from 
      `main` must be a list however, so it is wrapped in a list.

    """
    return [generate(ctx.vars["request"])]
```

Although starlark is an interpreter language, the protobuf message types are stongly typed: it is an error to set/get fields that are not part of the message definition.  See [stackb/grpc-starlark](https://github.com/stackb/grpc-starlark) and [stripe/skycfg](https://github.com/stripe/skycfg) for more details about this.

The sample protoc invocation might look something like:

```sh
$ export PROTOC_GEN_STARLARK_SCRIPT=foo.star
$ protoc \
  --foo_out=./gendir \
  --plugin=protoc-gen-foo=/usr/bin/protoc-gen-starlark
```

In this case `protoc-gen-starlark` discovers which script to evaluate using the
`PROTOC_GEN_STARLARK_SCRIPT` environment variable.

Another strategy is to copy/rename `protoc-gen-starlark` and the script to a
common name.  If a file named `$0.star` exists where `$0` is the name of the
executable itself, this will be loaded.  For example:

```sh
$ ln -s /usr/bin/protoc-gen-starlark tools/protoc-gen-foo
$ mv foo.plugin.star                 tools/protoc-gen-foo.star
$ ln -s /usr/bin/protoc-gen-starlark tools/protoc-gen-bar
$ mv bar.plugin.star                 tools/protoc-gen-bar.star

$ protoc \
  --foo_out=./gendir \
  --plugin=protoc-gen-foo=./tools/protoc-gen-foo
  --bar_out=./gendir \
  --plugin=protoc-gen-bar=./tools/protoc-gen-bar
```

