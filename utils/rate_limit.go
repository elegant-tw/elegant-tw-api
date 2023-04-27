package utils

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	limiter "github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	mredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimitInit(router *gin.Engine, cfg *Config) {
	rate, err := limiter.NewRateFromFormatted(cfg.RateLimitFormatted)

	if err != nil {
		logrus.Fatal(err)
	}

	var store limiter.Store

	if cfg.RateLimitStoreRedis {
		logrus.Info("Rate limit stores in redis.")
		client := redisInit(cfg)
		mstore, err := mredis.NewStoreWithOptions(client, limiter.StoreOptions{
			Prefix: cfg.RateLimitRedisPrefix,
		})
		if err != nil {
			logrus.Panic(err)
		}
		store = mstore
	} else {
		logrus.Info("Rate limit stores in memory.")
		store = memory.NewStore()
	}
	instance := limiter.New(store, rate)
	middleware := mgin.NewMiddleware(instance)
	router.Use(middleware)
	logrus.Infof("Request limit is %d reqs / %s", rate.Limit, rate.Period)
}

func redisInit(cfg *Config) *redis.Client {
	logrus.Info("Connecting redis...")
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RateLimitRedisHost + ":" + cfg.RateLimitRedisPort,
	})
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Redis message: %+v", pong)
	logrus.Infof("Connected redis.")
	return rdb
}
