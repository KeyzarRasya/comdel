package mock

import (
	"comdel-backend/internal/config"
	"comdel-backend/internal/model"
)

type MockVideoRepository struct {
	GetByIdFunc func(videoId string) 					(*model.Videos, error)
	SaveFunc	func(tx config.DBTx, video model.Videos) 	error
	UpdateCommentFunc func(tx config.DBTx, commentsId []string, videoId string, isDeleted bool) error
}

func (mvr *MockVideoRepository)GetById(videoId string) (*model.Videos, error) {
	return mvr.GetByIdFunc(videoId)
}

func (mvr *MockVideoRepository)Save(tx config.DBTx, video model.Videos) error {
	return mvr.SaveFunc(tx, video)
}

func (mvr *MockVideoRepository) UpdateComment(tx config.DBTx, commentsId []string, videoId string, isDeletd bool) error {
	return mvr.UpdateCommentFunc(tx, commentsId, videoId, isDeletd)
}