package helpers

import "github.com/gin-gonic/gin"

func UserIDFromContext(c *gin.Context) (uint, bool) {
	v, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	userID, ok := v.(uint)
	return userID, ok
}
