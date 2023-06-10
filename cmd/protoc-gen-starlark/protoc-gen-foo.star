pb = proto.package("google.protobuf.compiler")

def generate_foo(proto_file):
    return pb.CodeGeneratorResponse.File(
        name = proto_file.name + ".foo.txt",
        content = "GENERATED FILE - DO NOT EDIT",
    )

def generate(request):
    generated_files = []
    for proto_file in request.proto_file:
        generated_files.append(generate_foo(proto_file))
    return [pb.CodeGeneratorResponse(
        file = generated_files,
    )]

def main(ctx):
    return generate(ctx.vars["request"])
