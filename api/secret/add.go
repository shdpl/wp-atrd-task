package secret

import (
	"github.com/gin-gonic/gin"
)

func Add(c *gin.Context) {
	c.Status(501)
	// c.JSON(200, gin.H{
	// 	"createdAt":      "2017-07-21T17:32:28Z",
	// 	"expiresAt":      "2017-07-21T18:32:28Z",
	// 	"hash":           "b75ce598-f349-4c61-9246-2053e230187d",
	// 	"remainingViews": 0,
	// 	"secretText":     "secret",
	// })
}
