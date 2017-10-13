package script

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var myFileFileMode = os.FileMode(int(0700))

func TestFileExists(t *testing.T) {
	path := "/tmp/my-path/subdir"
	filename := "my-file"
	fullFilename := filepath.Join(path, filename)
	sc := NewContext()
	sc.fs = afero.NewMemMapFs()

	// make a file and check with absolute paths
	sc.SetWorkingDir("/not-needed")
	makeDirectory(sc, path)
	assert.False(t, sc.FileExists(fullFilename))
	makeFile(sc, fullFilename, "")
	assert.True(t, sc.FileExists(fullFilename))

	// check with relative paths using working dir
	sc.SetWorkingDir(path)
	sc.fs = afero.NewMemMapFs()
	makeDirectory(sc, path)
	assert.False(t, sc.FileExists(filename))
	makeFile(sc, fullFilename, "")
	assert.True(t, sc.FileExists(filename))
}

func TestDirExists(t *testing.T) {
	path := "/tmp/my-path/subdir"
	dirname := "my-directory"
	fullDirname := filepath.Join(path, dirname)
	sc := NewContext()
	sc.fs = afero.NewMemMapFs()

	// make a directory and check with absolute paths
	sc.SetWorkingDir("/not-needed")
	assert.False(t, sc.DirExists(fullDirname))
	makeDirectory(sc, fullDirname)
	assert.True(t, sc.DirExists(fullDirname))

	// check with relative paths using working dir
	sc.SetWorkingDir(path)
	sc.fs = afero.NewMemMapFs()
	assert.False(t, sc.DirExists(dirname))
	makeDirectory(sc, fullDirname)
	assert.True(t, sc.DirExists(dirname))
}

func TestAbsPath(t *testing.T) {
	sc := NewContext()
	sc.SetWorkingDir("/wd")
	assert.Equal(t, "/wd/file", sc.AbsPath("file"))
	assert.Equal(t, "/wd/dir", sc.AbsPath("dir"))
	assert.Equal(t, "/wd/dir", sc.AbsPath("dir/"))
	assert.Equal(t, "/abc/file", sc.AbsPath("/abc/file"))
	assert.Equal(t, "/abc/dir", sc.AbsPath("/abc/dir"))
	assert.Equal(t, "/abc/dir", sc.AbsPath("/abc/dir/"))
}

func TestAbsPathSep(t *testing.T) {
	sc := NewContext()
	sc.SetWorkingDir("/wd")
	assert.Equal(t, "/wd/dir/", sc.AbsPathSep("dir"))
	assert.Equal(t, "/wd/dir/", sc.AbsPathSep("dir/"))
	assert.Equal(t, "/abc/dir/", sc.AbsPathSep("/abc/dir"))
	assert.Equal(t, "/abc/dir/", sc.AbsPathSep("/abc/dir/"))
}

func TestWithTrailingPathSep(t *testing.T) {
	sc := NewContext()
	assert.Equal(t, "dir/", sc.WithTrailingPathSep("dir"))
	assert.Equal(t, "dir/", sc.WithTrailingPathSep("dir/"))
	assert.Equal(t, "/abc/dir/", sc.WithTrailingPathSep("/abc/dir"))
	assert.Equal(t, "/abc/dir/", sc.WithTrailingPathSep("/abc/dir/"))
}

func TestWithoutTrailingPathSep(t *testing.T) {
	sc := NewContext()
	assert.Equal(t, "dir", sc.WithoutTrailingPathSep("dir"))
	assert.Equal(t, "dir", sc.WithoutTrailingPathSep("dir/"))
	assert.Equal(t, "/abc/dir", sc.WithoutTrailingPathSep("/abc/dir"))
	assert.Equal(t, "/abc/dir", sc.WithoutTrailingPathSep("/abc/dir/"))
}

func TestResolveSymlinks(t *testing.T) {
	sc := NewContext()
	fs := afero.NewMemMapFs()
	sc.fs = fs
	sc.SetWorkingDir("/test")

	err := sc.ResolveSymlinks("dir-non-existing")
	assert.Nil(t, err)

	fs.MkdirAll("/test/dir", 0700)
	afero.WriteFile(fs, "/test/dir/file.txt", []byte("This is my content"), os.FileMode(0644))
	err = sc.ResolveSymlinks("dir")
	assert.Nil(t, err)

	// TODO test actual symlinks (can afero mock os.Symlink to its MemFS?)
}

func TestEnsureDirExists(t *testing.T) {
	path := "/start"
	dir := "abcde"
	fullPath := filepath.Join(path, dir)
	sc := NewContext()
	sc.fs = afero.NewMemMapFs()
	sc.SetWorkingDir(path)
	err := sc.EnsureDirExists(dir, myFileFileMode)
	assert.Nil(t, err)
	assert.True(t, sc.DirExists(fullPath))
}

// TODO func TestEnsureDirExistsFailure(t *testing.T) {}

func TestEnsurePathForFile(t *testing.T) {
	path := "/root/"
	file := "xyz.zip"
	sc := NewContext()
	sc.fs = afero.NewMemMapFs()
	sc.SetWorkingDir(path)
	err := sc.EnsurePathForFile(file, myFileFileMode)
	assert.Nil(t, err)
	assert.True(t, sc.DirExists(path))
}

func TestMoveFile(t *testing.T) {
	fileA := "/dir1/FileA"
	fileB := "/dir2/FileB"
	sc := NewContext()
	sc.fs = afero.NewMemMapFs()
	sc.SetWorkingDir("/outside")
	sc.EnsurePathForFile(fileA, myFileFileMode)
	sc.EnsurePathForFile(fileB, myFileFileMode)

	makeFile(sc, fileA, "insidethefile")
	assert.True(t, fileExists(sc, fileA))
	assert.False(t, fileExists(sc, fileB))

	sc.MoveFile(fileA, fileB)

	assert.False(t, fileExists(sc, fileA))
	assert.True(t, fileExists(sc, fileB))
}

func TestCopyFile(t *testing.T) {
	fileA := "/dir1/FileA"
	fileB := "/dir2/FileB"
	sc := NewContext()
	sc.fs = afero.NewMemMapFs()
	sc.SetWorkingDir("/outside")
	sc.EnsurePathForFile(fileA, myFileFileMode)
	sc.EnsurePathForFile(fileB, myFileFileMode)

	makeFile(sc, fileA, "insidethefile")
	assert.True(t, fileExists(sc, fileA))
	assert.False(t, fileExists(sc, fileB))

	sc.CopyFile(fileA, fileB)

	assert.True(t, fileExists(sc, fileA))
	assert.True(t, fileExists(sc, fileB))
}

func TestMoveDir(t *testing.T) {
	dirA := "/dir1"
	dirB := "/dir2"
	fileA := filepath.Join(dirA, "FileA")
	fileB := filepath.Join(dirA, "subdir/FileB")
	fileAAfter := filepath.Join(dirB, "FileA")
	fileBAfter := filepath.Join(dirB, "subdir/FileB")

	sc := NewContext()
	sc.fs = afero.NewMemMapFs()
	sc.SetWorkingDir("/outside")
	sc.EnsurePathForFile(fileA, myFileFileMode)
	sc.EnsurePathForFile(fileB, myFileFileMode)

	makeFile(sc, fileA, "insidethefilea")
	makeFile(sc, fileB, "insidethefileb")
	assert.True(t, fileExists(sc, fileA))
	assert.True(t, fileExists(sc, fileB))

	sc.MoveDir(dirA, dirB)

	assert.False(t, fileExists(sc, fileA))
	assert.False(t, fileExists(sc, fileB))
	assert.True(t, fileExists(sc, fileAAfter))
	assert.True(t, fileExists(sc, fileBAfter))
}

func TestCopyDir(t *testing.T) {
	dirA := "/dir1"
	dirB := "/dir2"
	fileA := filepath.Join(dirA, "FileA")
	fileB := filepath.Join(dirA, "subdir/FileB")
	fileAAfter := filepath.Join(dirB, "FileA")
	fileBAfter := filepath.Join(dirB, "subdir/FileB")

	sc := NewContext()
	sc.fs = afero.NewMemMapFs()
	sc.SetWorkingDir("/outside")
	sc.EnsurePathForFile(fileA, myFileFileMode)
	sc.EnsurePathForFile(fileB, myFileFileMode)

	makeFile(sc, fileA, "insidethefilea")
	makeFile(sc, fileB, "insidethefileb")
	assert.True(t, fileExists(sc, fileA))
	assert.True(t, fileExists(sc, fileB))

	sc.CopyDir(dirA, dirB)

	assert.True(t, fileExists(sc, fileA))
	assert.True(t, fileExists(sc, fileB))
	assert.True(t, fileExists(sc, fileAAfter))
	assert.True(t, fileExists(sc, fileBAfter))
}

func TestTempFile(t *testing.T) {
	content := "abcfilecontent"
	sc := NewContext()
	sc.fs = afero.NewMemMapFs()
	file, err := sc.tempFileInternal()
	assert.Nil(t, err)
	file.WriteString(content)
	file.Close()

	// verify
	r, err := afero.ReadFile(sc.fs, file.Name())
	assert.Nil(t, err)
	assert.Equal(t, content, string(r))
}

func TestTempDir(t *testing.T) {
	content := []byte("abcfilecontent")
	sc := NewContext()
	sc.fs = afero.NewMemMapFs()
	dir, err := sc.TempDir()
	assert.Nil(t, err)
	filename := filepath.Join(dir, "myfilename")
	afero.WriteFile(sc.fs, filename, content, myFileFileMode)

	// verify
	r, err := afero.ReadFile(sc.fs, filename)
	assert.Nil(t, err)
	assert.Equal(t, content, r)
}

func makeFile(c *Context, filename, content string) {
	afero.WriteFile(c.fs, filename, []byte(content), myFileFileMode)
}

func makeDirectory(c *Context, path string) {
	c.fs.MkdirAll(path, myFileFileMode)
}

func fileExists(c *Context, filename string) bool {
	r, err := afero.Exists(c.fs, filename)
	return r && err == nil
}
