package files

import "io/ioutil"

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
