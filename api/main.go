package api

import (
	// "context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pawmart/wp-atrd-task/service"
	"net/http"
)

type Api interface {
	Run(addr ...string) (err error)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type api struct {
	*gin.Engine
	secrets service.Secret
	// config  service.Config
}

// func (this *api) FetchSpecification(ctx context.Context) (err error, ok bool) {
// 	return nil, false
// }

func NewApi(secret service.Secret) Api {
	this := api{gin.Default(), secret}
	this.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	this.
		StaticFile("/swagger.yml", "./swagger/swagger.yml")
		// StaticFS("/swagger.yml", this.FetchSpecification)
	this.Group("/v1").
		POST("/secret", this.CreateSecret).
		GET("/secret/:hash", this.FetchSecret)

	return this
}
