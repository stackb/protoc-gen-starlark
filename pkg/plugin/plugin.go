package plugin

import (
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/stackb/grpc-starlark/pkg/program"
	"github.com/stackb/grpc-starlark/pkg/protodescriptorset"
	"github.com/stripe/skycfg/go/protomodule"
	"go.starlark.net/starlark"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed plugin_descriptor.pb
var pluginDescriptor []byte

type Plugin struct {
	Stdin  io.Reader
	Stdout io.Writer
	Files  *protoregistry.Files // Additional Files to be merged into the final set
}

func (p *Plugin) Run(args []string) error {
	if p.Stdin == nil {
		p.Stdin = os.Stdin
	}
	if p.Stdout == nil {
		p.Stdout = os.Stdout
	}

	files, err := protodescriptorset.ParseFiles(pluginDescriptor)
	if err != nil {
		return err
	}
	files = protodescriptorset.MergeFilesIgnoreConflicts(p.Files, files)
	req, err := readRequest(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading CodeGeneratorRequest: %w", err)
	}
	request, err := protomodule.NewMessage(req)
	if err != nil {
		return fmt.Errorf("starlarkifying CodeGeneratorRequest: %w", err)
	}

	cfg := program.NewConfig()
	if err := cfg.ParseArgs(args); err != nil {
		return fmt.Errorf("parsing args: %w", err)
	}
	cfg.ProtoFiles = protodescriptorset.MergeFilesIgnoreConflicts(cfg.ProtoFiles, files)
	cfg.ProtoTypes = protodescriptorset.FileTypes(cfg.ProtoFiles)
	cfg.OutputType = program.OutputProto
	cfg.Vars = starlark.StringDict{
		"request": request,
	}

	cmd := os.Args[0]
	if cfg.File == "" && fileExists(cmd+".star") {
		cfg.File = cmd + ".star"
	}

	pg, err := program.NewProgram(cfg)
	if err != nil {
		return err
	}

	msgs, err := pg.Exec()
	if err != nil {
		return err
	}

	if len(msgs) != 1 {
		return fmt.Errorf("expected single CodeGeneratorResponse, got %v", len(msgs))
	}

	data, err := proto.Marshal(msgs[0])
	if err != nil {
		return fmt.Errorf("marshaling response: %w", err)
	}
	if _, err := os.Stdout.Write(data); err != nil {
		return fmt.Errorf("writing response: %w", err)
	}
	return nil
}

func readRequest(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading request data: %w", err)
	}

	req := new(pluginpb.CodeGeneratorRequest)
	if err = proto.Unmarshal(data, req); err != nil {
		return nil, fmt.Errorf("unmarshaling CodeGeneratorRequest: %w", err)
	}

	return req, nil
}

// fileExists checks if a file exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info == nil {
		return false
	}
	return !info.IsDir()
}
