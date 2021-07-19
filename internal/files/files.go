package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
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

	sort.Sort(SortByFileType(result))

	return result, nil
}
