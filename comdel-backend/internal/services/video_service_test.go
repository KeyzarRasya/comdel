package services

import (
	"comdel-backend/internal/config"
	"comdel-backend/internal/model"
)

/*
	GetById(videoId string)											(*model.Videos, error)
	Save(tx pgx.Tx, video model.Videos)								error
*/

type MockVideoRepository struct {
	GetByIdFunc func(videoId string) 					(*model.Videos, error)
	SaveFunc	func(tx config.DBTx, video model.Videos) 	error
}

func (mvr *MockVideoRepository)GetById(videoId string) (*model.Videos, error) {
	return mvr.GetByIdFunc(videoId)
}

func (mvr *MockVideoRepository)Save(tx config.DBTx, video model.Videos) error {
	return mvr.SaveFunc(tx, video)
}