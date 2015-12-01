package script_test

import (
	"fmt"
	"github.com/jojomi/go-script"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandExists(t *testing.T) {
	sc := script.NewContext()
	goodBinary := "ls"
	if runtime.GOOS == "windows" {
		goodBinary = "dir"
	}
	badBinary := "not-any-binary-named-like-this"

	assert.Equal(t, true, sc.CommandExists(goodBinary), fmt.Sprintf("Expected command '%s' to exist, but it did not", goodBinary))
	assert.Equal(t, false, sc.CommandExists(badBinary), fmt.Sprintf("Expected command '%s' not to exist, but it did", badBinary))

	sc.MustCommandExist(goodBinary)

	defer func() { recover() }()
	sc.MustCommandExist(badBinary)
}

func TestExists(t *testing.T) {
	testFilename := "test/file.sample"
	testDirectory := "test/dir"

	sc := script.NewContext()

	assert.Equal(t, true, sc.FileExists(testFilename), fmt.Sprintf("Expected file '%s' to exist, but it did not", testFilename))
	sc.MustFileExist(testFilename)

	assert.Equal(t, true, sc.DirExists(testDirectory), fmt.Sprintf("Expected directory '%s' to exist, but it did not", testDirectory))
	sc.MustDirExist(testDirectory)

	defer func() { recover() }()
	sc.MustFileExist(testFilename + "-non-existing")
	sc.MustDirExist(testDirectory + "-non-existing")
}

func TestExecute(t *testing.T) {
	sc := script.NewContext()
	executeFunctions := []func(string, ...string) (*script.ProcessResult, error){
		sc.ExecuteDebug,
		sc.ExecuteSilent,
		sc.ExecuteFullySilent,
	}
	mustExecuteFunctions := []func(string, ...string) *script.ProcessResult{
		sc.MustExecuteDebug,
		sc.MustExecuteSilent,
		sc.MustExecuteFullySilent,
	}
	for _, function := range mustExecuteFunctions {
		pr := function("test/bin/output.sh")
		assert.Equal(t, "output\nalright", pr.Output(), "Unexpected output")
		assert.Equal(t, "error\nwrong", pr.Error(), "Unexpected error output")
	}
	for _, function := range executeFunctions {
		pr, err := function("test/bin/output.sh")
		assert.Equal(t, nil, err, "Command execution returned error")
		assert.Equal(t, true, pr.Successful(), "Command execution unsuccessful")
		assert.NotEqual(t, 0, pr.ProcessState.Pid(), "Command PID incorrect")
		assert.Equal(t, "output\nalright", pr.Output(), "Unexpected output")
		assert.Equal(t, "error\nwrong", pr.Error(), "Unexpected error output")
	}
}

func TestStateString(t *testing.T) {
	sc := script.NewContext()
	sc.ExecuteFullySilent("test/bin/output.sh")
	/*if actual := sc.WorkingDir(); actual != tempDir {
		t.Errorf("Expected WorkingDir: %s, Actual: %s", tempDir, actual)
	}*/
}

func TestWorkingDir(t *testing.T) {
	sc := script.NewContext()
	tempDir := sc.MustGetTempDir()
	defer os.RemoveAll(tempDir)
	sc.SetWorkingDir(tempDir)
	assert.Equal(t, tempDir, sc.WorkingDir(), "Unexpected working dir")

	// copy and move a dir
	from, _ := filepath.Abs("test/dir")
	to := path.Join(tempDir, "dir")
	err := sc.CopyDir(from, to)
	assert.Equal(t, nil, err, "Error on CopyDir")
	checkPaths := []string{
		path.Join(tempDir, "dir", "dir.txt"),
		path.Join(tempDir, "dir", "subdir", "subdir-file"),
	}
	for _, checkFile := range checkPaths {
		assert.Equal(t, true, sc.FileExists(checkFile), fmt.Sprintf("File not existing after CopyDir: %s", checkFile))
	}

	from = to
	to = path.Join(tempDir, "dir-moved")
	err = sc.MoveDir(from, to)
	assert.Equal(t, nil, err)
	checkPaths = []string{
		path.Join(tempDir, "dir-moved", "dir.txt"),
		path.Join(tempDir, "dir-moved", "subdir", "subdir-file"),
	}
	for _, checkFile := range checkPaths {
		assert.Equal(t, true, sc.FileExists(checkFile), fmt.Sprintf("File not existing after MoveDir: %s", checkFile))
	}

	// copy and move a file
	from, _ = filepath.Abs("test/file.sample")
	to = path.Join(tempDir, "output.txt")
	err = sc.CopyFile(from, to)
	assert.Equal(t, nil, err)
	checkPaths = []string{from, to}
	for _, checkFile := range checkPaths {
		assert.Equal(t, true, sc.FileExists(checkFile), fmt.Sprintf("File not existing after CopyFile: %s", checkFile))
	}

	from = to
	to = sc.MustGetTempFile().Name()
	defer os.Remove(to)
	err = sc.MoveFile(from, to)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, sc.FileExists(from), fmt.Sprintf("File existing after MoveFile: %s", from))
	assert.Equal(t, true, sc.FileExists(to), fmt.Sprintf("File not existing after MoveDir: %s", to))
}

func TestSuccessful(t *testing.T) {
	sc := script.NewContext()

	prSuccess, _ := sc.ExecuteFullySilent("test/bin/success.sh")
	assert.Equal(t, true, prSuccess.Successful(), "Command execution should be successful")

	prFail, _ := sc.ExecuteFullySilent("test/bin/fail.sh")
	assert.Equal(t, false, prFail.Successful(), "Command execution should be unsuccessful")
}

func TestResolveSymlinks(t *testing.T) {
	sc := script.NewContext()

	// create symlink for testing
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Symlink(path.Join(wd, "test/dir/dir.txt"), "test/dir/subdir/symlink.txt")
	if err != nil {
		panic(err)
	}
	defer os.Remove("test/dir/subdir/symlink.txt")

	tempDir := sc.MustGetTempDir()
	defer os.RemoveAll(tempDir)
	sc.SetWorkingDir(tempDir)
	from, _ := filepath.Abs("test/dir")
	to := path.Join(tempDir, "dir")
	err = sc.CopyDir(from, to)
	if err != nil {
		panic(err)
	}
	sc.ResolveSymlinks(to)

	symlinkSourcePath := path.Join(to, "dir.txt")
	ioutil.WriteFile(symlinkSourcePath, []byte("test"), 0640)

	// test
	content, err := ioutil.ReadFile(symlinkSourcePath)
	assert.Equal(t, nil, err)

	assert.Equal(t, "test", string(content), "Resolving symlinks did not work.")
	content, err = ioutil.ReadFile(path.Join(to, "subdir", "symlink.txt"))
	assert.Equal(t, "dir.txt content", strings.TrimSpace(string(content)), "Resolving symlinks did not work.")
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
