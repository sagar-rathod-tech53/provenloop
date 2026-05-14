package config

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func ConnectRedis() {

	RDB = redis.NewClient(&redis.Options{
		Addr:     AppConfig.RedisAddr,
		Password: AppConfig.RedisPassword,
		DB:       AppConfig.RedisDB,
		PoolSize: 100,
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("✅ Redis Connected")
}
