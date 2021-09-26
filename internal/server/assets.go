package server

import (
	"embed"
	"fmt"
	"hash/crc32"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

func (s *Server) handleAssets() http.Handler {
	return http.FileServer(http.FS(HashedAssetsHandler(handleAsset)))
}

func handleAsset(fingerprinted string) (fs.File, error) {
	return assetsFS.Open(assetFsPrefix + fingerprintedToAsset[fingerprinted])
	// TODO: cache headers
}

//go:embed assets
var assetsFS embed.FS

const (
	assetPathPrefix = "/filez/assets"
	assetFsPrefix   = "assets/"
)

var assetToFingerprinted = map[string]string{}
var fingerprintedToAsset = map[string]string{}

func assetPath(p string) (string, error) {
	fingerprintedPath, found := assetToFingerprinted[p]
	if !found {
		content, err := assetsFS.ReadFile(assetFsPrefix + p)
		if err != nil {
			return "", err
		}

		newFingerprint := addFingerprint(p, fmt.Sprint(crc32.ChecksumIEEE(content)))

		fingerprintedPath, assetToFingerprinted[p] = newFingerprint, newFingerprint
		fingerprintedToAsset[fingerprintedPath] = p
	}

	return path.Join(assetPathPrefix, fingerprintedPath), nil
}

func addFingerprint(p string, fingerprint string) string {
	ext := path.Ext(p)

	return strings.TrimSuffix(p, ext) + "-" + fingerprint + ext
}

type HashedAssetsHandler func(string) (fs.File, error)

func (f HashedAssetsHandler) Open(name string) (fs.File, error) {
	return f(name)
}
