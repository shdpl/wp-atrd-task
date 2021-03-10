package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRedisSecretConnectivity(t *testing.T) {
	var config Config
	err := config.Init()
	if err != nil {
		panic(err)
	}

	err = config.Unmarshal("../config.test")
	if err != nil {
		panic(err)
	}

	svc := NewRedisSecret(config.Redis)
	err, status := svc.Ping(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, "PONG", status)
}
