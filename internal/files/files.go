package files

import (
	"io/fs"
	"io/ioutil"
	"os"
)

type File struct {
	Name string
	Type FileType
}

type FileType = int

const (
	Directory = iota
	SomeFile
)

func Info(path string) (File, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return File{}, err
	}

	return fromFileInfo(fileInfo), nil
}

func List(dir string) ([]File, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	result := make([]File, len(files))
	for i, file := range files {
		result[i] = fromFileInfo(file)
	}

	return result, nil
}

func fromFileInfo(fileInfo fs.FileInfo) File {
	file := File{}

	file.Name = fileInfo.Name()
	if fileInfo.IsDir() {
		file.Type = Directory
	} else {
		file.Type = SomeFile
	}

	return file
}
