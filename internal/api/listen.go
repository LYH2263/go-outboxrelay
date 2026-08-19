package api

import (
	"net"
	"net/http"
	"time"
)

// ListenAndServe 带超时的监听。
func ListenAndServe(addr string, h http.Handler) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return srv.Serve(ln)
}
