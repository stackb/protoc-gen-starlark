.PHONY: build
build:
	bazel build ...

.PHONY: test
test:
	bazel test ... --runs_per_test=30

golden:
	bazel run //pkg/plugin:plugin_test \
		-- \
		--update

.PHONY: tidy
tidy:
	go mod tidy
	bazel run update_go_repositories
	bazel run gazelle

.PHONY: plugin_proto_descriptor
plugin_proto_descriptor:
	bazel build @protoapis//google/protobuf/compiler:plugin_descriptor
	cp -f bazel-bin/external/protoapis/google/protobuf/compiler/plugin_descriptor.pb pkg/plugin

.PHONY: unittest_proto_descriptor
unittest_proto_descriptor:
	bazel build @protoapis//google/protobuf:unittest_descriptor
	cp -f bazel-bin/external/protoapis/google/protobuf/unittest_descriptor.pb pkg/plugin


