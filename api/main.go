package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pawmart/wp-atrd-task/api/secret"
	"net/http"
)

type Api interface {
	Run(addr ...string) (err error)
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

func NewApi() Api {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))

	r.
		StaticFile("/swagger.yml", "./swagger/swagger.yml").
		POST("/secret", secret.Add).
		GET("/secret", secret.Get)

	return r
}
