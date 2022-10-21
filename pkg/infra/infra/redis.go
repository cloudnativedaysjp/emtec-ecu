package infra

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/xerrors"
)

type RedisClient struct {
	Client *redis.Client
}

const (
	RedisExpiration                 = 10 * time.Minute
	NextTalkNotificationKey         = "nextTalkNotificationAlreadySentFlag"
	NextTalkNotificationAlreadySent = true
)

func NewRedisClient(addr string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	if err := client.Ping(context.TODO()).Err(); err != nil {
		return nil, xerrors.Errorf("fail to connect to redis. message: %w", err)
	}
	return &RedisClient{
		Client: client,
	}, nil
}
