package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that sets CORS headers so frontend (e.g. localhost:3000) can call this API.
func CORS(allowOrigins []string) gin.HandlerFunc {
	originSet := make(map[string]bool)
	for _, o := range allowOrigins {
		originSet[o] = true
	}
	if len(originSet) == 0 {
		originSet["http://localhost:3000"] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		// Set Allow-Origin: allow list, or any localhost for development
		if origin != "" {
			if originSet[origin] || originSet["*"] {
				c.Header("Access-Control-Allow-Origin", origin)
			} else if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "http://127.0.0.1:") {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
