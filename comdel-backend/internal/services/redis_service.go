package services

import (
	"comdel-backend/internal/model"
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

type UserStore interface {
	SaveUser(user model.User) error;
	GetUser(userId string) (*model.User, error)
}

type RedisUserStore struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) RedisUserStore {
	return RedisUserStore{client: client}
}

func (rc *RedisUserStore) SaveUser(user model.User) error {
	log.Info("Redis:")
	log.Info(user.UserId)
	res, err := rc.client.HSet(context.Background(), model.RedisKey(user.UserId), user.RedisHashString()).Result()

	if err != nil {
		return err;
	}

	log.Info("Response : ", res);

	return nil
}

func (rc *RedisUserStore) GetUser(userId string) (*model.User, error) {
	var user model.User;
	err := rc.client.HGetAll(context.Background(), model.RedisKey(userId)).Scan(&user)

	if err != nil {
		return nil, err;
	}

	return &user, nil
}