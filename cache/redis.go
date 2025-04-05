package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",               
		DB:       0,                
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Redis connection failed: %v", err))
	}
	fmt.Println("Connected to Redis")
}
