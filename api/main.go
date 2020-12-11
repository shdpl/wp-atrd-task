package api

import (
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
}

func NewApi() Api {
	this := api{gin.Default(), &service.RedisSecret{}}
	this.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	this.
		StaticFile("/swagger.yml", "./swagger/swagger.yml")
	this.Group("/v1").
		POST("/secret", this.AddSecret).
		GET("/secret/:hash", this.GetSecret)

	return this
}
