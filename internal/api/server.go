package api

import (
	"net/http"

	"github.com/LYH2263/go-outboxrelay"
)

// Options HTTP 服务选项。
type Options struct {
	WebDir    string
	AllowCORS bool
}

// Server 管理 API。
type Server struct {
	box  *outboxrelay.Outbox
	opts Options
	mux  *http.ServeMux
}

// New 构造。
func New(box *outboxrelay.Outbox, opts Options) *Server {
	s := &Server{box: box, opts: opts, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.opts.AllowCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	s.mux.ServeHTTP(w, r)
}
