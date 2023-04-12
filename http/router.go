package http

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"package-service/http/controllers"
	"package-service/http/middlewares"
	"package-service/http/validators"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("gtin", validators.Gtin)
		v.RegisterValidation("sgtin", validators.Sgtin)
		v.RegisterValidation("sscc", validators.Sscc)
	}

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

		api.POST("/aggregate", controllers.Aggregate)
	}

	return r
}
