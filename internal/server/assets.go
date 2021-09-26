package server

import (
	"embed"
	"fmt"
	"hash/crc32"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"sync"
)

func (s *Server) handleAssets() http.Handler {
	fileServer := http.FileServer(http.FS(HashedAssetsFS(func(p string) (fs.File, error) {
		return s.assets.openAsset(p)
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
	lock                 sync.Mutex
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
		newFingerprint, err := fp.addAsset(p)
		if err != nil {
			return "", err
		}

		fingerprintedPath = newFingerprint
	}

	return path.Join(assetPathPrefix, fingerprintedPath), nil
}

func (fp *assetFingerprinter) openAsset(p string) (fs.File, error) {
	return assetsFS.Open(assetFsPrefix + fp.fingerprintedToAsset[p])
}

func (fp *assetFingerprinter) addAsset(p string) (string, error) {
	fp.lock.Lock()
	defer fp.lock.Unlock()

	content, err := assetsFS.ReadFile(assetFsPrefix + p)
	if err != nil {
		return "", err
	}

	ext := path.Ext(p)
	fingerprint := fmt.Sprint(crc32.ChecksumIEEE(content))
	fingerprintedPath := strings.TrimSuffix(p, ext) + "-" + fingerprint + ext

	fp.assetToFingerprinted[p] = fingerprintedPath
	fp.fingerprintedToAsset[fingerprintedPath] = p

	return fingerprintedPath, nil
}

type HashedAssetsFS func(string) (fs.File, error)

func (f HashedAssetsFS) Open(name string) (fs.File, error) {
	return f(name)
}
