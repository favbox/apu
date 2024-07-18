package hertz

import (
	"context"
	"net/http"

	payment2 "apu/kitex_gen/payment"
	"apu/payment"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func Handlers(s payment.UseCase) {
	h := server.Default(server.WithHostPorts(":8080"))
	h.GET("/health", healthHandler)
	h.GET("/", getBooks(s))
	h.Spin()
}

func getBooks(s payment.UseCase) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		req := &payment2.QueryOrderReq{
			OutOrderNo: c.Param("out_order_no"),
		}
		orderResp, err := s.QueryOrder(context.Background(), req)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
		}
		c.JSON(http.StatusOK, orderResp)
	}
}

func healthHandler(_ context.Context, c *app.RequestContext) {
	c.String(http.StatusOK, "Hertz App is healthy")
}
