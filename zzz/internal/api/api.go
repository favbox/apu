package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/negroni"
)

const TIMEOUT = 30 * time.Second

type ServerOption func(srv *http.Server)

// Start 使用默认参数启动一个带有平滑关闭功能的http服务器。
func Start(port string, handler http.Handler, options ...ServerOption) error {
	n := negroni.New()
	n.Use(negroni.HandlerFunc(MyMiddleware))
	n.UseHandler(handler) // 路由器排在最后

	srv := &http.Server{
		ReadTimeout:  TIMEOUT,
		WriteTimeout: TIMEOUT,
		Addr:         ":" + port,
		Handler:      n,
	}

	for _, o := range options {
		o(srv)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		log.Println("Stopping server...")
		err := srv.Shutdown(context.Background())
		if err != nil {
			panic(err)
		}
	}()

	log.Println(fmt.Sprintf("Starting server on port %s", port))
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// WithReadTimeout 配置 http.Server 读超时时长。
func WithReadTimeout(t time.Duration) ServerOption {
	return func(srv *http.Server) {
		srv.ReadTimeout = t
	}
}

// WithWriteTimeout 配置 http.Server 读超时时长。
func WithWriteTimeout(t time.Duration) ServerOption {
	return func(srv *http.Server) {
		srv.WriteTimeout = t
	}
}
