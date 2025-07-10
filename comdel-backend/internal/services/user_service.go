package services

import (
	"context"

	"github.com/KeyzarRasya/comdel-server/internal/config"
	"github.com/KeyzarRasya/comdel-server/internal/dto"
	"github.com/KeyzarRasya/comdel-server/internal/helper"
	"github.com/KeyzarRasya/comdel-server/internal/model"
	"github.com/KeyzarRasya/comdel-server/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type UserService interface{
	SaveUser(user dto.GoogleProfile, oauthToken *oauth2.Token) dto.Response
	GetUser(cookies string) dto.Response
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository;
	TokenRepository repository.TokenRepository;
	VideoRepository repository.VideoRepository;
}

/*
	Constructor for Creating
	NewUserService
	=====
	Also injecting dependency
*/
func NewUserService(
	userRepository 		repository.UserRepository,
	tokenRepository 	repository.TokenRepository,
	videoRepository 	repository.VideoRepository,
) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository, 
		TokenRepository: tokenRepository, 
		VideoRepository: videoRepository}
}


func (us *UserServiceImpl) SaveUser(user dto.GoogleProfile, oauthToken *oauth2.Token) dto.Response {
	conn := config.LoadDatabase();		/*Load Database*/
	oauthConfig := config.OAuthConfig()	/*Load OAuth Config*/

	var googleId string;		/*value to store g_id (it is available or not)*/
	var userId string;

	tx, err := conn.Begin(context.Background());	/* Starting database transaction */

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "failed to start transaction",
			Data: nil,
		}
	}

	defer tx.Rollback(context.Background())

	tokenSource := oauthConfig.TokenSource(context.Background(), oauthToken);
	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(tokenSource));

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to create youtube service",
			Data: err,
		}
	}

	channel, err := youtubeService.Channels.List([]string{"id", "snippet"}).Mine(true).Do();

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get youtube channel information",
			Data: err,
		}
	}

	_, err = us.UserRepository.IsGIDAvail(tx, user.GId, &googleId)
	
	if err != nil {
		return dto.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Database error checking g_id",
			Data:    err.Error(),
		}
	}

	if googleId == "" {
		log.Info("Creating a new account");

		modelUser := user.Parse()
		modelUser.YoutubeId = channel.Items[0].Id;
		modelUser.TitleSnippet = channel.Items[0].Snippet.Title;

		if  err := us.UserRepository.SaveReturningId(tx, modelUser, &userId); err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "failed to insert a profile",
				Data: err.Error(),
			}
		}

		if err := us.TokenRepository.Save(tx, oauthToken, userId); err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to save token",
				Data: err,
			}	
		}
	}

	userId, err = us.UserRepository.GetIDByGID(tx, googleId);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message:"Failed to get id by",
			Data: err.Error(),
		}
	}

	jwt, err := helper.GenerateToken(userId);
	if err != nil{
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to generate JWT Token for authentication",
			Data: nil,
		}
	}

	user.Token = jwt;

	if err := us.UserRepository.GrantSubscriptionAccess(tx, userId); err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to grant video access",
			Data: err.Error(),
		}
	}

	tx.Commit(context.Background());

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Success",
		Data: user,
	}
}


func (us *UserServiceImpl) GetUser(cookies string) dto.Response {
	var videosId []string;
	var videos []*model.Videos;
	
	if cookies == "" {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get credentials informations",
			Data: nil,
		}
	}

	userId, err := helper.VerifyAndGet(cookies)
	if err != nil {
		log.Info("Error : ", err);
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Cannot verify JWT Token",
			Data: err,
		}
	}

	user, videosId, err := us.UserRepository.GetByIdWithVideo(userId);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get user information",
			Data: err,
		}
	}


	for _, id := range videosId {
		var video *model.Videos;

		video, err = us.VideoRepository.GetById(id);
		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to fetch video",
				Data: nil,
			}
		}

		video.Id = id;
		videos = append(videos, video);
	}

	user.Videos = videos;

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Success fetching user data",
		Data: user,
	}
	
}
