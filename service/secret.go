package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pawmart/wp-atrd-task/models"
	"time"
)

type Secret interface {
	Create(context.Context, *models.Secret) error
	FetchByHash(context.Context, *models.Secret) (error, bool)
}

type RedisSecret struct {
	client *redis.Client
}

//TODO: SetLogger
func NewRedisSecret(address string) *RedisSecret {
	this := &RedisSecret{}
	this.client = redis.NewClient(&redis.Options{Addr: address})
	return this
}

func (this *RedisSecret) Create(ctx context.Context, secret *models.Secret) error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	secret.Hash = uuid.String()
	createdAt := fmt.Sprintf("secret_%s_createdAt", uuid)
	expiresAt := fmt.Sprintf("secret_%s_expiresAt", uuid)
	remainingViews := fmt.Sprintf("secret_%s_remainingViews", uuid)
	text := fmt.Sprintf("secret_%s_secretText", uuid)
	pipe := this.client.TxPipeline()
	err = pipe.MSet(
		ctx,
		map[string]interface{}{
			createdAt:      secret.CreatedAt.String(),
			expiresAt:      secret.ExpiresAt.String(),
			remainingViews: secret.RemainingViews,
			text:           secret.SecretText,
		},
	).Err()
	if err != nil {
		return err
	}
	err = pipe.ExpireAt(ctx, createdAt, time.Time(secret.ExpiresAt)).Err()
	if err != nil {
		return err
	}
	err = pipe.ExpireAt(ctx, expiresAt, time.Time(secret.ExpiresAt)).Err()
	if err != nil {
		return err
	}
	err = pipe.ExpireAt(ctx, remainingViews, time.Time(secret.ExpiresAt)).Err()
	if err != nil {
		return err
	}
	err = pipe.ExpireAt(ctx, text, time.Time(secret.ExpiresAt)).Err()
	if err != nil {
		return err
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (this *RedisSecret) FetchByHash(ctx context.Context, secret *models.Secret) (err error, ok bool) {
	createdAtKey := fmt.Sprintf("secret_%s_createdAt", secret.Hash)
	expiresAtKey := fmt.Sprintf("secret_%s_expiresAt", secret.Hash)
	remainingViewsKey := fmt.Sprintf("secret_%s_remainingViews", secret.Hash)
	secretTextKey := fmt.Sprintf("secret_%s_secretText", secret.Hash)
	err = this.client.Watch(
		ctx,
		func(tx *redis.Tx) error {
			remainingViews, err := tx.Decr(ctx, remainingViewsKey).Result()
			if err != nil {
				return err
			}
			values, err := tx.MGet(ctx, createdAtKey, expiresAtKey, secretTextKey).Result()
			if err != nil {
				return err
			}
			createdAt := values[0]
			expiresAt := values[1]
			secretText := values[2]
			notExists := createdAt == nil || expiresAt == nil || secretText == nil
			if notExists {
				err = tx.Del(ctx, remainingViewsKey).Err()
				if err != nil {
					return err
				}
				ok = false
				return nil
			} else {
				if remainingViews <= 0 {
					err = tx.Del(ctx, createdAtKey, expiresAtKey, remainingViewsKey, secretTextKey).Err()
					if err != nil {
						return err
					}
				}
				err = secret.CreatedAt.Scan(createdAt)
				if err != nil {
					return err
				}
				err = secret.ExpiresAt.Scan(expiresAt)
				if err != nil {
					return err
				}
				secret.SecretText = secretText.(string)
				secret.RemainingViews = int32(remainingViews)
				ok = true
				return nil
			}
		},
		remainingViewsKey,
	)
	if err != nil {
		return err, ok
	}
	return nil, ok
}

func (this *RedisSecret) Ping() (err error, pong string) {
	pong, err = this.client.Ping(context.TODO()).Result()
	return
}
