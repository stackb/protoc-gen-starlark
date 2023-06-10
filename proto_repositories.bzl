load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "proto_repository")

def proto_repositories():
    "proto_deps declares all external proto_repository externals"
    proto_repository(
        name = "googleapis",
        build_directives = [
            "gazelle:exclude google/example/endpointsapis/v1",
            "gazelle:exclude google/cloud/recommendationengine/v1beta1",  # is this a bug?
            "gazelle:proto_language descriptor enable true",
        ],
        build_file_expunge = True,
        build_file_proto_mode = "file",
        cfgs = ["//:rules_proto_config.yaml"],
        imports = ["@protoapis//:imports.csv"],
        override_go_googleapis = True,
        strip_prefix = "googleapis-02710fa0ea5312d79d7fb986c9c9823fb41049a9",
        type = "zip",
        urls = ["https://codeload.github.com/googleapis/googleapis/zip/02710fa0ea5312d79d7fb986c9c9823fb41049a9"],
    )

    proto_repository(
        name = "protoapis",
        build_directives = [
            "gazelle:exclude testdata",
            "gazelle:exclude google/protobuf/compiler/ruby",
            "gazelle:proto_language descriptor enable true",
        ],
        build_file_expunge = True,
        build_file_proto_mode = "file",
        cfgs = ["//:rules_proto_config.yaml"],
        deleted_files = [
            "google/protobuf/unittest_custom_options.proto",
            "google/protobuf/map_lite_unittest.proto",
            "google/protobuf/map_proto2_unittest.proto",
            "google/protobuf/test_messages_proto2.proto",
            "google/protobuf/test_messages_proto3.proto",
        ],
        strip_prefix = "protobuf-9650e9fe8f737efcad485c2a8e6e696186ae3862/src",
        type = "zip",
        urls = ["https://codeload.github.com/protocolbuffers/protobuf/zip/9650e9fe8f737efcad485c2a8e6e696186ae3862"],
    )
