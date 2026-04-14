package redisstore

import (
	"context"
	"fmt"
	"time"

	"github.com/AnastasiaDMW/auth-service/internal/store"
	"github.com/redis/go-redis/v9"
)

type RedisTokenStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisTokenStore(addr string, username string, password string) (*RedisTokenStore, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
	})

	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, store.ErrRedisConnection
	}

	return &RedisTokenStore{
		client: rdb,
		ctx:    ctx,
	}, nil
}

func (r *RedisTokenStore) key(token string) string {
	return "auth:refresh:" + token
}

func (r *RedisTokenStore) userKey(userID string) string {
	return fmt.Sprintf("user:%s:refresh_tokens", userID)
}

func (r *RedisTokenStore) SaveRefreshToken(token string, userID string, ttl time.Duration) error {
	key := r.key(token)

	pipe := r.client.TxPipeline()
	pipe.Set(r.ctx, key, userID, ttl)

	userKey := r.userKey(userID)
	pipe.SAdd(r.ctx, userKey, token)
	pipe.Expire(r.ctx, userKey, ttl)

	_, err := pipe.Exec(r.ctx)
	return err
}

func (r *RedisTokenStore) GetRefreshToken(token string) (string, error) {
	val, err := r.client.Get(r.ctx, r.key(token)).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisTokenStore) DeleteRefreshToken(token, userID string) error {
	key := r.key(token)
	userKey := r.userKey(userID)

	pipe := r.client.TxPipeline()

	pipe.Del(r.ctx, key)
	pipe.SRem(r.ctx, userKey, token)

	_, err := pipe.Exec(r.ctx)
	return err
}

func (r *RedisTokenStore) DeleteAllUserTokens(userID string) error {
	userKey := r.userKey(userID)

	tokens, err := r.client.SMembers(r.ctx, userKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	pipe := r.client.TxPipeline()

	for _, t := range tokens {
		if t != "" {
			pipe.Del(r.ctx, r.key(t))
		}
	}

	pipe.Del(r.ctx, userKey)

	_, err = pipe.Exec(r.ctx)
	return err
}

func (r *RedisTokenStore) HasUserTokens(userID string) (bool, error) {
	key := r.userKey(userID)

	count, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
