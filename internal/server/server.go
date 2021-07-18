package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mthadley/filez/internal/files"
)

type Server struct {
	baseDir string
	router  *mux.Router
}

func NewServer(dir string) *Server {
	return &Server{
		baseDir: dir,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.router == nil {
		s.initRoutes()
	}

	router := handlers.LoggingHandler(os.Stdout, s.router)
	router.ServeHTTP(w, r)
}

func (s *Server) initRoutes() {
	s.router = mux.NewRouter()

	s.router.PathPrefix("/").Handler(s.handleFile()).Methods("GET")
}

func (s *Server) handleFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(s.baseDir, r.URL.Path)
		file, err := files.Info(s.baseDir, path)
		if err != nil {
			s.handleFileError(err, w, r)
			return
		}

		switch file.Type {
		case files.Directory:
			// TODO: Show list of files in a table.
			contents, err := files.List(s.baseDir, path)
			if err != nil {
				s.handleFileError(err, w, r)
				return
			}

			s.render(w, "directory", struct {
				Files []files.File
			}{
				Files: contents,
			})
		case files.SomeFile:
			// TODO: Show file contents.
			fmt.Fprint(w, file)
		}
	}
}

func (s *Server) handleFileError(err error, w http.ResponseWriter, r *http.Request) {
	http.Error(w, "File not found", 404)
}
