package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func staticHandler(dir string) http.Handler {
	fs := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && !strings.Contains(r.URL.Path, ".") {
			// SPA fallback
			index := filepath.Join(dir, "index.html")
			if _, err := os.Stat(index); err == nil {
				http.ServeFile(w, r, index)
				return
			}
		}
		fs.ServeHTTP(w, r)
	})
}
