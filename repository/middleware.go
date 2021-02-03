package repository

import (
	"os"

	"github.com/gin-gonic/gin"
)

// GinMiddleware is used to set the database in the current gin context.
func GinMiddleware(connString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn := NewByConnString(connString)

		c.Set("mongoRepoConn", conn)
		c.Set(string(ContextDBName), os.Getenv("DB_NAME"))

		c.Next()
	}
}
