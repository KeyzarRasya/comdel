package mock

import "comdel-backend/internal/model"

type MockRedisUserStore struct {
	SaveUserFunc func(user model.User) error;
	GetUserAndVideoFunc func(userId string) (*model.User, []string, error);
	SaveVideoIdFunc func(videoId string, userId string) error
	GetVideoIdsFunc func(userId string) ([]string, error)
	IsCacheMissFunc func(err error) bool
	IsCacheHitFunc	func(err error) bool
}

func (mr *MockRedisUserStore) SaveUser(user model.User) error {
	return mr.SaveUserFunc(user);
}

func (mr *MockRedisUserStore) GetUserAndVideo(userId string) (*model.User, []string, error) {
	return mr.GetUserAndVideoFunc(userId)
}

func (mr *MockRedisUserStore) SaveVideoId(videoId string, userId string) error {
	return mr.SaveVideoIdFunc(videoId, userId);
}

func (mr *MockRedisUserStore) GetVideoIds(userId string) ([]string, error) {
	return mr.GetVideoIdsFunc(userId)
}

func (mr *MockRedisUserStore) IsCacheMiss(err error) bool {
	return mr.IsCacheMissFunc(err);
}

func (mr *MockRedisUserStore) IsCacheHit(err error) bool {
	return mr.IsCacheHitFunc(err);
}