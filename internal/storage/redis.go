package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

// InitRedis инициализирует подключение к Redis
func InitRedis(addr string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: addr, // например "localhost:6379"
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("не удалось подключиться к Redis: %v", err))
	}
}

// SaveToken сохраняет токен в Redis с Telegram ID как ключом
func SaveToken(telegramID int64, token string, ttl time.Duration) error {
	return rdb.Set(ctx, fmt.Sprintf("token:%d", telegramID), token, ttl).Err()
}

// GetToken получает токен по Telegram ID
func GetToken(telegramID int64) (string, error) {
	return rdb.Get(ctx, fmt.Sprintf("token:%d", telegramID)).Result()
}
