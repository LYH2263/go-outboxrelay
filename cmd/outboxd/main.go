package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LYH2263/go-outboxrelay"
	"github.com/LYH2263/go-outboxrelay/internal/api"
)

func main() {
	addr := flag.String("addr", ":8091", "HTTP 监听地址")
	web := flag.String("web", "web", "静态管理页目录")
	target := flag.String("target", "http://127.0.0.1:9999/hook", "默认下游 URL")
	persist := flag.String("persist", "", "事件快照 JSON 路径（可选）")
	timeout := flag.Duration("http-timeout", 10*time.Second, "投递超时")
	maxTry := flag.Int("max-attempts", 5, "最大尝试次数")
	flag.Parse()

	opts := []outboxrelay.Option{
		outboxrelay.WithTargetURL(*target),
		outboxrelay.WithHTTPTimeout(*timeout),
		outboxrelay.WithMaxAttempts(*maxTry),
		outboxrelay.WithAllowHTTP(true),
	}
	if *persist != "" {
		opts = append(opts, outboxrelay.WithPersistPath(*persist))
	}

	box := outboxrelay.New(opts...)
	defer box.Close()

	srv := api.New(box, api.Options{WebDir: *web, AllowCORS: true})
	httpSrv := &http.Server{
		Addr:              *addr,
		Handler:           srv,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("outboxd 监听 %s", *addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	_ = httpSrv.Close()
}
