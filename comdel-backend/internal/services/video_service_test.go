package services

import (
	"comdel-backend/internal/config"
	"comdel-backend/internal/dto"
	"comdel-backend/internal/model"
	"comdel-backend/mock"
	"context"
	"errors"
	"testing"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

/*
	GetById(videoId string)											(*model.Videos, error)
	Save(tx pgx.Tx, video model.Videos)								error
*/

func TestUploadVideo_VideoMetadataNotComplete(t *testing.T) {
	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return nil, nil
		},
	}

	mockAutheticator := mock.MockAuthenticator{
		GetUserIdByCookieFunc: func(cookie string) (string, error) {
			return "userid-123", nil
		},
	}

	metadata := dto.UploadVideos{
		Link: "https://youtube.com/keyzarrasya/?v=123",
		Scheduler: "24",
		Strategy: "",
	}

	videoService := NewVideoService(nil, nil, nil, nil, nil, nil, &mockDBLoader, &mockAutheticator)

	res := videoService.UploadVideo("cookie-123", metadata)

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("unexpected status, got %d", res.Status)
	}

	if res.Message != "invalid request query" {
		t.Errorf("Unexpected error message, found %s", res.Message)
	}

}

func TestUploadVideo_GetYoutubeIdFailed(t *testing.T) {
	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return nil, nil
		},
	}

	mockAutheticator := mock.MockAuthenticator{
		GetUserIdByCookieFunc: func(cookie string) (string, error) {
			return "userid-123", nil
		},
	}

	metadata := dto.UploadVideos{
		Link: "https://youtube.com/keyzarrasya/?v=123",
		Scheduler: "24",
		Strategy: "AUTO",
	}

	mockUserRepo := mock.MockUserRepository{
		GetYoutubeIdByIdFunc: func(id string) (string, error) {
			return "", errors.New("Failed to get youtube id")
		},
	}

	videoService := NewVideoService(&mockUserRepo, nil, nil, nil, nil,nil, &mockDBLoader, &mockAutheticator);
	res := videoService.UploadVideo("cookie-123", metadata)

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("Expected to error")
	}

	if res.Message != "failed to get youtube id" {
		t.Errorf("unexpected message, found %s", res.Message)
	}
}

func TestUploadVideo_VideoListFailed(t *testing.T) {
	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return nil, nil
		},
	}

	mockAutheticator := mock.MockAuthenticator{
		GetUserIdByCookieFunc: func(cookie string) (string, error) {
			return "userid-123", nil
		},
	}

	metadata := dto.UploadVideos{
		Link: "https://youtube.com/keyzarrasya/?v=123",
		Scheduler: "24",
		Strategy: "AUTO",
	}

	mockTokenRepo := mock.MockTokenRepository{
		GetByOwnerIdFunc: func(ownerId string) (*oauth2.Token, error) {
			return &oauth2.Token{}, nil
		},
	}

	mockUserRepo := mock.MockUserRepository{
		GetYoutubeIdByIdFunc: func(id string) (string, error) {
			return "", nil
		},
	}

	mockYtService := mock.MockYoutubeService{
		VideoFunc: func(t *oauth2.Token, s string) (*youtube.Video, error) {
			return nil, errors.New("Failed to get video list")
		},
	}

	videoService := NewVideoService(&mockUserRepo, nil, &mockTokenRepo, nil, nil, mockYtService, &mockDBLoader, &mockAutheticator)
	res := videoService.UploadVideo("cookie-123", metadata)

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("Expected to be error, found %d", res.Status)
	}

	if res.Message != "failed to get videos" {
		t.Errorf("Unexpected message, found %s", res.Message)
	}
}

func TestUploadVideo_NotOwnerOfVideo(t *testing.T) {
	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return nil, nil
		},
	}

	mockAutheticator := mock.MockAuthenticator{
		GetUserIdByCookieFunc: func(cookie string) (string, error) {
			return "userid-123", nil
		},
	}

	metadata := dto.UploadVideos{
		Link: "https://youtube.com/keyzarrasya/?v=123",
		Scheduler: "24",
		Strategy: "AUTO",
	}

	mockTokenRepo := mock.MockTokenRepository{
		GetByOwnerIdFunc: func(ownerId string) (*oauth2.Token, error) {
			return &oauth2.Token{}, nil
		},
	}

	mockUserRepo := mock.MockUserRepository{
		GetYoutubeIdByIdFunc: func(id string) (string, error) {
			return "channel-1235", nil
		},
	}

	mockYtService := mock.MockYoutubeService{
		VideoFunc: func(t *oauth2.Token, s string) (*youtube.Video, error) {
			vidItem := youtube.Video{
				Snippet: &youtube.VideoSnippet{
					ChannelId: "channel-123",
				},
			}
			return &vidItem, nil
		},
	}

	videoService := NewVideoService(&mockUserRepo, nil, &mockTokenRepo, nil, nil, mockYtService, &mockDBLoader, &mockAutheticator)
	res := videoService.UploadVideo("cookie-123", metadata)

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("expected to error, found %d", res.Status)
	}

	if res.Message != "action not permitted, you are not the owner of the video" {
		t.Errorf("unexpected error message, found %s", res.Message)
	}

}

func TestUploadVideo_SaveVideo(t *testing.T) {
	mockDBTx := mock.MockDBTx{
		RollbackFunc: func(ctx context.Context) error {
			return nil
		},
	}
	mockDBConn := mock.MockDBConn{
		BeginFunc: func(ctx context.Context) (config.DBTx, error) {
			return &mockDBTx, nil
		},
	}
	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return &mockDBConn, nil
		},
	}


	mockAutheticator := mock.MockAuthenticator{
		GetUserIdByCookieFunc: func(cookie string) (string, error) {
			return "userid-123", nil
		},
	}

	metadata := dto.UploadVideos{
		Link: "https://youtube.com/keyzarrasya/?v=123",
		Scheduler: "24",
		Strategy: "AUTO",
	}

	mockTokenRepo := mock.MockTokenRepository{
		GetByOwnerIdFunc: func(ownerId string) (*oauth2.Token, error) {
			return &oauth2.Token{}, nil
		},
	}

	mockUserRepo := mock.MockUserRepository{
		GetYoutubeIdByIdFunc: func(id string) (string, error) {
			return "channel-123", nil
		},
	}

	mockYtService := mock.MockYoutubeService{
		VideoFunc: func(t *oauth2.Token, s string) (*youtube.Video, error) {
			vidItem := youtube.Video{
				Id: "123",
				Snippet: &youtube.VideoSnippet{
					Title: "title-123",
					ChannelId: "channel-123",
					Thumbnails: &youtube.ThumbnailDetails{
						Maxres: &youtube.Thumbnail{
							Url: "https://host/img",
						},
					},
				},
			}
			return &vidItem, nil
		},
	}

	mockVideoService := mock.MockVideoRepository{
		SaveFunc: func(tx config.DBTx, video model.Videos) error {
			return errors.New("Failed to save item")
		},
	}

	videoService := NewVideoService(&mockUserRepo, &mockVideoService, &mockTokenRepo, nil, nil, mockYtService, &mockDBLoader, &mockAutheticator)
	res := videoService.UploadVideo("cookie-123", metadata)

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("should return error, found %d", res.Status)
	}

	if res.Message != "Failed to add a new videos" {
		t.Errorf("unexpected message, found %s", res.Message)
	}
}