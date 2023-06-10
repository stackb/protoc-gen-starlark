package plugin

import (
	_ "embed"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/stackb/grpc-starlark/pkg/protodescriptorset"
	libtime "go.starlark.net/lib/time"
)

var update = flag.Bool("update", false, "update golden files")

//go:embed unittest_descriptor.pb
var unittestDescriptor []byte

func TestGoldens(t *testing.T) {
	flag.Parse()
	workspaceDir := os.Getenv("BUILD_WORKING_DIRECTORY")

	start := time.Now()
	epoch := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)

	libtime.NowFunc = func() time.Time {
		delta := time.Since(start).Round(100 * time.Millisecond)
		now := epoch.Add(delta)
		return now
	}

	type goldenTest struct {
		file          string
		goldenOutFile string
		goldenErrFile string
		outFile       string
		errFile       string
	}
	var tests []*goldenTest

	entries, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no test files found!")
	}

	for _, file := range entries {
		if strings.HasSuffix(file.Name(), ".plugin.star") {
			tests = append(tests, &goldenTest{
				file:          file.Name(),
				goldenOutFile: file.Name() + ".out",
				goldenErrFile: file.Name() + ".err",
				outFile:       file.Name() + ".stdout",
				errFile:       file.Name() + ".stderr",
			})
		}
	}
	if len(tests) == 0 {
		t.Fatal("no tests found!")
	}

	for _, pair := range tests {
		t.Run(pair.file, func(t *testing.T) {
			outFile, err := os.Create(pair.outFile)
			if err != nil {
				t.Fatal(err)
			}
			errFile, err := os.Create(pair.errFile)
			if err != nil {
				t.Fatal(err)
			}
			stdout := os.Stdout
			stderr := os.Stderr
			os.Stdout = outFile
			os.Stderr = errFile
			defer func() {
				os.Stderr = stdout
				os.Stderr = stderr
			}()

			files, err := protodescriptorset.ParseFiles(unittestDescriptor)
			if err != nil {
				t.Fatal(err)
			}

			plugin := &Plugin{
				Files: files,
			}
			if err := plugin.Run([]string{
				"-file=" + filepath.Join("testdata", pair.file),
				"-o", "stablejson",
			}); err != nil {
				t.Fatal(err)
			}
			outFile.Close()
			errFile.Close()

			gotOut, err := os.ReadFile(pair.outFile)
			if err != nil {
				t.Fatal("reading out file:", err)
			}
			gotErr, err := os.ReadFile(pair.errFile)
			if err != nil {
				t.Fatal("reading err file:", err)
			}

			if *update {
				if workspaceDir == "" {
					t.Fatal("BUILD_WORKING_DIRECTORY not set!")
				}
				dir := filepath.Join(workspaceDir, "pkg", "plugin", "testdata")
				if err := os.WriteFile(filepath.Join(dir, pair.goldenOutFile), gotOut, os.ModePerm); err != nil {
					t.Fatal("writing goldenOut file:", err)
				}
				if err := os.WriteFile(filepath.Join(dir, pair.goldenErrFile), gotErr, os.ModePerm); err != nil {
					t.Fatal("writing goldenErr file:", err)
				}
			} else {
				wantOut, err := os.ReadFile(filepath.Join("testdata", pair.goldenOutFile))
				if err != nil {
					t.Fatal("reading goldenOut file:", err)
				}
				wantErr, err := os.ReadFile(filepath.Join("testdata", pair.goldenErrFile))
				if err != nil {
					t.Fatal("reading goldenErr file:", err)
				}
				if diff := cmp.Diff(string(wantOut), string(gotOut)); diff != "" {
					t.Errorf("stdout (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(string(wantErr), string(gotErr)); diff != "" {
					t.Log("want stderr:\n", string(wantErr))
					t.Log("got stderr:\n", string(gotErr))
					t.Errorf("stderr (-want +got):\n%s", diff)
				}
			}
		})
	}
}
