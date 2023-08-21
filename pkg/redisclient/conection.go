package pkg

import (
	"context"
	"strconv"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/arvan/qoute/config"
)

var (
	Cfg              *config.Config
	Ctx              = context.Background()
	zlogger          = log.With().Str("service", "Arvan-Qoute").Logger()
	RedisClient      *redis.Ring
	RedisCacheQoute  *cache.Cache
	RedisCacheCheck  *cache.Cache
	RedisCacheVolume *cache.Cache
)

func Connect(ctx context.Context, cfg *config.Config) *redis.Ring {

	Cfg = cfg
	var host = cfg.Redis.Ip
	var port = cfg.Redis.Port
	Ctx = ctx
	client := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"shard1": host + ":" + strconv.Itoa(port),
		},
	})

	res, err := client.Ping(Ctx).Result()

	if err != nil {
		zlogger.Info().Msgf("Redis connection failed")
		return nil
	}

	RedisCacheQoute = cache.New(&cache.Options{
		Redis: client,
	})

	RedisCacheCheck = cache.New(&cache.Options{
		Redis: client,
	})

	RedisCacheVolume = cache.New(&cache.Options{
		Redis: client,
	})

	zlogger.Info().Msgf("Redis Server Connected ....... ")
	zlogger.Info().Msgf("Ping to Redis Server :" + res)

	if cfg.Redis.Flush {
		client.FlushAll(Ctx)
		zlogger.Info().Msgf("Redis DB Flushed....")
	}

	RedisClient = client
	return client
}
