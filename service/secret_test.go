package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewRedisSecretConnectivity(t *testing.T) {
	svc := NewRedisSecret("wp-atrd-task-database:6379")
	err, status := svc.Ping(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, status, "PONG")
}
