package main

import (
	"context"
	"embed"
	_ "embed"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
)

//go:embed assets/*
var Assets embed.FS

//go:embed static/index.html
var IndexHtml embed.FS

//go:embed static/index.html
var IndexByte []byte

//go:embed static/favicon.svg
var FaviconBytes []byte

func main() {
	h := server.Default()

	h.GET("/favicon.svg", func(ctx context.Context, c *app.RequestContext) {
		c.Header("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
		c.Data(http.StatusOK, "image/svg+xml", FaviconBytes)
		err := c.Flush()
		if err != nil {
			return
		}
	})

	h.GET("/", func(ctx context.Context, c *app.RequestContext) {
		c.Header("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		c.Data(http.StatusOK, "text/html", IndexByte)
		err := c.Flush()
		if err != nil {
			return
		}
	})

	h.GET("/assets/*filepath", func(ctx context.Context, c *app.RequestContext) {
		staticServer := http.FileServer(http.FS(Assets))
		req, err := adaptor.GetCompatRequest(&c.Request)
		if err != nil {
			return
		}
		staticServer.ServeHTTP(adaptor.GetCompatResponseWriter(&c.Response), req)
	})

	h.NoRoute(func(ctx context.Context, c *app.RequestContext) {
		c.Status(http.StatusOK)
		c.Header("Accept", "text/html")
		_, _ = c.WriteString("敬请期待")
		err := c.Flush()
		if err != nil {
			return
		}
	})

	h.Spin()
}
