package services

import (
	"context"

	"comdel-backend/internal/config"
	"comdel-backend/internal/dto"
	"comdel-backend/internal/model"
	"comdel-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/oauth2"
)

type UserService interface{
	SaveUser(user dto.GoogleProfile, oauthToken *oauth2.Token) dto.Response
	GetUser(cookies string) dto.Response
	RedisGetAll(userId string) dto.Response
}

type UserServiceImpl struct {
	UserRepository repository.UserRepository;
	TokenRepository repository.TokenRepository;
	VideoRepository repository.VideoRepository;
	YtService YoutubeService
	OAuth config.OAuthProvider
	DBLoader config.DBLoader
	Authenticator Authenticator
	Redis UserStore
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
	authenticator 		Authenticator,
	dbLoader 			config.DBLoader,
	oAuth 				config.OAuthProvider,
	ytService			YoutubeService,
	redis				UserStore,
) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository, 
		TokenRepository: tokenRepository, 
		VideoRepository: videoRepository,
		Authenticator: authenticator,
		DBLoader: dbLoader,
		OAuth: oAuth,
		YtService: ytService,
		Redis: redis,
	}
}

func (us *UserServiceImpl) RedisGetAll(userId string) dto.Response {
	user,_, err := us.Redis.GetUserAndVideo(userId)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Redis Cache Miss",
			Data: err,
		}
	}

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Redis Cache Hit",
		Data: user,
	}
}


func (us *UserServiceImpl) SaveUser(user dto.GoogleProfile, oauthToken *oauth2.Token) dto.Response {
	conn, err := us.DBLoader.Load()

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to load database",
			Data: err.Error(),
		}
	}

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

	channel, err := us.YtService.ChannelInfo(oauthToken)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get channel info",
			Data: err.Error(),
		}
	}


	isAvail, err := us.UserRepository.IsGIDAvail(tx, user.GId, &googleId)
	if err != nil {
		return dto.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Database error checking g_id",
			Data:    err.Error(),
		}
	}

	log.Info(googleId)
	modelUser := user.Parse()
	if !isAvail {
		modelUser.YoutubeId = channel.Id;
		modelUser.TitleSnippet = channel.Snippet.Title;

		userId, err = us.UserRepository.SaveReturningId(tx, modelUser)
		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "failed to insert a profile",
				Data: err.Error(),
			}
		}

		modelUser.UserId = userId;
		log.Info("UserId", userId)

		if err := us.TokenRepository.Save(tx, oauthToken, userId); err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to save token",
				Data: err,
			}	
		}
	} else {
		userId, err = us.UserRepository.GetIDByGID(tx, googleId);
		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message:"Failed to get id by",
				Data: err.Error(),
			}
		}
		modelUser.UserId = userId;
	}


	

	jwt, err := us.Authenticator.GenerateToken(userId);
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
	
	if err := us.Redis.SaveUser(modelUser); err != nil {
		return dto.Response{
			Status: fiber.StatusOK,
			Message: "Success but redis error",
			Data: user,
		}
	}
	
	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Success",
		Data: user,
	}
}


func (us *UserServiceImpl) GetUser(cookies string) dto.Response {
	var videosId []string;
	var videos []*model.Videos;
	var user *model.User;
	
	if cookies == "" {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get credentials informations",
			Data: nil,
		}
	}

	userId, err := us.Authenticator.GetUserIdByCookie(cookies)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Cannot verify JWT Token",
			Data: err,
		}
	}

	user, videosId, err = us.Redis.GetUserAndVideo(userId)

	if us.Redis.IsCacheMiss(err) {
		user, videosId, err = us.UserRepository.GetByIdWithVideo(userId);
		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to get user information",
				Data: err,
			}
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
