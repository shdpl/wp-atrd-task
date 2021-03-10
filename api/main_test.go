package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pawmart/wp-atrd-task/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newApi() Api {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
	}))
	return router
}

func TestSwaggerCors(t *testing.T) {
	req, err := http.NewRequest("GET", "/swagger.yml", nil)
	req.Header.Set("Origin", "http://127.0.0.1:8081")
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	var config service.Config
	err = config.Init()
	if err != nil {
		panic(err)
	}

	err = config.Unmarshal("../config.test")
	if err != nil {
		panic(err)
	}

	NewApi(
		service.NewRedisSecret(config.Redis),
	).ServeHTTP(w, req)
	assert.Equal(t, "*", w.Result().Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, 200, w.Result().StatusCode)
}
