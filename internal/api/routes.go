package api

import "net/http"

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/events", s.handleEvents)
	s.mux.HandleFunc("/api/events/append", s.handleAppend)
	s.mux.HandleFunc("/api/relay", s.handleRelay)
	s.mux.HandleFunc("/api/stats", s.handleStats)
	if s.opts.WebDir != "" {
		s.mux.Handle("/", staticHandler(s.opts.WebDir))
	} else {
		s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	}
}
