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
	fileServer := http.FileServer(http.FS(HashedAssetsFS(func(p string) (fs.File, error) {
		return assetsFS.Open(assetFsPrefix + s.assets.fingerprintedToAsset[p])
	})))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "public, max-age=604800, immutable")

		fileServer.ServeHTTP(w, r)
	})
}

//go:embed assets
var assetsFS embed.FS

const (
	assetPathPrefix = "/filez/assets"
	assetFsPrefix   = "assets/"
)

type assetFingerprinter struct {
	assetToFingerprinted map[string]string
	fingerprintedToAsset map[string]string
}

func newAssetFingerprinter() *assetFingerprinter {
	return &assetFingerprinter{
		assetToFingerprinted: map[string]string{},
		fingerprintedToAsset: map[string]string{},
	}
}

func (fp *assetFingerprinter) assetPath(p string) (string, error) {
	fingerprintedPath, found := fp.assetToFingerprinted[p]
	if !found {
		content, err := assetsFS.ReadFile(assetFsPrefix + p)
		if err != nil {
			return "", err
		}

		newFingerprint := addFingerprint(p, fmt.Sprint(crc32.ChecksumIEEE(content)))

		fingerprintedPath, fp.assetToFingerprinted[p] = newFingerprint, newFingerprint
		fp.fingerprintedToAsset[fingerprintedPath] = p
	}

	return path.Join(assetPathPrefix, fingerprintedPath), nil
}

func addFingerprint(p string, fingerprint string) string {
	ext := path.Ext(p)

	return strings.TrimSuffix(p, ext) + "-" + fingerprint + ext
}

type HashedAssetsFS func(string) (fs.File, error)

func (f HashedAssetsFS) Open(name string) (fs.File, error) {
	return f(name)
}
