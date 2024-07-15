package gin

import (
	"net/http"

	"apu/book"
	"github.com/gin-gonic/gin"
)

func Handlers(s book.UseCase) *gin.Engine {
	r := gin.Default()
	r.GET("/health", healthHandler)
	r.Handle("GET", "/", getBooks(s))
	return r
}

func getBooks(s book.UseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := s.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		c.JSON(http.StatusOK, books)
	}
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "Gin App is healthy")
}
