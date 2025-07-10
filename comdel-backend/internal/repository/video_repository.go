package repository

import (
	"context"

	"github.com/KeyzarRasya/comdel-server/internal/model"
	"github.com/jackc/pgx/v5"
)

type VideoRepository interface {
	GetById(videoId string)											(*model.Videos, error)
	Save(tx pgx.Tx, video model.Videos)								error
	// UpdateVideo(videoId string, cooldown time.Time, userId string)	error
}

type VideoRepositoryImpl struct {
	conn *pgx.Conn;
}

func NewVideoRepository(pgxConn *pgx.Conn) VideoRepository {
	return &VideoRepositoryImpl{conn: pgxConn}
}

func (vr *VideoRepositoryImpl) Save(tx pgx.Tx, video model.Videos) error {
	_, err := tx.Exec(
		context.Background(),
		"INSERT INTO videos(videos_id, title, owner, thumbnail, strategy, scheduler) VALUES ($1, $2, $3, $4, $5, $6)",
		video.Id, video.Title, video.Owner, video.Thumbnail, video.Strategy, video.Scheduler,
	)

	return err;
}

func (vr *VideoRepositoryImpl) GetById(videoId string) (*model.Videos, error) {
	var video model.Videos;
	
	err := vr.conn.QueryRow(
		context.Background(),
		"SELECT title, thumbnail, strategy, scheduler, owner FROM videos WHERE videos_id=$1",
		videoId,
	).Scan(&video.Title, &video.Thumbnail, &video.Strategy, &video.Scheduler, &video.Owner);

	return &video, err;
}