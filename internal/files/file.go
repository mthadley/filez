package files

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type File struct {
	Name string
	Type FileType
	Path string
	base string
}

type FileType = int

const (
	Directory = iota
	SomeFile
)

func fromFileInfo(base, path string, fileInfo fs.FileInfo) File {
	file := File{}

	if fileInfo.IsDir() {
		file.Type = Directory
	} else {
		file.Type = SomeFile
	}

	fmt.Println(path, base)
	file.Name = fileInfo.Name()
	file.base = base
	file.Path = strings.Replace(path+"/"+file.Name, base, "", 1)

	return file
}
func (f File) Content() string {
	content, err := ioutil.ReadFile(f.base + f.Path)
	if err != nil {
		return "Unable to read file"
	}

	return string(content)
}

func (f File) EmojiIcon() (icon string) {
	switch f.Type {
	case Directory:
		icon = "📂"
	case SomeFile:
		icon = "📄"
	}
	return
}

func (f File) IsRoot() bool {
	return f.ParentPath() == "."
}

func (f File) ParentPath() string {
	return filepath.Dir(f.Path)
}

// Custom sort for File: Directories come firs, then sorting by Name alphabetically.
type SortByFileType []File

func (files SortByFileType) Len() int  { return len(files) }
func (f SortByFileType) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (files SortByFileType) Less(i, j int) bool {
	a, b := files[i], files[j]
	if a.Type == b.Type {
		return a.Name < b.Name
	} else {
		return a.Type < b.Type
	}
}
