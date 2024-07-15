package hertz

import (
	"context"
	"net/http"

	"apu/book"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Handlers(s book.UseCase) {
	h := server.Default(server.WithHostPorts(":8080"))
	h.GET("/health", healthHandler)
	h.GET("/", getBooks(s))
	h.Spin()
}

func getBooks(s book.UseCase) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		books, err := s.GetAll()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.JSON(http.StatusOK, books)
	}
}

func healthHandler(_ context.Context, c *app.RequestContext) {
	c.String(http.StatusOK, "Hertz App is healthy")
}
