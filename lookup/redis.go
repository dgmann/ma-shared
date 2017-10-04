package lookup

import (
	"github.com/go-redis/redis"
)

func newRedisClient(host string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

type RedisStorage struct {
	redis *redis.Client
}

func NewRedisStorage(host string) *RedisStorage {
	return &RedisStorage{redis:newRedisClient(host)}
}

func(storage *RedisStorage) Exists(keys ...string) bool {
	count, err := storage.redis.Exists(keys...).Result()
	if err != nil {
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
