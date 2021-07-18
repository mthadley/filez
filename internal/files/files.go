package files

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Info(base, path string) (File, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}

	return fromFileInfo(base, filepath.Dir(path), fileInfo), nil
}

func List(base, dir string) ([]File, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make([]File, len(files))
	for i, file := range files {
		result[i] = fromFileInfo(base, dir, file)
	}

	return result, nil
}

func fromFileInfo(base, path string, fileInfo fs.FileInfo) File {
	file := File{}

	if fileInfo.IsDir() {
		file.Type = Directory
	} else {
		file.Type = SomeFile
	}

	file.Name = fileInfo.Name()
	file.base = base
	file.Path = strings.Replace(path, base, "", 1) + "/" + file.Name

	return file
}

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
