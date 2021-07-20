package files_test

import (
	"log"
	"reflect"
	"testing"

	"github.com/mthadley/filez/internal/files"
	"github.com/spf13/afero"
)

func TestParentPaths(t *testing.T) {
	file := mockFile("foo/bar/baz.txt")

	expectedParentPaths := []files.ParentPath{
		files.ParentPath{Name: "foo", Path: "/foo"},
		files.ParentPath{Name: "bar", Path: "/foo/bar"},
	}

	if !reflect.DeepEqual(expectedParentPaths, file.ParentPaths()) {
		t.Errorf(
			"Parent paths did not match. Expected %v, got %v",
			expectedParentPaths,
			file.ParentPaths(),
		)
	}
}

func TestParentPathsWithRoot(t *testing.T) {
	file := mockFile("/")

	if file.ParentPaths() != nil {
		t.Errorf("Parent paths did not match. Expected nil, got %v", file.ParentPaths())
	}
}

func mockFile(path string) files.File {
	appFS := afero.NewMemMapFs()
	afero.WriteFile(appFS, path, []byte(""), 0644)
	fs := afero.IOFS{Fs: appFS}

	file, err := files.Info(fs, path)
	if err != nil {
		log.Fatal(err)
	}

	return file
}
