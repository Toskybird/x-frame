package repository

import (
	"database/sql"
	"errors"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/x-frame/x-frame/ent"
	"github.com/x-frame/x-frame/internal/config"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is the Wire provider set for scaffold repositories
var ProviderSet = wire.NewSet(
	NewUserRepository,
	NewAnnouncementRepository,
	NewAnnouncementReadRepository,
	NewIdempotencyRepository,
	NewSettingRepository,
	NewUserAttributeDefinitionRepository,
	NewUserAttributeValueRepository,

	// Cache implementations
	NewEmailCache,
	NewTotpCache,
	NewRefreshTokenCache,

	// Encryptors
	NewAESEncryptor,

	// Backup infrastructure
	NewPgDumper,
	NewS3BackupStoreFactory,

	ProvideEnt,
	ProvideSQLDB,
	ProvideRedis,
)

// ProvideEnt 为依赖注入提供 Ent 客户端。
func ProvideEnt(cfg *config.Config) (*ent.Client, error) {
	client, _, err := InitEnt(cfg)
	return client, err
}

// ProvideSQLDB 从 Ent 客户端提取底层的 *sql.DB 连接。
func ProvideSQLDB(client *ent.Client) (*sql.DB, error) {
	if client == nil {
		return nil, errors.New("nil ent client")
	}
	drv, ok := client.Driver().(*entsql.Driver)
	if !ok {
		return nil, errors.New("ent driver does not expose *sql.DB")
	}
	return drv.DB(), nil
}

// ProvideRedis 为依赖注入提供 Redis 客户端。
func ProvideRedis(cfg *config.Config) *redis.Client {
	return InitRedis(cfg)
}
