.PHONY: build
build:
	bazel build ...

.PHONY: test
test:
	bazel test ... --runs_per_test=30

golden:
	bazel run //cmd/protoc-gen-starlark:protoc-gen-starlark_test -- --update
	bazel run //pkg/plugin:plugin_test -- --update

.PHONY: tidy
tidy:
	go mod tidy
	bazel run update_go_repositories
	bazel run gazelle

.PHONY: plugin_descriptor
plugin_descriptor:
	bazel build @protoapis//google/protobuf/compiler:plugin_descriptor
	cp -f bazel-bin/external/protoapis/google/protobuf/compiler/plugin_descriptor.pb pkg/plugin

.PHONY: unittest_descriptor
unittest_descriptor:
	bazel build @protoapis//google/protobuf:unittest_descriptor
	cp -f bazel-bin/external/protoapis/google/protobuf/unittest_descriptor.pb pkg/plugin


