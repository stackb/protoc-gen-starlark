pb = proto.package("google.protobuf.compiler")

def generate(request):
    return [pb.CodeGeneratorResponse(
        error = "Nothing to see here!",
    )]

def main(ctx):
    # return generate(ctx.vars.request)
    return generate(pb.CodeGeneratorRequest())
