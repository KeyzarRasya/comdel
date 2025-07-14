package mock

import (
	"comdel-backend/internal/config"

	"golang.org/x/oauth2"
)

type MockTokenRepository struct {
	GetByOwnerIdFunc func(ownerId string) (*oauth2.Token, error);
	SaveFunc func(tx config.DBTx, token *oauth2.Token, userId string) error
}

func (mtr *MockTokenRepository) GetByOwnerId(ownerId string) (*oauth2.Token, error) {
	return mtr.GetByOwnerIdFunc(ownerId);
}

func (mtr *MockTokenRepository) Save(tx config.DBTx, token *oauth2.Token, userId string) error {
	return mtr.SaveFunc(tx, token, userId)
}