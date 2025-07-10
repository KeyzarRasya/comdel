package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
)

type TokenRepository interface {
	Save(tx pgx.Tx, token *oauth2.Token, userId string)		error;
	GetByOwnerId(ownerId string)							(*oauth2.Token, error)
}

type TokenRepositoryImpl struct {
	conn *pgx.Conn;
}

func NewTokenRepository(pgxConn *pgx.Conn) TokenRepository {
	return &TokenRepositoryImpl{conn: pgxConn}
}

func (tr *TokenRepositoryImpl) Save(tx pgx.Tx, token *oauth2.Token, userId string) error {
	_, err := tx.Exec(
		context.Background(),
		"INSERT INTO oauth_token (access_token, refresh_token, expiry, owner_id) VALUES ($1, $2, $3, $4)",
		token.AccessToken, token.RefreshToken, token.Expiry, userId,
	)

	return err;
}

func (tr *TokenRepositoryImpl) GetByOwnerId(ownerId string) (*oauth2.Token, error) {
	var token *oauth2.Token;
	err := tr.conn.QueryRow(
		context.Background(),
		"SELECT access_token, refresh_token, expiry FROM oauth_token WHERE owner_id=$1",
		ownerId,
	).Scan(&token.AccessToken, &token.RefreshToken, &token.Expiry)

	if err != nil {
		return nil, err;
	}

	return token, nil
}

