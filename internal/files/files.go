package files

import (
	"io/fs"
	"path/filepath"
	"sort"
)

func Info(base *fs.FS, path string) (File, error) {
	path = normalizePath(path)

	fileInfo, err := fs.Stat(*base, path)
	if err != nil {
		return File{}, err
	}

	return fromFileInfo(base, filepath.Dir(path), fileInfo), nil
}

func List(base *fs.FS, dir string) ([]File, error) {
	dir = normalizePath(dir)

	entries, err := fs.ReadDir(*base, dir)
	if err != nil {
		return nil, err
	}

	result := make([]File, len(entries))
	for i, dirEntry := range entries {
		fileInfo, err := dirEntry.Info()
		if err != nil {
			return nil, err
		}

		result[i] = fromFileInfo(base, dir, fileInfo)
	}

	sort.Sort(SortByFileType(result))

	return result, nil
}

func normalizePath(path string) string {
	if path == "/" {
		return "."
	}

	if rune(path[0]) == '/' {
		return path[1:]
	}

	return path
}
