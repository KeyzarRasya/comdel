package repository

import (
	"context"

	"github.com/KeyzarRasya/comdel-server/internal/model"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	GetByIdWithVideo(id string)										(*model.User, []string, error);
	GetIDByGID(gid string)											(string, error);
	GetNameById(userId string)										(string, error);
	GetYoutubeIdById(id string)										(string, error);
	GetVideos(userId string)										([]string, error)
	IsGIDAvail(tx pgx.Tx, gid string, googleId *string)				(bool, error);
	Save(user model.User)											error;
	SaveReturningId(tx pgx.Tx, user model.User, userId *string)		error;
}


type UserRepositoryImpl struct {
	conn *pgx.Conn;
}

func (ur *UserRepositoryImpl) SaveReturningId(tx pgx.Tx, user model.User, userId *string) error {
	err := tx.QueryRow(
		context.Background(), 
		"INSERT INTO user_info (name, given_name, email, isverified, picture, g_id, youtube_id, title_snippet) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING user_id",
		user.Name, user.GivenName, user.Email, user.VerifiedEmail, user.Picture, user.GId, user.YoutubeId, user.TitleSnippet,
	).Scan(&userId)

	return err;
}

func (ur *UserRepositoryImpl) GetByIdWithVideo(id string) (*model.User, []string, error) {
	var user 		model.User;
	var videosId	[]string

	err := ur.conn.QueryRow(
		context.Background(),
		"SELECT user_id, name, email, given_name, picture, videos, g_id, subscription, premium_plan FROM user_info WHERE user_id=$1",
		id,
	).Scan(&user.UserId,&user.Name, &user.Email, &user.GivenName, &user.Picture, &videosId, &user.GId, &user.Subscription, &user.PremiumPlan)


	if err != nil {
		return nil, []string{}, err;
	}

	return &user, videosId, err;
}

func (ur *UserRepositoryImpl) IsGIDAvail(tx pgx.Tx, gid string, googleId *string) (bool, error) {
	err := tx.QueryRow(
		context.Background(), 
		"SELECT g_id FROM user_info WHERE g_id=$1", 
		gid,
	).Scan(googleId)

	if err == pgx.ErrNoRows {
		return true, nil
	}

	if err != nil {
		return false, err;
	}

	return false, nil
}


