package repository

import (
	"github.com/KeyzarRasya/comdel-server/internal/dto"
	"github.com/KeyzarRasya/comdel-server/internal/model"
)

type CommentRepository interface {
	Save(comment dto.Comment, detected bool)	error		
	GetByVideoId(videoId string)				([]model.Comment, error);
}