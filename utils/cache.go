package utils
/*
import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"context"
	"strconv"
)

var RedisClient *redis.Client

func InitRedis() error {
	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	db := os.Getenv("REDIS_DB")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}
	return nil
}*/
