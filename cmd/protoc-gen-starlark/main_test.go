package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	libtime "go.starlark.net/lib/time"
)

var update = flag.Bool("update", false, "update golden files")

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

	for _, pair := range tests {
		t.Run(pair.file, func(t *testing.T) {
			wd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			listFiles(t, wd)

			tmpdir := t.TempDir()
			gotOut, gotErr, err := runProtoc(wd, tmpdir, pair.file)
			if err != nil {
				t.Log(string(gotOut))
				t.Log(string(gotErr))
				t.Fatal("protoc error:", err)
			}

			if *update {
				if workspaceDir == "" {
					t.Fatal("BUILD_WORKING_DIRECTORY not set!")
				}
				dir := filepath.Join(workspaceDir, "cmd", "protoc-gen-starlark", "testdata")
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

func runProtoc(cwd, dir string, filename string) ([]byte, []byte, error) {
	pluginPath := filepath.Join(cwd, "protoc-gen-starlark.exe")
	cmd := exec.Command("protoc.exe",
		"--descriptor_set_in=descriptor.pb",
		"--starlark_out="+dir,
		"--plugin=plugin-gen-starlark="+pluginPath,
		"google/protobuf/unittest.proto",
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return stdout.Bytes(), stderr.Bytes(), err
	}
	return stdout.Bytes(), stderr.Bytes(), nil
}

func listFiles(t *testing.T, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Logf("%v\n", err)
			return err
		}
		// if info.Mode()&os.ModeSymlink > 0 {
		// 	link, err := os.Readlink(path)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	t.Logf("%s -> %s", path, link)
		// 	return nil
		// }
		t.Log(strings.TrimPrefix(path, dir+"/"))
		return nil
	})
}
