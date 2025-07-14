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


func TestGetUser_EmptyCookie(t *testing.T) {
	service := NewUserService(nil, nil, nil, nil, nil, nil, nil)
	resp := service.GetUser("")

	if resp.Status != fiber.StatusBadRequest {
		t.Errorf("Expected StatusBadRequest, got %d", resp.Status)
	}
	if resp.Message != "Failed to get credentials informations" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}
}

func TestGetUser_InvalidCookie(t *testing.T) {
	auth := mock.MockAuthenticator{
		GetUserIdByCookieFunc: func(cookie string) (string, error) {
			return "", errors.New("invalid cookie")
		},
	}

	service := NewUserService(nil, nil, nil, &auth, nil, nil, nil)
	resp := service.GetUser("invalid-cookie")

	if resp.Status != fiber.StatusBadRequest {
		t.Errorf("Expected StatusBadRequest, got %d", resp.Status)
	}
	if resp.Message != "Cannot verify JWT Token" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}
}

func TestGetUser_UserNotFound(t *testing.T) {
	auth := mock.MockAuthenticator{
		GetUserIdByCookieFunc: func(cookie string) (string, error) {
			return "123", nil
		},
	}

	userRepo := mock.MockUserRepository{
		GetByIdWithVideoFunc: func(id string) (*model.User, []string, error) {
			return nil, nil, errors.New("user not found")
		},
	}

	service := NewUserService(&userRepo, nil, nil, &auth, nil, nil, nil)
	resp := service.GetUser("valid-cookie")

	if resp.Status != fiber.StatusBadRequest {
		t.Errorf("Expected StatusBadRequest, got %d", resp.Status)
	}
	if resp.Message != "Failed to get user information" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}
}

func TestGetUser_VideoFetchFailed(t *testing.T) {
	auth := mock.MockAuthenticator{
		GetUserIdByCookieFunc: func(cookie string) (string, error) {
			return "123", nil
		},
	}

	userRepo := mock.MockUserRepository{
		GetByIdWithVideoFunc: func(id string) (*model.User, []string, error) {
			return &model.User{UserId: "123"}, []string{"v1", "v2"}, nil
		},
	}

	videoRepo := mock.MockVideoRepository{
		GetByIdFunc: func(videoId string) (*model.Videos, error) {
			return nil, errors.New("video fetch failed")
		},
	}

	service := NewUserService(&userRepo, nil, &videoRepo, &auth, nil, nil, nil)
	resp := service.GetUser("valid-cookie")

	if resp.Status == fiber.StatusBadRequest {
		t.Errorf("Expected StatusBadRequest, got %d", resp.Status)
	}
	if resp.Message != "Failed to fetch video" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}
}

func TestGetUser_Success(t *testing.T) {
	auth := mock.MockAuthenticator{
		GetUserIdByCookieFunc: func(cookie string) (string, error) {
			return "123", nil
		},
	}

	userRepo := mock.MockUserRepository{
		GetByIdWithVideoFunc: func(id string) (*model.User, []string, error) {
			return &model.User{
				UserId: id,
				Name:   "Mock User",
			}, []string{"v1", "v2"}, nil
		},
	}

	videoRepo := mock.MockVideoRepository{
		GetByIdFunc: func(videoId string) (*model.Videos, error) {
			return &model.Videos{
				Id:    videoId,
				Title: "Mock Video",
			}, nil
		},
	}

	service := NewUserService(&userRepo, nil, &videoRepo, &auth, nil, nil, nil)
	resp := service.GetUser("valid-cookie")

	if resp.Status != fiber.StatusOK {
		t.Errorf("Expected StatusOK, got %d", resp.Status)
	}
	if resp.Message != "Success fetching user data" {
		t.Errorf("Unexpected message: %s", resp.Message)
	}

	user, ok := resp.Data.(*model.User)
	if !ok {
		t.Fatalf("Expected *model.User, got %T", resp.Data)
	}
	if user.UserId != "123" || user.Name != "Mock User" {
		t.Errorf("User data mismatch")
	}
	if len(user.Videos) != 2 {
		t.Errorf("Expected 2 videos, got %d", len(user.Videos))
	}
}

func TestSaveUser_LoadDBFailed(t *testing.T) {
	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return nil, errors.New("Failed to load Database")
		},
	}

	userService := NewUserService(nil, nil, nil, nil, &mockDBLoader, nil, nil)

	res := userService.SaveUser(dto.GoogleProfile{}, nil)

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("Expected to return bad reqesut, but got %d", res.Status)
	}

	if res.Message != "Failed to load database" {
		t.Errorf("Unexpected message : %s", res.Message)
	}
}

func TestSaveUser_TransactionFailed(t *testing.T) {
	dbConnBegin := mock.MockDBConn{
		BeginFunc: func(ctx context.Context) (config.DBTx, error) {
			return nil, errors.New("Failed to started transaction")
		},
	}

	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return &dbConnBegin, nil
		},
	}

	userService := NewUserService(nil, nil, nil, nil, &mockDBLoader, nil, nil);

	res := userService.SaveUser(dto.GoogleProfile{}, nil);

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("Expected to be bad request, founded that %d", res.Status);
	}

	if res.Message != "failed to start transaction" {
		t.Errorf("Unexpected message, got %s", res.Message);
	}
}

func TestSaveUser_GetChannelInfoFailed(t *testing.T) {
	dbTx := mock.MockDBTx{
		RollbackFunc: func(ctx context.Context) error {
			return nil
		},
	}
	dbConn := mock.MockDBConn{
		BeginFunc: func(ctx context.Context) (config.DBTx, error) {
			return &dbTx, nil
		},
	}

	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return &dbConn, nil
		},
	}

	mockYoutubeservice := mock.MockYoutubeService{
		ChannelInfoFunc: func(t *oauth2.Token) (*youtube.Channel, error) {
			return nil, errors.New("Failed to get channel info")
		},
	}

	UserService := NewUserService(nil, nil, nil, nil, &mockDBLoader, nil, mockYoutubeservice)
	res := UserService.SaveUser(dto.GoogleProfile{}, nil)

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("Unexpected response status, found : %d", res.Status)
	}

	if res.Message != "Failed to get channel info" {
		t.Errorf("Unexpected message, found %s", res.Message)
	}
}

func TestSaveUser_SaveFailed(t *testing.T) {
	dbTx := mock.MockDBTx{
		RollbackFunc: func(ctx context.Context) error {
			return nil
		},
	}
	dbConn := mock.MockDBConn{
		BeginFunc: func(ctx context.Context) (config.DBTx, error) {
			return &dbTx, nil
		},
	}

	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return &dbConn, nil
		},
	}

	mockYoutubeService := mock.MockYoutubeService{
		ChannelInfoFunc: func(t *oauth2.Token) (*youtube.Channel, error) {
			channel := youtube.Channel{
				Id: "123",
				Snippet: &youtube.ChannelSnippet{
					Title: "Keyzar",
				},
			}
			return &channel, nil
		},
	}

	mockUserRepo := mock.MockUserRepository{
		IsGIDAvailFunc: func(tx config.DBTx, gid string, googleId *string) (bool, error) {
			return false, nil
		},
		SaveReturningIdFunc: func(tx config.DBTx, user model.User, userId *string) error {
			return errors.New("Failed to save user")
		},
	}

	userService := NewUserService(&mockUserRepo, nil, nil, nil, &mockDBLoader, nil, mockYoutubeService)

	res := userService.SaveUser(dto.GoogleProfile{}, nil);

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("Unexpected error status, found %d", res.Status)
	}

	if res.Message != "failed to insert a profile" {
		t.Errorf("Unexpected message, found %s", res.Message)
	}
}

func TestSaveUser_SaveTokenFailed(t *testing.T) {
	dbTx := mock.MockDBTx{
		RollbackFunc: func(ctx context.Context) error {
			return nil
		},
	}
	dbConn := mock.MockDBConn{
		BeginFunc: func(ctx context.Context) (config.DBTx, error) {
			return &dbTx, nil
		},
	}

	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return &dbConn, nil
		},
	}

	mockYoutubeService := mock.MockYoutubeService{
		ChannelInfoFunc: func(t *oauth2.Token) (*youtube.Channel, error) {
			channel := youtube.Channel{
				Id: "123",
				Snippet: &youtube.ChannelSnippet{
					Title: "Keyzar",
				},
			}
			return &channel, nil
		},
	}

	mockUserRepo := mock.MockUserRepository{
		IsGIDAvailFunc: func(tx config.DBTx, gid string, googleId *string) (bool, error) {
			return false, nil
		},
		SaveReturningIdFunc: func(tx config.DBTx, user model.User, userId *string) error {
			return nil
		},
	}

	mockTokenRepo := mock.MockTokenRepository{
		SaveFunc: func(tx config.DBTx, token *oauth2.Token, userId string) error {
			return errors.New("Failed to save token")
		},
	}


	userService := NewUserService(&mockUserRepo, &mockTokenRepo, nil, nil, &mockDBLoader, nil, mockYoutubeService)

	res := userService.SaveUser(dto.GoogleProfile{}, nil);

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("Unexpected error status, found %d", res.Status)
	}

	if res.Message != "Failed to save token" {
		t.Errorf("Unexpected message, found %s", res.Message)
	}
}

func TestSaveUser_GetIDFailed(t *testing.T) {
	dbTx := mock.MockDBTx{
		RollbackFunc: func(ctx context.Context) error {
			return nil
		},
	}
	dbConn := mock.MockDBConn{
		BeginFunc: func(ctx context.Context) (config.DBTx, error) {
			return &dbTx, nil
		},
	}

	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return &dbConn, nil
		},
	}

	mockYoutubeService := mock.MockYoutubeService{
		ChannelInfoFunc: func(t *oauth2.Token) (*youtube.Channel, error) {
			channel := youtube.Channel{
				Id: "123",
				Snippet: &youtube.ChannelSnippet{
					Title: "Keyzar",
				},
			}
			return &channel, nil
		},
	}

	mockUserRepo := mock.MockUserRepository{
		IsGIDAvailFunc: func(tx config.DBTx, gid string, googleId *string) (bool, error) {
			return false, nil
		},
		SaveReturningIdFunc: func(tx config.DBTx, user model.User, userId *string) error {
			return nil
		},
		GetIDByGIDFunc: func(tx config.DBTx, googleId string) (string, error) {
			return "", errors.New("Failed to get id")
		},
	}

	mockTokenRepo := mock.MockTokenRepository{
		SaveFunc: func(tx config.DBTx, token *oauth2.Token, userId string) error {
			return nil
		},
	}

	userService := NewUserService(&mockUserRepo, &mockTokenRepo, nil, nil, &mockDBLoader, nil, &mockYoutubeService);
	res := userService.SaveUser(dto.GoogleProfile{}, nil)

	if res.Status != fiber.StatusBadRequest {
		t.Errorf("Unexpected status, found %d", res.Status)
	}

	if res.Message != "Failed to get id by" {
		t.Errorf("Unexpected message, found %s", res.Message)
	}
}

func TestSaveUser_Success(t *testing.T) {
	dbTx := mock.MockDBTx{
		RollbackFunc: func(ctx context.Context) error {
			return nil
		},
		CommitFunc: func(ctx context.Context) error {
			return nil
		},
	}
	dbConn := mock.MockDBConn{
		BeginFunc: func(ctx context.Context) (config.DBTx, error) {
			return &dbTx, nil
		},
	}

	mockDBLoader := mock.MockDBLoader{
		LoadFunc: func() (config.DBConn, error) {
			return &dbConn, nil
		},
	}

	mockYoutubeService := mock.MockYoutubeService{
		ChannelInfoFunc: func(t *oauth2.Token) (*youtube.Channel, error) {
			channel := youtube.Channel{
				Id: "123",
				Snippet: &youtube.ChannelSnippet{
					Title: "Keyzar",
				},
			}
			return &channel, nil
		},
	}

	mockUserRepo := mock.MockUserRepository{
		IsGIDAvailFunc: func(tx config.DBTx, gid string, googleId *string) (bool, error) {
			return false, nil
		},
		SaveReturningIdFunc: func(tx config.DBTx, user model.User, userId *string) error {
			return nil
		},
		GetIDByGIDFunc: func(tx config.DBTx, googleId string) (string, error) {
			return "", nil
		},
		GrantSubscriptionAccessFunc: func(tx config.DBTx, userId string) error {
			return nil
		},
	}

	mockAuthentication := mock.MockAuthenticator{
		GenerateTokenFunc: func(userId string) (string, error) {
			return "token-123", nil
		},
	}

	mockTokenRepo := mock.MockTokenRepository{
		SaveFunc: func(tx config.DBTx, token *oauth2.Token, userId string) error {
			return nil
		},
	}

	userService := NewUserService(&mockUserRepo, &mockTokenRepo, nil, &mockAuthentication, &mockDBLoader, nil, mockYoutubeService)
	res := userService.SaveUser(dto.GoogleProfile{}, nil)

	if res.Status != fiber.StatusOK {
		t.Errorf("Expected to be Successfull or status OK, but found %d", res.Status)
	}

	if res.Message != "Success" {
		t.Errorf("Unexpected message, found %s", res.Message)
	}
}