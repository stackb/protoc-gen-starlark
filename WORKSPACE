workspace(name = "build_stack_protoc_gen_starlark")

load("//:repositories.bzl", "repositories")

repositories()

# ----------------------------------------------------
# @rules_proto
# ----------------------------------------------------

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies")

rules_proto_dependencies()

# ----------------------------------------------------
# @io_bazel_rules_go
# ----------------------------------------------------

load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(version = "1.18.2")

# ----------------------------------------------------
# @build_stack_rules_proto
# ----------------------------------------------------

register_toolchains("@build_stack_rules_proto//toolchain:standard")

load("//:proto_repositories.bzl", "proto_repositories")

proto_repositories()

# ----------------------------------------------------
# @build_stack_rules_proto
# ----------------------------------------------------

load("@build_stack_rules_proto//:go_deps.bzl", "gazelle_protobuf_extension_go_deps")

gazelle_protobuf_extension_go_deps()

load("@build_stack_rules_proto//deps:go_core_deps.bzl", "go_core_deps")

go_core_deps()

# ----------------------------------------------------
# external go dependencies
# ----------------------------------------------------

load("//:go_repositories.bzl", "go_repositories")

go_repositories()

# ----------------------------------------------------
# @bazel_gazelle
# ----------------------------------------------------
# gazelle:repository_macro go_repositories.bzl%go_repositories

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
