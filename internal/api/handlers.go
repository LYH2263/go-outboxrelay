package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/LYH2263/go-outboxrelay"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "closed": s.box.Closed()})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.box.Snapshot())
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	st := r.URL.Query().Get("status")
	var (
		evs []outboxrelay.Event
		err error
	)
	if st == "" || st == "pending" {
		evs, err = s.box.ListPending(limit)
	} else {
		evs, err = s.box.ListByStatus(outboxrelay.Status(st), limit)
	}
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, evs)
}

type appendBody struct {
	Topic   string            `json:"topic"`
	Payload json.RawMessage   `json:"payload"`
	Headers map[string]string `json:"headers"`
	URL     string            `json:"url"`
}

func (s *Server) handleAppend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	var body appendBody
	if err := readJSON(r, &body); err != nil {
		writeErr(w, err)
		return
	}
	payload := []byte(body.Payload)
	if body.URL != "" {
		id, err := s.box.AppendEvent(outboxrelay.Event{
			Topic:     body.Topic,
			Payload:   payload,
			Headers:   body.Headers,
			TargetURL: body.URL,
		})
		if err != nil {
			writeErr(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"id": id})
		return
	}
	id, err := s.box.Append(body.Topic, payload, body.Headers)
	if err != nil {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (s *Server) handleRelay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	max, _ := strconv.Atoi(r.URL.Query().Get("max"))
	if max <= 0 {
		max = 1
	}
	res, err := s.box.RelayContext(r.Context(), max)
	if err != nil && len(res) == 0 {
		writeErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"results": res, "err": errString(err)})
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func readJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	data, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
