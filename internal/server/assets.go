package server

import (
	"embed"
	"net/http"
)

func (s *Server) handleAssets() http.Handler {
	return http.FileServer(http.FS(assetsFS))
}

//go:embed assets
var assetsFS embed.FS
