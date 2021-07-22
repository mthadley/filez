package files

import (
	"bytes"
	"html/template"
	"io/fs"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/dustin/go-humanize"
)

type File struct {
	Name    string
	Type    FileType
	Mime    string
	path    string
	base    fs.FS
	size    uint64
	modTime time.Time
}

type FileType = int

const (
	Directory = iota
	SomeFile
	SpecialFile
	Symlink
)

func FromFileInfo(base fs.FS, path string, fileInfo fs.FileInfo) File {
	var type_ FileType
	switch {
	case fileInfo.IsDir():
		type_ = Directory
	case fileInfo.Mode()&fs.ModeSymlink != 0:
		type_ = Symlink
	case fileInfo.Mode().IsRegular():
		type_ = SomeFile
	default:
		type_ = SpecialFile
	}

	name := fileInfo.Name()

	return File{
		Type:    type_,
		Name:    name,
		Mime:    mime.TypeByExtension(filepath.Ext(name)),
		path:    filepath.Join(path, name),
		base:    base,
		size:    uint64(fileInfo.Size()),
		modTime: fileInfo.ModTime(),
	}
}

func (f File) Content() template.HTML {
	content, err := fs.ReadFile(f.base, normalizePath(f.path))
	if err != nil {
		return "Unable to read file"
	}

	fallback := template.HTML("<pre>" + strings.TrimSpace(string(content)) + "</pre>")

	buf := new(bytes.Buffer)

	lexer := lexers.Match(f.Name)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	it, err := lexer.Tokenise(nil, string(content))
	if err != nil {
		return fallback
	}

	formatter := html.New(
		html.Standalone(false),
		html.WithLineNumbers(true),
		html.LinkableLineNumbers(true, "l"),
	)
	style := styles.Friendly

	err = formatter.Format(buf, style, it)
	if err != nil {
		return fallback
	}

	return template.HTML(buf.String())
}

func (f File) EmojiIcon() (icon string) {
	switch f.Type {
	case Directory:
		icon = "📂"
	case SomeFile:
		icon = emojiFromMime(f.Mime)
	case Symlink:
		icon = "🔗"
	case SpecialFile:
		icon = "⚙️"
	}
	return
}

func emojiFromMime(m string) (icon string) {
	switch m {
	case "application/x-ruby":
		icon = "💎"
	case "text/x-python":
		icon = "🐍"
	default:
		switch {
		case strings.HasPrefix(m, "image/"):
			icon = "🖼️"
		case strings.HasPrefix(m, "video/"):
			icon = "🎞️"
		default:
			icon = "📄"
		}
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

type ParentPath struct {
	Name string
	Path string
}

func (f File) ParentPaths() (paths []ParentPath) {
	if f.IsRoot() {
		return
	}

	path := filepath.Dir(f.Path())

	for path != "/" {
		paths = append([]ParentPath{ParentPath{filepath.Base(path), path}}, paths...)
		path = filepath.Dir(path)
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
