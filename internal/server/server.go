package server

import (
	"errors"
	"io/fs"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mthadley/filez/internal/files"
)

type Server struct {
	base   fs.FS
	router *mux.Router
	assets *assetFingerprinter
}

func NewServer(base fs.FS) *Server {
	return &Server{
		base:   base,
		assets: newAssetFingerprinter(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.router == nil {
		s.initRoutes()
	}

	s.router.ServeHTTP(w, r)
}

func (s *Server) initRoutes() {
	s.router = mux.NewRouter()

	filez := s.router.PathPrefix("/filez").Subrouter()
	filez.PathPrefix("/assets").
		Handler(http.StripPrefix(assetPathPrefix+"/", s.handleAssets())).
		Methods("GET")
	filez.PathPrefix("/raw").
		Handler(http.StripPrefix("/filez/raw", s.handleFileRaw())).
		Methods("GET")

	s.router.PathPrefix("/").Handler(s.handleFile()).Methods("GET")

	s.router.Use(
		handlers.CompressHandler,
		func(h http.Handler) http.Handler { return handlers.LoggingHandler(os.Stdout, h) },
	)
}

func (s *Server) handleFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		file, err := files.Info(s.base, path)

		if err != nil {
			s.handleFileError(err, w)
			return
		}

		switch file.Type {
		case files.Directory:
			contents, err := files.List(s.base, path)
			if err != nil {
				s.handleFileError(err, w)
				return
			}

			s.render(w, "directory", struct {
				File  files.File
				Files []files.File
			}{
				File:  file,
				Files: contents,
			})
		case files.SomeFile, files.Symlink:
			s.render(w, "file", file)
		case files.SpecialFile:
			s.handleError(
				w,
				errors.New("Cannot open this kind of file."),
				http.StatusForbidden,
			)
		}
	}
}

func (s *Server) handleFileRaw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, err := files.Info(s.base, r.URL.Path)
		if err != nil {
			s.handleFileError(err, w)
			return
		}
		steam, err := file.OpenSteam()
		if err != nil {
			s.handleFileError(err, w)
			return
		}
		defer steam.Close()

		http.ServeContent(w, r, file.Name, file.ModTime, steam)
	}
}

func (s *Server) handleFileError(err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, fs.ErrNotExist):
		err = errors.New("File not found.")
		status = http.StatusNotFound
	case errors.Is(err, fs.ErrPermission):
		err = errors.New("No permission to view this file.")
		status = http.StatusNotFound
	default:
		err = errors.New("Unable to view file or page.")
	}

	s.handleError(w, err, status)
}

func (s *Server) handleError(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)
	s.render(w, "error", err)
}
