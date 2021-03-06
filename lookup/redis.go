package lookup

import (
	"github.com/go-redis/redis"
	"fmt"
)

func newRedisClient(host string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     host + ":6379",
		Password: password, // no password set
		DB:       0,  // use default DB
	})
}

type RedisStorage struct {
	redis *redis.Client
}

func NewRedisStorage(host string, password string) *RedisStorage {
	return &RedisStorage{redis:newRedisClient(host, password)}
}

func(storage *RedisStorage) Exists(keys ...string) bool {
	count, err := storage.redis.Exists(keys...).Result()
	if err != nil {
		fmt.Printf("Redis error: %v", err)
		return false
	}
	return int(count) == len(keys)
}

func(storage *RedisStorage) Close() {
	storage.redis.Close()
}

func(storage *RedisStorage) Add(key string) {
	storage.redis.Set(key, true, 0)
}
