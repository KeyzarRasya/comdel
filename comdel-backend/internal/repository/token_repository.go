package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
)

type TokenRepository interface {
	Save(tx pgx.Tx, token *oauth2.Token, userId string)		error;
	GetByOwnerId(ownerId string)							(oauth2.Token, error)
}

type TokenRepositoryImpl struct {
	conn pgx.Conn;
}

func (tr *TokenRepositoryImpl) Save(tx pgx.Tx, token *oauth2.Token, userId string) error {
	_, err := tx.Exec(
		context.Background(),
		"INSERT INTO oauth_token (access_token, refresh_token, expiry, owner_id) VALUES ($1, $2, $3, $4)",
		token.AccessToken, token.RefreshToken, token.Expiry, userId,
	)

	return err;
}

