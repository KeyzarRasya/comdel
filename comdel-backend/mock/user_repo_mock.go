package mock

import (
	"comdel-backend/internal/config"
	"comdel-backend/internal/model"
	"time"
)

type MockUserRepository struct {
	GetByIdWithVideoFunc 	func(id string) 					(*model.User, []string, error);
	GetIDByGIDFunc 			func(tx config.DBTx, googleId string)	(string, error);
	GetNameByIdFunc 		func(userId string)					(string, error);
	GetYoutubeIdByIdFunc 	func(id string)						(string, error);
	GetSubsIdByIdFunc 		func(userId string)					(string, error);

	IsCooldownFunc 			func(userId string)								(bool, error);
	IsGIDAvailFunc			func(tx config.DBTx, gid string, googleId *string)	(bool, error);

	SaveReturningIdFunc 		func(tx config.DBTx, user model.User)	(string, error);
	DeactivateSubscriptionFunc 	func(tx config.DBTx, userId string)						error;

	UpdateVideosFunc 			func(tx config.DBTx, videoId string, userId string, cooldown time.Time)	error
	GrantSubscriptionAccessFunc func(tx config.DBTx, userId string) error;
}

func (mur *MockUserRepository) GetByIdWithVideo(id string) (*model.User, []string, error){
	return mur.GetByIdWithVideoFunc(id)
}

func (mur *MockUserRepository)GetIDByGID(tx config.DBTx, googleId string) (string, error) {
	return mur.GetIDByGIDFunc(tx, googleId);
}

func (mur *MockUserRepository) GetNameById(userId string) (string, error) {
	return mur.GetNameByIdFunc(userId)
}

func (mur *MockUserRepository) GetYoutubeIdById(id string) (string, error) {
	return mur.GetYoutubeIdByIdFunc(id)
}

func (mur *MockUserRepository) GetSubsIdById(userId string) (string, error) {
	return mur.GetSubsIdByIdFunc(userId)
}

func (mur *MockUserRepository) IsCooldown(userId string) (bool, error) {
	return mur.IsCooldownFunc(userId)
}

func (mur *MockUserRepository) IsGIDAvail(tx config.DBTx, gid string, googleId *string) (bool, error) {
	return mur.IsGIDAvailFunc(tx, gid, googleId);
}

func (mur *MockUserRepository) SaveReturningId(tx config.DBTx, user model.User) (string, error) {
	return mur.SaveReturningIdFunc(tx, user);
}

func (mur *MockUserRepository)DeactivateSubscription(tx config.DBTx, userId string) error {
	return mur.DeactivateSubscriptionFunc(tx, userId)
}

func (mur *MockUserRepository) UpdateVideos(tx config.DBTx, videoId string, userId string, cooldown time.Time) error {
	return mur.UpdateVideosFunc(tx, videoId, userId, cooldown)
}

func (mur *MockUserRepository) GrantSubscriptionAccess(tx config.DBTx, userId string) error {
	return mur.GrantSubscriptionAccessFunc(tx, userId)
}
