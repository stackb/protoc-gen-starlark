package main

import (
	"bytes"
	"flag"
	"io/fs"
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

type testFile struct {
	path    string
	content []byte
}

type goldenTest struct {
	pluginFile     testFile
	outFile        testFile
	errFile        testFile
	genfiles       []testFile
	goldenOutFile  testFile
	goldenErrFile  testFile
	goldenGenfiles []testFile
}

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

	var tests []*goldenTest

	testdata := "testdata"
	entries, err := os.ReadDir(testdata)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("no test files found!")
	}

	for _, file := range entries {
		if !file.IsDir() {
			continue
		}
		if !strings.HasPrefix(file.Name(), "protoc-gen-") {
			continue
		}

		var tc goldenTest
		tc.pluginFile = testFile{path: filepath.Join(testdata, file.Name(), "plugin.star")}
		tc.errFile = testFile{path: tc.pluginFile.path + ".stderr.tmp"}
		tc.outFile = testFile{path: tc.pluginFile.path + ".stdout.tmp"}
		tc.goldenOutFile = testFile{path: tc.pluginFile.path + ".stdout"}
		tc.goldenErrFile = testFile{path: tc.pluginFile.path + ".stderr"}
		tests = append(tests, &tc)
	}

	for _, tc := range tests {
		t.Run(tc.pluginFile.path, func(t *testing.T) {
			tmpdir := t.TempDir()
			wd, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			// listFiles(t, wd)

			gotOut, gotErr, gotFiles, err := runProtoc(t, wd, tmpdir, tc.pluginFile.path)
			tc.outFile.content = gotOut.Bytes()
			tc.errFile.content = gotErr.Bytes()

			if err != nil {
				t.Log(gotOut.String())
				t.Log(gotErr.String())
				t.Fatal("protoc error:", err)
			}

			if *update {
				if workspaceDir == "" {
					t.Fatal("BUILD_WORKING_DIRECTORY not set!")
				}

				dir := filepath.Join(workspaceDir, "cmd", "protoc-gen-starlark")
				if err := os.WriteFile(filepath.Join(dir, tc.goldenOutFile.path), gotOut.Bytes(), os.ModePerm); err != nil {
					t.Fatal("writing goldenOut file:", err)
				}
				if err := os.WriteFile(filepath.Join(dir, tc.goldenErrFile.path), gotErr.Bytes(), os.ModePerm); err != nil {
					t.Fatal("writing goldenErr file:", err)
				}

				genfiles := filepath.Join(dir, "genfiles")
				for _, gotFile := range gotFiles {
					filename := filepath.Join(genfiles, gotFile.path)
					dirname := filepath.Dir(filename)
					if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
						t.Fatalf("prepare dir %s: %v", filename, err)
					}
					if err := os.WriteFile(filename, gotFile.content, os.ModePerm); err != nil {
						t.Fatalf("writing %s: %v", filename, err)
					}
					t.Log("wrote:", filename)
				}
				panic("wip")
			} else {
				wantOut, err := os.ReadFile(tc.goldenOutFile.path)
				if err != nil {
					t.Fatal("reading goldenOut file:", err)
				}
				wantErr, err := os.ReadFile(tc.goldenErrFile.path)
				if err != nil {
					t.Fatal("reading goldenErr file:", err)
				}

				if diff := cmp.Diff(string(wantOut), gotOut.String()); diff != "" {
					t.Errorf("stdout (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(string(wantErr), gotErr.String()); diff != "" {
					t.Log("want stderr:\n", string(wantErr))
					t.Log("got stderr:\n", gotErr.String())
					t.Errorf("stderr (-want +got):\n%s", diff)
				}

				got := make(map[string]testFile)
				for _, file := range gotFiles {
					got[file.path] = file
				}
				genfiles := filepath.Join(filepath.Dir(tc.pluginFile.path), "genfiles")
				wantFiles := collectFiles(t, genfiles)
				for _, wantFile := range wantFiles {
					gotFile, ok := got[wantFile.path]
					if !ok {
						t.Error("wanted generated file (not found):", wantFile.path)
					}
					delete(got, gotFile.path)
					if diff := cmp.Diff(string(wantFile.content), string(gotFile.content)); diff != "" {
						t.Errorf("%s (-want +got):\n%s", wantFile.path, diff)
					}
				}
				for _, gotFile := range got {
					t.Errorf("unexpected generated file: %s", gotFile.path)
				}
			}
		})
	}
}

func runProtoc(t *testing.T, cwd, dir string, filename string) (stdout bytes.Buffer, stderr bytes.Buffer, files []testFile, err error) {
	listFiles(t, ".")
	cmd := exec.Command("protoc.exe",
		"--proto_path=.",
		"--descriptor_set_in=unittest_descriptor.pb",
		"--unittest_out="+dir,
		"google/protobuf/unittest.proto",
	)
	cmd.Env = []string{
		"PROTOC_GEN_STARLARK_FILE=" + filename,
		"PATH=" + cwd,
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = cwd

	if err = cmd.Run(); err != nil {
		return
	}

	files = collectFiles(t, dir)

	return
}

func collectFiles(t *testing.T, dir string) (files []testFile) {
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files = append(files, testFile{
			path:    strings.TrimPrefix(path[len(dir):], "/"),
			content: data,
		})
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	return
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
