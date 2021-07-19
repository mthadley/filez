package files

import (
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type File struct {
	Name    string
	Type    FileType
	path    string
	base    *fs.FS
	size    uint64
	modTime time.Time
}

type FileType = int

const (
	Directory = iota
	SomeFile
)

func fromFileInfo(base *fs.FS, path string, fileInfo fs.FileInfo) File {
	var type_ FileType
	if fileInfo.IsDir() {
		type_ = Directory
	} else {
		type_ = SomeFile
	}

	return File{
		Type:    type_,
		Name:    fileInfo.Name(),
		path:    filepath.Join(path, fileInfo.Name()),
		base:    base,
		size:    uint64(fileInfo.Size()),
		modTime: fileInfo.ModTime(),
	}
}

func (f File) Content() string {
	content, err := fs.ReadFile(*f.base, normalizePath(f.path))
	if err != nil {
		return "Unable to read file"
	}

	return strings.TrimSpace(string(content))
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

func (f File) HumanSize() string {
	return humanize.Bytes(f.size)
}

func (f File) HumanLastModified() string {
	return humanize.Time(f.modTime)
}

func (f File) IsRoot() bool {
	return f.path == "."
}

func (f File) Path() string {
	return "/" + f.path
}

func (f File) ParentPath() string {
	return filepath.Dir(f.Path())
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
