package repository

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/x-frame/x-frame/internal/service"
	"github.com/redis/go-redis/v9"
)

type refreshTokenCache struct {
	rdb *redis.Client
}

// NewRefreshTokenCache creates a Redis-backed RefreshTokenCache.
func NewRefreshTokenCache(rdb *redis.Client) service.RefreshTokenCache {
	return &refreshTokenCache{rdb: rdb}
}

func (c *refreshTokenCache) tokenKey(hash string) string {
	return "rt:token:" + hash
}

func (c *refreshTokenCache) userSetKey(userID int64) string {
	return "rt:user:" + strconv.FormatInt(userID, 10)
}

func (c *refreshTokenCache) familySetKey(familyID string) string {
	return "rt:family:" + familyID
}

func (c *refreshTokenCache) StoreRefreshToken(ctx context.Context, tokenHash string, data *service.RefreshTokenData, ttl time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, c.tokenKey(tokenHash), b, ttl).Err()
}

func (c *refreshTokenCache) GetRefreshToken(ctx context.Context, tokenHash string) (*service.RefreshTokenData, error) {
	val, err := c.rdb.Get(ctx, c.tokenKey(tokenHash)).Result()
	if err == redis.Nil {
		return nil, service.ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}
	var data service.RefreshTokenData
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *refreshTokenCache) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	return c.rdb.Del(ctx, c.tokenKey(tokenHash)).Err()
}

func (c *refreshTokenCache) AddToUserTokenSet(ctx context.Context, userID int64, tokenHash string, ttl time.Duration) error {
	key := c.userSetKey(userID)
	pipe := c.rdb.Pipeline()
	pipe.SAdd(ctx, key, tokenHash)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	return err
}

func (c *refreshTokenCache) AddToFamilyTokenSet(ctx context.Context, familyID, tokenHash string, ttl time.Duration) error {
	key := c.familySetKey(familyID)
	pipe := c.rdb.Pipeline()
	pipe.SAdd(ctx, key, tokenHash)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	return err
}

func (c *refreshTokenCache) DeleteTokenFamily(ctx context.Context, familyID string) error {
	key := c.familySetKey(familyID)
	members, err := c.rdb.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	pipe := c.rdb.Pipeline()
	for _, m := range members {
		pipe.Del(ctx, c.tokenKey(m))
	}
	pipe.Del(ctx, key)
	_, err = pipe.Exec(ctx)
	return err
}

func (c *refreshTokenCache) DeleteUserRefreshTokens(ctx context.Context, userID int64) error {
	key := c.userSetKey(userID)
	members, err := c.rdb.SMembers(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	pipe := c.rdb.Pipeline()
	for _, m := range members {
		pipe.Del(ctx, c.tokenKey(m))
	}
	pipe.Del(ctx, key)
	_, err = pipe.Exec(ctx)
	return err
}
