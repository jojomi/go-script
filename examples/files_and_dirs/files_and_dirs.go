package main

import (
	"fmt"

	"github.com/jojomi/go-script/v2"
)

func main() {
	// make sure panics are printed in a human friendly way
	defer script.RecoverFunc()

	sc := script.NewContext()

	files := []script.File{sc.FileAt("/tmp/absolute"), sc.FileAt("~/in-home.jpeg"), sc.FileAt("tmp/relative.png"), sc.FileAt("tmp/archive.tar.xz")}

	for _, file := range files {
		fmt.Println("File", file)
		fmt.Println("File absolute?", file.IsAbs())
		fmt.Println("File absolute", file.AbsPath())
		fmt.Println("Filename", file.Filename())

		newExt := script.FileExtensionFrom(".php")
		checkExt := script.FileExtensionFrom(".png")
		fmt.Println("has any Ext", file.HasAnyExtension())
		fmt.Println("has Ext", checkExt, file.HasExtension(checkExt))
		fmt.Println("with Ext", newExt, file.WithExtension(newExt))

		parentDir := file.Dir()
		fmt.Println("Parent Dir", parentDir, parentDir.AbsPath())
		fmt.Println("Dir exists?", parentDir.Exists())
		fmt.Println("File in Dir?", parentDir.MustFileAt("subdir/file.pdf"))
		subDir := parentDir.MustDirAt("subdir")
		fmt.Println("Dir from Dir?", subDir, subDir.Exists())

		fmt.Println()
		_ = file
	}

	// chaining test
	f := sc.DirAt("~").
		MustFileAt("testfile.log").
		AssertExists().
		AssertExtension(script.Log).
		AssertReadable()

	// writing data
	fmt.Println(f.MustStringContent())
	fmt.Println(f.AssertWritable().AppendStringln("abc"))
	fmt.Println(f.MustStringContent())
	fmt.Println(f.SetContentString("content\n"))
	fmt.Println(f.MustStringContent())
}
