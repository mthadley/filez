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
