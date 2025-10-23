package resource

import (
	"context"
	"fmt"
	"productfc/config"
	"productfc/infrastructure/log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(cfg config.RedisConfig) *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})

	ctx := context.Background()
	res, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to connect Redis")
	}

	log.Logger.Info().Msgf("Redis connected: %s", res)

	return RedisClient
}
