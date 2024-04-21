package producers

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisProducer interface {
	Set(key string, value interface{}, expiration time.Duration)
	Get(key string) string
	Incr(key string)
	Decr(key string)
}

type redisProducer struct {
	client *redis.Client
}

func NewRedisProducer(host string, port string) RedisProducer {
	addr := fmt.Sprintf("%s:%s", host, port)
	rs := redisProducer{}
	rs.client = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	return &rs
}

func (s *redisProducer) Set(key string, value interface{}, expiration time.Duration) {
	err := s.client.Set(context.Background(), key, value, expiration).Err()
	if err != nil {
		log.Println("redisProducer set error", err)
	}
}

func (s *redisProducer) Get(key string) string {
	val, err := s.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		fmt.Println("redisProducer", key, "does not exist")
	} else if err != nil {
		log.Println("redisProducer get error", err)
	}
	return val
}

func (s *redisProducer) Incr(key string) {
	err := s.client.Incr(context.Background(), key).Err()
	if err != nil {
		log.Println("redisProducer incr error", err)
	}
}

func (s *redisProducer) Decr(key string) {
	err := s.client.Decr(context.Background(), key).Err()
	if err != nil {
		log.Println("redisProducer decr error", err)
	}
}
