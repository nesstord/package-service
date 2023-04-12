package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"package-service/http/middlewares"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.CORS())

	r.GET("/healthcheck", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.Use(middlewares.Authorization())

	api := r.Group("/api")
	{
		api.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
		})
	}

	return r
}
