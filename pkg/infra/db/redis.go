package db

import (
	"context"
	"strconv"
	"time"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/model"
	redis "github.com/go-redis/redis/v8"
	"golang.org/x/xerrors"
)

type RedisClient struct {
	Client *redis.Client
}

const (
	RedisExpiration         = 10 * time.Minute
	NextTalkNotificationKey = "nextTalkNotificationAlreadySentFlag:"
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

func (rc *RedisClient) SetNextTalkNotification(ctx context.Context, n model.NotificationOnDkTimetable) error {
	nextTalkId := int(n.Next().Id)
	return rc.Client.Set(ctx, NextTalkNotificationKey+strconv.Itoa(nextTalkId), true, RedisExpiration).Err()
}

func (rc *RedisClient) HasNextTalkNotificationAlreadyBeenSent(ctx context.Context, n model.NotificationOnDkTimetable) (bool, error) {
	nextTalkId := int(n.Next().Id)
	result, err := rc.Client.Exists(ctx, NextTalkNotificationKey+strconv.Itoa(nextTalkId)).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
