package mock

import "comdel-backend/internal/model"

type MockRedisUserStore struct {
	SaveUserFunc func(user model.User) error;
	GetUserFunc func(userId string) (*model.User, error);
}

func (mr *MockRedisUserStore) SaveUser(user model.User) error {
	return mr.SaveUserFunc(user);
}

func (mr *MockRedisUserStore) GetUser(userId string) (*model.User, error) {
	return mr.GetUserFunc(userId)
}