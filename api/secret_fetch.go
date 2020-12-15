package api

import (
	"github.com/gin-gonic/gin"
	"github.com/pawmart/wp-atrd-task/models"
	"net/http"
)

func (this *api) FetchSecret(c *gin.Context) {
	hash := c.Param("hash")

	secret := models.Secret{
		Hash: hash,
	}

	err, ok := this.secrets.FetchByHash(c, &secret)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	if ok {
		c.Negotiate(http.StatusOK, gin.Negotiate{
			Offered: []string{gin.MIMEJSON, gin.MIMEXML},
			Data:    secret,
		})
	} else {
		c.Negotiate(http.StatusNotFound, gin.Negotiate{
			Offered: []string{gin.MIMEJSON, gin.MIMEXML},
			Data:    gin.H{},
		})
	}
}
