package service

import (
	"context"
	"github.com/pawmart/wp-atrd-task/models"
)

type Secret interface {
	Put(context.Context, *models.Secret) error
}

type RedisSecret struct {
}

func (this *RedisSecret) Put(ctx context.Context, secret *models.Secret) error {
	// secret.Hash = "b75ce598-f349-4c61-9246-2053e230187d"
	return nil
}
