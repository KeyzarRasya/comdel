package repository

import (
	"context"
	"time"

	"comdel-backend/internal/config"
	"comdel-backend/internal/model"
	"comdel-backend/internal/status"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	GetByIdWithVideo(id string)										(*model.User, []string, error);
	GetIDByGID(tx config.DBTx, googleId string)							(string, error);
	GetNameById(userId string)										(string, error);
	GetYoutubeIdById(id string)										(string, error);

	GetSubsIdById(userId string)									(string, error);
	IsCooldown(userId string)										(bool, error)

	IsGIDAvail(tx config.DBTx, gid string, googleId *string)				(bool, error);
	// Save(user model.User)											error;
	SaveReturningId(tx config.DBTx, user model.User)		(string, error);
	DeactivateSubscription(tx config.DBTx, userId string)				error;

	UpdateVideos(tx config.DBTx, videoId string, userId string, cooldown time.Time)	error

	/*
		Restricted Function
	*/
	GrantSubscriptionAccess(tx config.DBTx, userId string) error;
}


type UserRepositoryImpl struct {
	conn config.DBConn;
}

func NewUserRepository(pgxConn config.DBConn) UserRepository {
	return &UserRepositoryImpl{conn: pgxConn}
}

func (ur *UserRepositoryImpl) SaveReturningId(tx config.DBTx, user model.User) (string, error) {
	var userId string
	err := tx.QueryRow(
		context.Background(), 
		"INSERT INTO user_info (name, given_name, email, isverified, picture, g_id, youtube_id, title_snippet) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING user_id",
		user.Name, user.GivenName, user.Email, user.VerifiedEmail, user.Picture, user.GId, user.YoutubeId, user.TitleSnippet,
	).Scan(&userId)

	return userId, err;
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

func (ur *UserRepositoryImpl) GetNameById(userId string) (string, error) {
	var username string;
	err := ur.conn.QueryRow(
		context.Background(),
		"SELECT name FROM user_info WHERE user_id=$1",
		userId,
	).Scan(&username)

	if err != nil {
		return "", err
	}

	return username, nil
}

func (ur *UserRepositoryImpl) GetIDByGID(tx config.DBTx, googleId string) (string, error) {
	var userId string;
	err := tx.QueryRow(
		context.Background(),
		"SELECT user_id FROM user_info WHERE g_id=$1",
		googleId,
	).Scan(&userId);

	if err != nil {
		return "", err;
	}

	return userId, nil;
}

func (ur *UserRepositoryImpl) GetYoutubeIdById(id string) (string, error) {
	var youtubeId string;

	err := ur.conn.QueryRow(
		context.Background(),
		"SELECT youtube_id FROM user_info WHERE user_id=$1",
		id,
	).Scan(&youtubeId)

	if err != nil {
		return "", nil
	}

	return youtubeId, err;
}

func (ur *UserRepositoryImpl) IsGIDAvail(tx config.DBTx, gid string, googleId *string) (bool, error) {
	err := tx.QueryRow(
		context.Background(), 
		"SELECT g_id FROM user_info WHERE g_id=$1", 
		gid,
	).Scan(googleId)

	if err == pgx.ErrNoRows {
		log.Info("True")
		return false, nil
	}

	if err != nil {
		return false, err;
	}

	return true, nil
}

func (ur *UserRepositoryImpl) IsCooldown(userId string) (bool, error) {
	var subscription string
	var cooldown *time.Time

	err := ur.conn.QueryRow(
		context.Background(),
		`SELECT subscription, upload_cooldown FROM user_info WHERE user_id = $1`,
		userId,
	).Scan(&subscription, &cooldown)
	if err != nil {
		return false, err
	}

	if subscription == "NONE" || cooldown == nil || time.Now().After(*cooldown) {
		return false, nil
	}
	return true, nil
}

func (ur *UserRepositoryImpl) GetSubsIdById(userId string) (string, error) {
	var subsId pgtype.Text

	err := ur.conn.QueryRow(
		context.Background(),
		"SELECT subs_id FROM user_info WHERE user_id=$1",
		userId,
	).Scan(&subsId)

	if err == pgx.ErrNoRows {
		return "", nil 
	}

	if err != nil {
		return "", err
	}

	if !subsId.Valid {
		return "", status.ErrNotSubscribed
	}

	return subsId.String, nil
}

func (ur *UserRepositoryImpl) UpdateVideos(tx config.DBTx, videoId string, userId string, cooldown time.Time) error {
	_, err := tx.Exec(
		context.Background(),
		"UPDATE user_info SET videos = array_append(videos, $1), upload_cooldown=$2 WHERE user_id=$3",
		videoId, cooldown, userId,
	)

	return err;
}

func (ur *UserRepositoryImpl) DeactivateSubscription(tx config.DBTx, userId string) error {
	_, err := tx.Exec(
		context.Background(),
		"UPDATE user_info SET subscription = 'NONE', premium_plan='NONE', subs_id=null WHERE user_id=$1",
		userId,
	)

	return err
}

func (ur *UserRepositoryImpl) GrantSubscriptionAccess(tx config.DBTx, userId string) error {
	var subscriptionExpiry time.Time = time.Now().Add((time.Hour * 24) * 30);

	_, err := tx.Exec(
		context.Background(),
		"UPDATE user_info SET premium_plan = 'NEWBIE', subscription = 'ACTIVE', subscription_expiry = $1 WHERE user_id=$2",
		subscriptionExpiry, userId,
	)

	return err
}