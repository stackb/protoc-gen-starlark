.PHONY: build
build:
	bazel build ...

.PHONY: test
test:
	bazel test ... --runs_per_test=30

.PHONY: tidy
tidy:
	go mod tidy
	bazel run update_go_repositories
	bazel run gazelle

golden2:
	bazel run //cmd/protoc-gen-starlark:protoc-gen-starlark_test \
		-- \
		--update

golden:
	bazel run //cmd/grpcstar:grpcstar_test \
		-- \
		--update

.PHONY: serve
serve: build
	GODEBUG=http2debug=2 \
	bazel-bin/cmd/grpcstar/grpcstar_/grpcstar \
		-p bazel-bin/example/routeguide/routeguide_proto_descriptor.pb \
		-f cmd/grpcstar/testdata/routeguide.grpc.star

.PHONY: routeguide_proto_descriptor
routeguide_proto_descriptor:
	bazel build //example/routeguide:routeguide_proto_descriptor
	cp -f bazel-bin/example/routeguide/routeguide_proto_descriptor.pb pkg/starlarkgrpc/

.PHONY: plugin_proto_descriptor
plugin_proto_descriptor:
	bazel build @protoapis//google/protobuf/compiler:plugin_descriptor
	cp -f bazel-bin/external/protoapis/google/protobuf/compiler/plugin_descriptor.pb cmd/protoc-gen-starlark

