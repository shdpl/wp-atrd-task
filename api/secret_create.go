package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/pawmart/wp-atrd-task/models"
	"net/http"
	"time"
)

func (this *api) CreateSecret(c *gin.Context) {
	var err error
	var params struct {
		Secret           string `form:"secret" binding:"required"`
		ExpireAfterViews int32  `form:"expireAfterViews" binding:"required"`
		ExpireAfter      int32  `form:"expireAfter" binding:"required"`
	}
	err = c.ShouldBind(&params)
	if err != nil {
		c.AbortWithError(http.StatusMethodNotAllowed, err)
		return
	}

	now := time.Now()

	secret := models.Secret{
		SecretText:     params.Secret,
		RemainingViews: params.ExpireAfterViews,
		CreatedAt:      strfmt.DateTime(now),
		ExpiresAt:      strfmt.DateTime(now.Add(time.Duration(params.ExpireAfter) * time.Minute)),
	}

	err = this.secrets.Create(c, &secret)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Negotiate(http.StatusOK, gin.Negotiate{
		Offered: []string{gin.MIMEJSON, gin.MIMEXML},
		Data:    secret,
	})
}
