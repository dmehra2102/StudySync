package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dmehra2102/StudySync/internal/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	client *redis.Client
	ctx    context.Context
}

func New(cfg config.RedisConfig) *Client {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		opt = &redis.Options{
			Addr:     cfg.URL,
			Password: "",
			DB:       0,
		}
	}

	client := redis.NewClient(opt)

	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	return &Client{
		client: client,
		ctx:    ctx,
	}
}

func (c *Client) Publish(channel string, message any) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.client.Publish(c.ctx, channel, jsonMessage).Err()
}

func (c *Client) Subscribe(channel string) *redis.PubSub {
	return c.client.Subscribe(c.ctx, channel)
}

func (c *Client) Set(key string, value any, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(c.ctx, key, jsonValue, expiration).Err()
}

func (c *Client) Get(key string, dest any) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (c *Client) Close() error {
	return c.client.Close()
}
