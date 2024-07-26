package fiber

import (
	"context"
	"net/http"

	payment2 "apu/kitex_gen/payment"
	"apu/payment"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func Handlers(s payment.UseCase) http.HandlerFunc {
	app := fiber.New()
	app.Get("/health", healthHandler)
	app.Get("/", getBooks(s))
	return adaptor.FiberApp(app)
}

func getBooks(s payment.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := &payment2.QueryOrderReq{
			OutOrderNo: c.Query("out_order_no"),
		}
		orderResp, err := s.QueryOrder(context.Background(), req)
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(http.StatusOK).JSON(orderResp)
	}
}

func healthHandler(c *fiber.Ctx) error {
	return c.SendString("Fiber App is healthy")
}
