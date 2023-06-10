package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/stackb/protoc-gen-starlark/pkg/plugin"
)

func main() {
	if err := run(os.Args[0], os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cmd string, args []string) error {
	p := plugin.Plugin{}
	return p.Run(args)
}
