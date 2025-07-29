package services

import (
	"comdel-backend/internal/model"
	"context"

	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

type UserStore interface {
	SaveUser(user model.User) error;
	GetUserAndVideo(userId string) (*model.User, []string, error)
	SaveVideoId(videoId string, userId string) error
	GetVideoIds(userId string) ([]string, error)
	IsCacheMiss(err error) bool
	IsCacheHit(err error) bool
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
	res, err := rc.client.HSet(context.Background(), model.RedisUserKey(user.UserId), user.RedisHashString()).Result()

	if err != nil {
		return err;
	}

	log.Info("Response : ", res);

	return nil
}

func (rc *RedisUserStore) GetUserAndVideo(userId string) (*model.User, []string, error) {
	var user model.User;
	err := rc.client.HGetAll(context.Background(), model.RedisUserKey(userId)).Scan(&user)

	if err != nil {
		return nil, nil, err;
	}

	redisUserVideoKey := model.RedisUserVideoKey(userId)
	videoIds, err := rc.client.LRange(context.Background(), redisUserVideoKey, 0, -1).Result()

	return &user, videoIds, nil
}

func (rc *RedisUserStore) SaveVideoId(videoId string, userId string) error {
	redisUserVideoKey := model.RedisUserVideoKey(userId)
	_, err := rc.client.LPush(context.Background(), redisUserVideoKey, videoId).Result()

	if err != nil {
		return err;
	}

	return nil
}

func (rc *RedisUserStore) GetVideoIds(userId string) ([]string, error) {
	redisUserVideoKey := model.RedisUserVideoKey(userId);
	videosId, err := rc.client.LRange(context.Background(), redisUserVideoKey, 0, -1).Result()

	if err != nil {
		return nil, err;
	}

	return videosId, nil
}

func (rc *RedisUserStore) IsCacheMiss(err error) bool {return err == redis.Nil}
func (rc *RedisUserStore) IsCacheHit(err error) bool {return err == nil}

// func (rc *RedisUserStore) VideoOnly(videoId string) (*model.Videos, error) {
// 	var video model.Videos;
// 	redisVideoKey := model.RedisVideoKey(videoId)
// 	if err := rc.client.HGetAll(context.Background(), redisVideoKey).Scan(&video); err != nil {
// 		return nil, err;
// 	}
// 	return &video, nil;
// }

// func (rc *RedisUserStore) AllVideo(userId string) ([]*model.Videos, error){
// 	var videos []*model.Videos;
// 	var errs error;

// 	redisUserVideoKey := model.RedisUserVideoKey(userId);
// 	videoLists, err := rc.client.LRange(context.Background(), redisUserVideoKey, 0, -1).Result();

// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, video := range videoLists {
		
// 			var vid model.Videos;
// 			if err := rc.client.HGetAll(context.Background(), video).Scan(&vid); err != nil {
// 				errs = errors.New(err.Error())
// 				continue;
// 			}
// 			videos = append(videos, &vid);
// 	}

// 	return videos, errs;
// }