package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrMiss = redis.Nil

type Client struct {
	client *redis.Client
}

func New(redisURL string) (*Client, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	return &Client{client: redis.NewClient(options)}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Client) Get(ctx context.Context, key string, destination any) error {
	value, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(value, destination); err != nil {
		return errors.Join(errors.New("decode cached value"), err)
	}
	return nil
}

func (c *Client) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Join(errors.New("encode cached value"), err)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *Client) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}
