pb = proto.package("google.protobuf")
compilerpb = proto.package("google.protobuf.compiler")

fake_request = compilerpb.CodeGeneratorRequest(
    proto_file = [
        pb.FileDescriptorProto(
            name = "a.proto",
        ),
    ],
)

def generate_txt_file(proto_file):
    return compilerpb.CodeGeneratorResponse.File(
        name = proto_file.name + ".txt",
        content = "Fake Content",
    )

def generate(request):
    generated_files = []
    for file in request.proto_file:
        generated_files.append(generate_txt_file(file))

    response = compilerpb.CodeGeneratorResponse(
        file = generated_files,
    )
    return [response]

def main(ctx):
    return generate(ctx.vars.request)
    # return generate(fake_request)
