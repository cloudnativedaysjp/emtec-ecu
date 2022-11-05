package db

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/model"
	redis "github.com/go-redis/redis/v8"
	"golang.org/x/xerrors"
)

type RedisClient struct {
	Client *redis.Client
}

const (
	RedisExpiration = 10 * time.Minute
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

func (rc *RedisClient) NextTalkNotificationJustWasSent(ctx context.Context, n model.NotificationOnDkTimetable) error {
	return rc.Client.Set(ctx, rc.nextTalkNotificationKey(n.Next().Id), true, RedisExpiration).Err()
}

func (rc *RedisClient) HasNextTalkNotificationAlreadyBeenSent(ctx context.Context, n model.NotificationOnDkTimetable) (bool, error) {
	result, err := rc.Client.Exists(ctx, rc.nextTalkNotificationKey(n.Next().Id)).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (rc *RedisClient) nextTalkNotificationKey(nextTalkId int32) string {
	return fmt.Sprintf("nextTalkNotificationAlreadySentFlag:%d", nextTalkId)
}

func (rc *RedisClient) MoveToNextSceneJustWasDone(ctx context.Context, next model.Talk) error {
	return rc.Client.Set(ctx, rc.moveToNextSceneKey(next.Id), true, RedisExpiration).Err()
}

func (rc *RedisClient) HasMoveToNextSceneBeenDone(ctx context.Context, next model.Talk) (bool, error) {
	result, err := rc.Client.Exists(ctx, rc.moveToNextSceneKey(next.Id)).Result()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func (rc *RedisClient) moveToNextSceneKey(nextTalkId int32) string {
	return fmt.Sprintf("moveToNextSceneFlag:%d", nextTalkId)
}
