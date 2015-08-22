// Package script is a library facilitating the creation of programs that resemble
// bash scripts.
package script

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/termie/go-shutil"
)

// Context for script operations. A Context includes the working directory and provides
// access the buffers and results of commands run in the Context.
// Using different Contexts it is possible to handle multiple separate environments.
type Context struct {
	workingDir       string
	lastProcessState *os.ProcessState
	lastProcessError error
	stdoutBuffer     bytes.Buffer
	stderrBuffer     bytes.Buffer
}

// NewContext returns a pointer to a new Context.
func NewContext() (context *Context) {
	// initialize Context
	context = &Context{}

	cwd, err := os.Getwd()
	if err == nil {
		context.SetWorkingDir(cwd)
	}
	return
}

/* WorkingDir */

// SetWorkingDir changes the current working dir.
func (c *Context) SetWorkingDir(workingDir string) {
	c.workingDir = workingDir
}

// WorkingDir retrieves the current working dir.
func (c *Context) WorkingDir() string {
	return c.workingDir
}

/* Commands */

// CommandPath finds the full path of a binary given its name.
func (c *Context) CommandPath(name string) string {
	count := 0
	for count < 10 {
		count++
		c.ExecuteFullySilent("which", name)
		if output := strings.Trim(c.LastOutput().String(), "\n"); output != "" {
			return output
		}
	}
	return ""
}

// CommandExists checks if a given binary exists in PATH.
func (c *Context) CommandExists(name string) bool {
	return c.CommandPath(name) != ""
}

// MustCommandExist ensures a given binary exists in PATH, otherwise panics.
func (c *Context) MustCommandExist(name string) {
	if !c.CommandExists(name) {
		panic(fmt.Errorf("Command %s is not available. Please make sure it is installed and accessible. Ouptut of which: %s", name, c.LastOutput().String()))
	}
}

/* Files and Directories */

// FileExists checks if a given filename exists (being a file).
func (c *Context) FileExists(filename string) bool {
	filename = c.AbsPath(filename)
	fi, err := os.Stat(filename)
	return !os.IsNotExist(err) && !fi.IsDir()
}

// MustFileExist ensures if a given filename exists (being a file), panics otherwise.
func (c *Context) MustFileExist(filename string) {
	if !c.FileExists(filename) {
		panic(fmt.Errorf("File %s does not exist.", filename))
	}
}

// DirExists checks if a given filename exists (being a directory).
func (c *Context) DirExists(path string) bool {
	path = c.AbsPath(path)
	fi, err := os.Stat(path)
	return !os.IsNotExist(err) && fi.IsDir()
}

// MustDirExist checks if a given filename exists (being a directory).
func (c *Context) MustDirExist(path string) {
	if !c.FileExists(path) {
		panic(fmt.Errorf("Directory %s does not exist.", path))
	}
}

// MustGetTempFile guarantees to return a temporary file, panics otherwise
func (c *Context) MustGetTempFile() (tempFile *os.File) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	return
}

// MustGetTempDir guarantees to return a temporary directory, panics otherwise
func (c *Context) MustGetTempDir() (tempDir string) {
	tempDir, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	return
}

// AbsPath returns the absolute path of the path given. If the input path
// is absolute, it is returned untouched. Otherwise the absolute path is
// built relative to the current working directory of the Context.
func (c *Context) AbsPath(filename string) string {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return filename
	}
	isAbsolute := absPath == filename
	if !isAbsolute {
		absPath, err := filepath.Abs(path.Join(c.workingDir, filename))
		if err != nil {
			return filename
		}
		return absPath
	}
	return filename
}

// ResolveSymlinks resolve symlinks in a directory. All symlinked files are
// replaced with copies of the files they point to. Only one level symlinks
// are currently supported.
func (c *Context) ResolveSymlinks(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// symlink?
		if info.Mode()&os.ModeSymlink == os.ModeSymlink {
			// resolve
			linkTargetPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				panic(err)
			}
			targetInfo, err := os.Stat(linkTargetPath)
			if err != nil {
				panic(err)
			}
			os.Remove(path)
			// directory?
			if targetInfo.IsDir() {
				c.CopyDir(linkTargetPath, path)
			} else {
				c.CopyFile(linkTargetPath, path)
			}
		}
		return err
	})
	return err
}

/* Move/Copy Files and Directories */

// MoveFile moves a file. Cross-device moving is supported, so files
// can be moved from and to tmpfs mounts.
func (c *Context) MoveFile(from, to string) error {
	from = c.AbsPath(from)
	to = c.AbsPath(to)

	// work around "invalid cross-device link" for os.Rename
	err := shutil.CopyFile(from, to, true)
	if err != nil {
		return err
	}
	err = os.Remove(from)
	if err != nil {
		return err
	}
	return nil
}

// MoveDir moves a directory. Cross-device moving is supported, so directories
// can be moved from and to tmpfs mounts.
func (c *Context) MoveDir(from, to string) error {
	from = c.AbsPath(from)
	to = c.AbsPath(to)

	// work around "invalid cross-device link" for os.Rename
	options := &shutil.CopyTreeOptions{
		Symlinks:               true,
		Ignore:                 nil,
		CopyFunction:           shutil.Copy,
		IgnoreDanglingSymlinks: false,
	}
	err := shutil.CopyTree(from, to, options)
	if err != nil {
		return err
	}
	err = os.RemoveAll(from)
	if err != nil {
		return err
	}
	return nil
}

// CopyFile copies a file. Cross-device copying is supported, so files
// can be copied from and to tmpfs mounts.
func (c *Context) CopyFile(from, to string) error {
	return shutil.CopyFile(from, to, true) // don't follow symlinks
}

// CopyDir copies a directory. Cross-device copying is supported, so directories
// can be copied from and to tmpfs mounts.
func (c *Context) CopyDir(src, dst string) error {
	options := &shutil.CopyTreeOptions{
		Symlinks:               true,
		Ignore:                 nil,
		CopyFunction:           shutil.Copy,
		IgnoreDanglingSymlinks: false,
	}
	err := shutil.CopyTree(src, dst, options)
	return err
}

/* Executing Commands */

// ExecuteNormal execute a system command, stdout and stderr are output just
// normally
func (c *Context) ExecuteNormal(name string, args ...string) (err error) {
	err = c.Execute(false, false, name, args...)
	return
}

// ExecuteSilent executes a  system command without outputting stdout (it is
// still captured and can be retrieved using LastOutput())
func (c *Context) ExecuteSilent(name string, args ...string) (err error) {
	err = c.Execute(true, false, name, args...)
	return
}

// ExecuteFullySilent executes a system command without outputting stdout or
// stderr (both are still captured and can be retrieved using LastOutput() and
// LastError())
func (c *Context) ExecuteFullySilent(name string, args ...string) (err error) {
	err = c.Execute(true, true, name, args...)
	return
}

// MustExecute ensures a system command to be executed, otherwise panics
func (c *Context) MustExecute(name string, args ...string) {
	err := c.Execute(false, false, name, args...)
	if err != nil {
		panic(err)
	}
}

// MustExecuteSilent ensures a system command to be executed without outputting
// stdout, otherwise panics
func (c *Context) MustExecuteSilent(name string, args ...string) {
	err := c.Execute(true, false, name, args...)
	if err != nil {
		panic(err)
	}
}

// MustExecuteFullySilent ensures a system command to be executed without
// outputting stdout and stderr, otherwise panics
func (c *Context) MustExecuteFullySilent(name string, args ...string) {
	err := c.Execute(true, true, name, args...)
	if err != nil {
		panic(err)
	}
}

// Execute exceutes a system command with configurable stdout and stderr output
func (c *Context) Execute(stdoutSilent bool, stderrSilent bool, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Dir = c.workingDir

	// handling Stdout and Stderr
	// idea from http://nathanleclaire.com/blog/2014/12/29/shelled-out-commands-in-golang/
	cmdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	outScanner := bufio.NewScanner(cmdOutReader)
	c.stdoutBuffer = *new(bytes.Buffer)
	go func() {
		outputHandler(outScanner, !stdoutSilent, c.stdoutBuffer)
	}()

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	errScanner := bufio.NewScanner(cmdErrReader)
	c.stderrBuffer = *new(bytes.Buffer)
	go func() {
		outputHandler(errScanner, !stderrSilent, c.stderrBuffer)
	}()
	cmd.Start()
	err = cmd.Wait()

	// flush buffers
	outputHandler(outScanner, !stdoutSilent, c.stdoutBuffer)
	outputHandler(errScanner, !stderrSilent, c.stderrBuffer)
	c.lastProcessState = cmd.ProcessState
	return err
}

// internal
func outputHandler(scanner *bufio.Scanner, output bool, buffer bytes.Buffer) {
	for scanner.Scan() {
		if output {
			fmt.Println(scanner.Text())
		}
		buffer.WriteString(scanner.Text() + "\n")
	}
}

// LastOutput returns the output buffer of the last command executed using one
// of the Execute* functions. stdout is captured for any command run.
func (c *Context) LastOutput() *bytes.Buffer {
	return &c.stdoutBuffer
}

// LastError returns the error buffer of the last command executed using one
// of the Execute* functions. stderr is captured for any command run.
func (c *Context) LastError() *bytes.Buffer {
	return &c.stderrBuffer
}

// LastSuccessful determins if the last command run was successful. Success is
// defined as the process' return code being zero.
func (c *Context) LastSuccessful() bool {
	return c.LastExitCode() == 0
}

// LastExitCode returns the exit code of the last command run.
func (c *Context) LastExitCode() int {
	var waitStatus syscall.WaitStatus
	if exitError, ok := c.lastProcessError.(*exec.ExitError); ok {
		waitStatus = exitError.Sys().(syscall.WaitStatus)
	} else {
		waitStatus = c.lastProcessState.Sys().(syscall.WaitStatus)
	}
	return waitStatus.ExitStatus()
}

// LastProcessState returns the *os.ProcessState of the last command run.
func (c *Context) LastProcessState() *os.ProcessState {
	return c.lastProcessState
}

// PrintLastState conveniently prints statistics on the last command run. Useful
// for debugging purposes.
func (c *Context) PrintLastState() {
	state := c.lastProcessState
	fmt.Println(
		"PID:", state.Pid(),
		"Exited:", state.Exited(),
		"Exit Code:", c.LastExitCode(),
		"Success:", state.Success(),
		"User Time:", state.UserTime())
}
