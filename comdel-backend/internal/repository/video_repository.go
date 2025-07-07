package repository

import (
	"context"
	"time"

	"github.com/KeyzarRasya/comdel-server/internal/model"
	"github.com/jackc/pgx/v5"
)

type VideoRepository interface {
	GetById(videoId string)											(*model.Videos, error)
	GetYoutubeIdById(id string)										(string, error)
	Save(video model.Videos)										error
	UpdateVideo(videoId string, cooldown time.Time, userId string)	error
	CheckOwnership(userChId string, vidId string)					(bool, error)
	IsCanUpload(userd string)										(bool, error)
}

type VideoRepositoryImpl struct {
	conn *pgx.Conn;
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