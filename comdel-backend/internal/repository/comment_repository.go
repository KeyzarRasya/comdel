package repository

import (
	"context"

	"comdel-backend/internal/config"
	"comdel-backend/internal/dto"
	"comdel-backend/internal/model"
)

type CommentRepository interface {
	GetByVideoId(videoId string)				([]*model.Comment, error);
	Save(tx config.DBTx, comment dto.Comment, isDetected bool) (string, error);
}

type CommentRepositoryImpl struct {
	conn config.DBConn;
}

func NewCommentRepository(pgxConn config.DBConn) CommentRepository {
	return &CommentRepositoryImpl{conn: pgxConn}
}

func (cr *CommentRepositoryImpl) GetByVideoId(videoId string) ([]*model.Comment, error) {
	var comments []*model.Comment;

	rows, err := cr.conn.Query(
		context.Background(),
		"SELECT DISTINCT y_comment_id, published_at, channel_id, author_channel_url, display_name, profile_url, text_display, isdetected FROM comments WHERE video_id=$1",
		videoId,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment model.Comment;
		err := rows.Scan(&comment.Yid, &comment.PublishedAt, &comment.ChannelId, &comment.ChannelUrl, &comment.DisplayName, &comment.ProfileUrl, &comment.TextDisplay, &comment.Isdetected)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	
	if err := rows.Err() ;err != nil {
		return nil, err;
	}

	return comments, nil

}

func (cr *CommentRepositoryImpl) Save(tx config.DBTx, comment dto.Comment, isDetected bool) (string, error){
	var commentId string;

	err := tx.QueryRow(
		context.Background(),
		"INSERT INTO comments(y_comment_id, published_at, channel_id, author_channel_url, display_name, profile_url, text_display, video_id, isdetected) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING comment_id",
		comment.Yid, comment.PublishedAt, comment.ChannelId, comment.ChannelUrl, comment.DisplayName, comment.ProfileUrl, comment.TextDisplay, comment.VideoId, isDetected,
	).Scan(&commentId);

	if err != nil {
		return "", err
	}

	return commentId, nil
}