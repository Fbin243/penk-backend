package redisrepo

import "github.com/go-redis/redis/v8"

type ReminderCache struct {
	client *redis.Client
}

func NewReminderCache(client *redis.Client) *ReminderCache {
	return &ReminderCache{client: client}
}
