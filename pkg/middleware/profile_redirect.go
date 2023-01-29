package middleware

import (
	"beneburg/pkg/database/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ProfileRedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("currentUser").(*model.User)
		if user.Status == model.UserStatusNew {
			c.Redirect(http.StatusFound, "/profile")
			return
		}
	}
}
