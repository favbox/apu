package gin

import (
	"context"
	"net/http"

	payment2 "apu/kitex_gen/payment"
	"apu/payment"
	"github.com/gin-gonic/gin"
)

func Handlers(s payment.UseCase) *gin.Engine {
	r := gin.Default()
	r.GET("/health", healthHandler)
	r.Handle("GET", "/", queryOrder(s))
	return r
}

func queryOrder(s payment.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &payment2.QueryOrderReq{
			OutOrderNo: c.Param("out_order_no"),
		}
		orderResp, err := s.QueryOrder(context.Background(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		c.JSON(http.StatusOK, orderResp)
	}
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "Gin App is healthy")
}
