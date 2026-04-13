package redis

import (
	"context"
	"log"
	"strconv"

	"github.com/Aryan-Gupta4460/letstalk/pkg/config"
	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()

func InitRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Cache.Host + ":" + strconv.Itoa(cfg.Cache.Port),
	})

	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatal("Redis connection failed:", err)
	}

	log.Println("Connected to Redis")

	return rdb
}
