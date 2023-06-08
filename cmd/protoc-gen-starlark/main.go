package main

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
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed plugin_descriptor.pb
var pluginDescriptor []byte

func main() {
	if err := run(os.Args[0], os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd string, args []string) error {
	descriptor, err := protodescriptorset.Unmarshal(pluginDescriptor)
	if err != nil {
		return err
	}
	files, err := protodesc.NewFiles(descriptor)
	if err != nil {
		return err
	}

	req, err := readRequest(os.Stdin)
	if err != nil {
		return fmt.Errorf("reading CodeGeneratorRequest: %w", err)
	}
	request, err := protomodule.NewMessage(req)
	if err != nil {
		return fmt.Errorf("starlarkifying CodeGeneratorRequest: %w", err)
	}

	cfg := program.NewConfig()
	cfg.ProtoFiles = files
	cfg.ProtoTypes = protodescriptorset.FileTypes(files)
	cfg.OutputType = program.OutputProto
	cfg.Vars = starlark.StringDict{
		"request": request,
	}
	if err := cfg.ParseArgs(args); err != nil {
		return fmt.Errorf("parsing args: %w", err)
	}

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

	if err := pg.Format(msgs); err != nil {
		return err
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

	// if len(req.GetFileToGenerate()) == 0 {
	// 	return nil, fmt.Errorf("no files were supplied to the generator")
	// }

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
