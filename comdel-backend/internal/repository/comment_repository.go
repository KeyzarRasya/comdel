package repository

import (
	"context"

	"github.com/KeyzarRasya/comdel-server/internal/model"
	"github.com/jackc/pgx/v5"
)

type CommentRepository interface {
	GetByVideoId(videoId string)				([]*model.Comment, error);
}

type CommentRepositoryImpl struct {
	conn *pgx.Conn;
}

func NewCommentRepository(pgxConn *pgx.Conn) CommentRepository {
	return &CommentRepositoryImpl{conn: pgxConn}
}

func (cr *CommentRepositoryImpl) GetByVideoId(videoId string) ([]*model.Comment, error) {
	var comments []*model.Comment;
	var yCommentId, publishedAt, channelId, channelUrl, displayName, profileUrl, textDisplay string;
	var isDetected bool;

	rows, err := cr.conn.Query(
		context.Background(),
		"SELECT DISTINCT y_comment_id, published_at, channel_id, author_channel_url, display_name, profile_url, text_display, isdetected FROM comments WHERE video_id=$1",
		videoId,
	)

	if err != nil {
		return nil, err
	}

	_, err = pgx.ForEachRow(rows, []any{&yCommentId, &publishedAt, &channelId, &channelUrl, &displayName, &profileUrl, &textDisplay, &isDetected}, func () error  {
		comments = append(comments, &model.Comment{
			Yid: yCommentId,
			PublishedAt: publishedAt,
			ChannelId: channelId,
			ChannelUrl: channelUrl,
			DisplayName: displayName,
			ProfileUrl: profileUrl,
			TextDisplay: textDisplay,
			Isdetected: isDetected,
		})
		return nil
	})

	if err != nil {
		return nil, err;
	}

	return comments, nil

}