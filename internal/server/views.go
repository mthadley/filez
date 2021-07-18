package server

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

const (
	layoutViewFile = "views/layout.go.tmpl"
	layoutViewName = "layout"
)

func (s *Server) render(w http.ResponseWriter, name string, data interface{}) {
	views, err := template.ParseFS(viewFS, layoutViewFile, "views/"+name+".go.tmpl")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), 500)
		return
	}

	if err := views.ExecuteTemplate(w, layoutViewName, data); err != nil {
		http.Error(w, fmt.Sprintf("Error executing view: %v", err), 500)
	}
}

//go:embed views
var viewFS embed.FS