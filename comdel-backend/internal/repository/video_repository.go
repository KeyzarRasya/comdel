package repository

import (
	"context"

	"comdel-backend/internal/config"
	"comdel-backend/internal/model"
)

type VideoRepository interface {
	GetById(videoId string)											(*model.Videos, error)
	Save(tx config.DBTx, video model.Videos)								error
	UpdateComment(tx config.DBTx, commentsId []string, videoId string, isDeleted bool) error
	// UpdateVideo(videoId string, cooldown time.Time, userId string)	error
}

/*_, err = tx.Exec(
		context.Background(),
		"UPDATE videos SET deleted_comments = $1 WHERE videos_id=$2",
		deletedCommentsId, videoId,
	)*/

type VideoRepositoryImpl struct {
	conn config.DBConn;
}

func NewVideoRepository(pgxConn config.DBConn) VideoRepository {
	return &VideoRepositoryImpl{conn: pgxConn}
}

func (vr *VideoRepositoryImpl) Save(tx config.DBTx, video model.Videos) error {
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

func (vr *VideoRepositoryImpl) UpdateComment(tx config.DBTx, commentsId []string, videoId string, isDetected bool) error {
	var err error;
	if isDetected {
		_, err = tx.Exec(
			context.Background(),
			"UPDATE videos SET deleted_comments = $1 WHERE videos_id=$2",
			commentsId, videoId,
		)
	} else {
		_, err = tx.Exec(
			context.Background(),
			"UPDATE videos SET comments = $1 WHERE videos_id=$2",
			commentsId, videoId,
		)
	}
	return err
}