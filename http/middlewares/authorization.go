package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, tokenExists := os.LookupEnv("AUTH_TOKEN")
		if !tokenExists {
			panic("Auth token not exists")
		}

		t := ""
		if len(c.Request.Header["Authorization"]) > 0 {
			t = c.Request.Header["Authorization"][0]
		}

		if t != token {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Next()
	}
}
