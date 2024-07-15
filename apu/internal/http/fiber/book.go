package fiber

import (
	"net/http"

	"apu/book"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

func Handlers(s book.UseCase) http.HandlerFunc {
	app := fiber.New()
	app.Get("/health", healthHandler)
	app.Get("/", getBooks(s))
	return adaptor.FiberApp(app)
}

func getBooks(s book.UseCase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		books, err := s.GetAll()
		if err != nil {
			return c.Status(http.StatusInternalServerError).SendString(err.Error())
		}
		return c.Status(http.StatusOK).JSON(books)
	}
}

func healthHandler(c *fiber.Ctx) error {
	return c.SendString("Fiber App is healthy")
}
