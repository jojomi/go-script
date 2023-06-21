package script

import (
	"fmt"
	"github.com/juju/errors"
	"io"
	"os"
)

func (x File) IsReadable() bool {
	f, err := os.Open(x.AbsPath())
	if err == nil {
		_ = f.Close()
	}
	return err == nil
}

func (x File) AssertReadable() File {
	if !x.IsReadable() {
		panic(fmt.Errorf("file %s should have been readable", x))
	}
	return x
}

func (x File) IsWritable() bool {
	filePath := x.AbsPath()
	existed := x.Exists()
	f, err := os.OpenFile(filePath, os.O_RDWR, x.createPermissions)
	if err == nil {
		_ = f.Close()
		if !existed {
			err = os.Remove(filePath)
			if err != nil {
				panic(err)
			}
		}
	}
	return err == nil
}

func (x File) AssertWritable() File {
	if !x.IsWritable() {
		panic(fmt.Errorf("file %s should have been writable", x))
	}
	return x
}

func (x File) Content() ([]byte, error) {
	return os.ReadFile(x.AbsPath())
}

func (x File) MustContent() []byte {
	content, err := os.ReadFile(x.AbsPath())
	if err != nil {
		panic(errors.Annotatef(err, "could not read content of %s", x))
	}
	return content
}

func (x File) StringContent() (string, error) {
	content, err := os.ReadFile(x.AbsPath())
	return string(content), err
}

func (x File) MustStringContent() string {
	content, err := os.ReadFile(x.AbsPath())
	if err != nil {
		panic(errors.Annotatef(err, "could not read content of %s", x))
	}
	return string(content)
}

func (x File) Append(newContent []byte) error {
	f, err := os.OpenFile(x.AbsPath(), os.O_RDWR, x.createPermissions)
	if err != nil {
		return err
	}

	// seek to end of file
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	_, err = f.Write(newContent)
	return err
}

func (x File) AppendString(newContent string) error {
	return x.Append([]byte(newContent))
}

func (x File) AppendStringln(newContent string) error {
	return x.AppendString(newContent + "\n")
}

func (x File) SetContent(newContent []byte) error {
	return os.WriteFile(x.AbsPath(), newContent, x.createPermissions)
}

func (x File) SetContentString(newContent string) error {
	return x.SetContent([]byte(newContent))
}
