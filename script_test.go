package script_test

import (
	_ "fmt"
	"github.com/jojomi/go-script"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCommandExists(t *testing.T) {
	sc := script.NewContext()
	goodBinary := "ls"
	if runtime.GOOS == "windows" {
		goodBinary = "dir"
	}
	badBinary := "not-any-binary-named-like-this"

	if !sc.CommandExists(goodBinary) {
		t.Errorf("Expected command '%s' to exists, but it did not", goodBinary)
	}
	if sc.CommandExists(badBinary) {
		t.Errorf("Expected command '%s' not to exists, but it did", badBinary)
	}

	sc.MustCommandExist(goodBinary)

	defer func() { recover() }()
	sc.MustCommandExist(badBinary)
}

func TestExists(t *testing.T) {
	testFilename := "test/file.sample"
	testDirectory := "test/dir"

	sc := script.NewContext()
	if !sc.FileExists(testFilename) {
		t.Errorf("Expected file '%s' to exist, but it did not", testFilename)
	}
	sc.MustFileExist(testFilename)

	if !sc.DirExists(testDirectory) {
		t.Errorf("Expected directory '%s' to exist, but it did not", testDirectory)
	}
	sc.MustDirExist(testDirectory)

	defer func() { recover() }()
	sc.MustFileExist(testFilename + "-non-existing")
	sc.MustDirExist(testDirectory + "-non-existing")
}

func TestExecute(t *testing.T) {
	sc := script.NewContext()
	executeFunctions := []func(string, ...string) error{
		sc.ExecuteDebug,
		sc.ExecuteSilent,
		sc.ExecuteFullySilent,
	}
	mustExecuteFunctions := []func(string, ...string){
		sc.MustExecuteDebug,
		sc.MustExecuteSilent,
		sc.MustExecuteFullySilent,
	}
	for _, function := range mustExecuteFunctions {
		function("test/output.sh")

		expected := "output\nalright"
		if actual := sc.LastOutput(); actual != expected {
			t.Errorf("Expected Output: %s, Actual: %s", expected, actual)
		}
		expected = "error\nwrong"
		if actual := sc.LastError(); actual != expected {
			t.Errorf("Expected Error: %s, Actual: %s", expected, actual)
		}
	}
	for _, function := range executeFunctions {
		function("test/output.sh")
		if !sc.LastSuccessful() {
			t.Errorf("Command execution unsuccessful")
		}
		if sc.LastProcessState().Pid() == 0 {
			t.Errorf("Command PID incorrect")
		}
		expected := "output\nalright"
		if actual := sc.LastOutput(); actual != expected {
			t.Errorf("Expected Output: %s, Actual: %s", expected, actual)
		}
		expected = "error\nwrong"
		if actual := sc.LastError(); actual != expected {
			t.Errorf("Expected Error: %s, Actual: %s", expected, actual)
		}
	}
}

func TestPrintLastState(t *testing.T) {
	sc := script.NewContext()
	sc.ExecuteFullySilent("test/output.sh")
	sc.PrintLastState()
}

func TestWorkingDir(t *testing.T) {
	sc := script.NewContext()
	tempDir := sc.MustGetTempDir()
	defer os.RemoveAll(tempDir)
	sc.SetWorkingDir(tempDir)
	if actual := sc.WorkingDir(); actual != tempDir {
		t.Errorf("Expected WorkingDir: %s, Actual: %s", tempDir, actual)
	}

	// copy and move a dir
	from, _ := filepath.Abs("test/dir")
	to := path.Join(tempDir, "dir")
	err := sc.CopyDir(from, to)
	if err != nil {
		t.Errorf("Error on CopyDir: %q", err)
	}
	checkPaths := []string{
		path.Join(tempDir, "dir", "dir.txt"),
		path.Join(tempDir, "dir", "subdir", "subdir-file"),
	}
	for _, checkFile := range checkPaths {
		if !sc.FileExists(checkFile) {
			t.Errorf("File not existing after CopyDir: %s", checkFile)
		}
	}

	from = to
	to = path.Join(tempDir, "dir-moved")
	err = sc.MoveDir(from, to)
	if err != nil {
		t.Errorf("Error on MoveDir: %q", err)
	}
	checkPaths = []string{
		path.Join(tempDir, "dir-moved", "dir.txt"),
		path.Join(tempDir, "dir-moved", "subdir", "subdir-file"),
	}
	for _, checkFile := range checkPaths {
		if !sc.FileExists(checkFile) {
			t.Errorf("File not existing after MoveDir: %s", checkFile)
		}
	}

	// copy and move a file
	from, _ = filepath.Abs("test/file.sample")
	to = path.Join(tempDir, "output.txt")
	err = sc.CopyFile(from, to)
	if err != nil {
		t.Errorf("Error on CopyFile: %q", err)
	}
	checkPaths = []string{from, to}
	for _, checkFile := range checkPaths {
		if !sc.FileExists(checkFile) {
			t.Errorf("File not existing after CopyFile: %s", checkFile)
		}
	}

	from = to
	to = sc.MustGetTempFile().Name()
	defer os.Remove(to)
	err = sc.MoveFile(from, to)
	if err != nil {
		t.Errorf("Error on MoveFile: %q", err)
	}
	if sc.FileExists(from) {
		t.Errorf("File existing after MoveFile: %s", from)
	}
	if !sc.FileExists(to) {
		t.Errorf("File not existing after MoveDir: %s", to)
	}
}

func TestResolveSymlinks(t *testing.T) {
	sc := script.NewContext()
	tempDir := sc.MustGetTempDir()
	defer os.RemoveAll(tempDir)
	sc.SetWorkingDir(tempDir)
	from, _ := filepath.Abs("test/dir")
	to := path.Join(tempDir, "dir")
	err := sc.CopyDir(from, to)
	if err != nil {
		panic(err)
	}
	sc.ResolveSymlinks(to)

	symlinkSourcePath := path.Join(to, "dir.txt")
	ioutil.WriteFile(symlinkSourcePath, []byte("test"), 0640)

	// test
	content, err := ioutil.ReadFile(symlinkSourcePath)
	if err != nil {
		panic(err)
	}
	actual := string(content)
	expected := "test"
	if actual != expected {
		t.Errorf("Resolving symlinks did not work. File content expected: '%s', actual '%s'", expected, actual)
	}
	content, err = ioutil.ReadFile(path.Join(to, "subdir", "symlink.txt"))
	actual = strings.TrimSpace(string(content))
	expected = "dir.txt content"
	if actual != expected {
		t.Errorf("Resolving symlinks did not work. File content expected: '%s', actual '%s'", expected, actual)
	}
}

func createFile(filename, content string) {
	ioutil.WriteFile(filename, []byte(content), 0640)
}

func createDirectory(name string) {
	err := os.MkdirAll(name, 0640)
	if err != nil {
		panic(err)
	}
}
