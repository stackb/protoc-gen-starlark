pb = proto.package("google.protobuf")
compilerpb = proto.package("google.protobuf.compiler")

def generate_enum_type(enum):
    return [
        "## enum: " + enum.name,
        "",
    ]

def generate_message_type(message):
    return [
        "## message: " + message.name,
        "",
    ]

def generate_md(proto_file):
    lines = []
    for m in proto_file.message_type:
        lines += generate_message_type(m)
    for m in proto_file.enum_type:
        lines += generate_enum_type(m)

    return compilerpb.CodeGeneratorResponse.File(
        name = proto_file.name + ".md",
        content = "\n".join(lines),
    )

def generate(request):
    generated_files = []
    for file in request.proto_file:
        generated_files.append(generate_md(file))

    return [compilerpb.CodeGeneratorResponse(
        file = generated_files,
    )]

def main(ctx):
    return generate(ctx.vars["request"])
