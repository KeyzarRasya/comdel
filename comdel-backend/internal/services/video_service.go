package services

import (
	"context"
	"errors"
	"net/url"
	"time"

	"comdel-backend/internal/config"
	"comdel-backend/internal/dto"
	"comdel-backend/internal/helper"
	"comdel-backend/internal/model"
	"comdel-backend/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type VideoService interface {
	UploadVideo(cookie string, upload dto.UploadVideos) 	dto.Response;
	IsCanUpload(link string, cookie string) 				dto.Response;
	processURL(link string) 								(string, error);
	Info(videoId string, cookies string)					dto.Response;
}

type VideoServiceImpl struct {
	UserRepository 		repository.UserRepository;
	VideoRepository 	repository.VideoRepository;
	TokenRepository 	repository.TokenRepository;
	CommentRepository 	repository.CommentRepository
	CommentService 		CommentService
	YtService 			YoutubeService
	DBLoader 			config.DBLoader
	Authentication		Authenticator
}

/*
	Constructor for Creating
	VideoService
	=====
	Also injecting dependency
*/
func NewVideoService(
	userRepository repository.UserRepository,
	videoRepository repository.VideoRepository,
	tokenRepository repository.TokenRepository,
	commentService CommentService,
	commentRepository repository.CommentRepository,
	ytService YoutubeService,
	dbLoader config.DBLoader,
	authentication Authenticator,
) VideoService {
	return &VideoServiceImpl{
		UserRepository: userRepository,
		VideoRepository: videoRepository,
		TokenRepository: tokenRepository,
		CommentService: commentService,
		CommentRepository: commentRepository,
		YtService: ytService,
		DBLoader: dbLoader,
		Authentication: authentication,
	}
}


func (vs *VideoServiceImpl) UploadVideo(cookie string, upload dto.UploadVideos) dto.Response {
	conn, err := vs.DBLoader.Load()
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to load database",
			Data: err.Error(),
		}
	}
	
	userId, err := vs.Authentication.GetUserIdByCookie(cookie)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "invalid authentication token",
			Data: nil,
		}
	}

	if upload.Link == "" || upload.Scheduler == "" || upload.Strategy == "" {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "invalid request query",
			Data: nil,
		}
	}

	/* Variable Declaration */
	var channelId string;
	var video model.Videos;
	
	/* 
		=== START ===
		Parsing Url End Get the value 
	*/
	videoId, err := vs.processURL(upload.Link)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Faiiled to process given url",
			Data: err.Error(),
		}
	}
	/* === END ===*/

	/*
		=== START ===
		Validate user cookies
	*/

	channelId, err = vs.UserRepository.GetYoutubeIdById(userId)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "failed to get youtube id",
			Data: err.Error(),
		}
	}
	/* === END ===*/

	token, err := vs.TokenRepository.GetByOwnerId(userId);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get token",
			Data: err.Error(),
		}
	}

	vidItem, err := vs.YtService.Video(token, videoId)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "failed to get videos",
			Data: err.Error(),
		}
	}

	if channelId != vidItem.Snippet.ChannelId {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "action not permitted, you are not the owner of the video",
			Data: nil,
		}
	}

	video.Id = vidItem.Id;
	video.Title = vidItem.Snippet.Title
	video.Owner = vidItem.Snippet.ChannelId
	video.Thumbnail = vidItem.Snippet.Thumbnails.Maxres.Url;
	video.Strategy = upload.Strategy;
	video.Scheduler = upload.Scheduler;

	tx, err := conn.Begin(context.Background());
	if err != nil {
		return dto.Response{
			Status: fiber.StatusInternalServerError,
			Message: "Failed enstablish database connection",
			Data: nil,
		}
	}
	defer tx.Rollback(context.Background());

	err = vs.VideoRepository.Save(tx, video)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to add a new videos",
			Data:err,
		}
	}

	var cooldown time.Time = time.Now().Add((time.Hour * 24) * 7);

	err = vs.UserRepository.UpdateVideos(tx, video.Id, userId, cooldown)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to update a new videos",
			Data:nil,
		}
	}
	
	var initialDetection dto.Response = vs.CommentService.FetchAndDeleteComment(cookie, videoId)
	
	if initialDetection.Status == fiber.StatusBadRequest {
		return initialDetection;
	}
	
	tx.Commit(context.Background());

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Successfully uploaded a new video",
		Data: video,
	}
}

func (vs *VideoServiceImpl) IsCanUpload(link string, cookie string) dto.Response {
	var channelId string;

	userId, err := vs.Authentication.GetUserIdByCookie(cookie);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get user information",
			Data: nil,
		}
	}

	isCooldown, err := vs.UserRepository.IsCooldown(userId)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed in check video cooldown",
			Data: err.Error(),
		}
	}

	if isCooldown {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "you only permitted upload video only once per week",
			Data: nil,
		}
	}


	channelId, err = vs.UserRepository.GetYoutubeIdById(userId);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get Channel Id",
			Data: nil,
		}
	}

	videoId, err := vs.processURL(link);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "failed to process URL",
			Data: err.Error(),
		}
	}

	
	token, err := vs.TokenRepository.GetByOwnerId(userId);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get token",
			Data: err.Error(),
		}
	}
	

	vidItem, err := vs.YtService.Video(token, videoId)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to fetch video",
			Data: err.Error(),
		}
	}

	if channelId != vidItem.Snippet.ChannelId {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "action not permitted, you are not the owner of the video",
			Data: nil,
		}
	}

	video := dto.Videos{
		Thumbnail: vidItem.Snippet.Thumbnails.Maxres.Url,
		Title: vidItem.Snippet.Title,
	}

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Video terverifikasi milik anda",
		Data: video,
	}

}

func (vs *VideoServiceImpl) processURL(link string) (string, error) {

	youtubeLink, err := url.Parse(link);

	if err != nil {
		return "", err;
	}

	queryLen := len(youtubeLink.Query()["v"])
	
	if queryLen == 0 {
		return "", errors.New("there is no Video ID")
	}

	var videoId string = youtubeLink.Query()["v"][0];

	return videoId, nil
}

func (vs *VideoServiceImpl) Info(videoId string, cookies string) dto.Response {
	userId, err := helper.VerifyAndGet(cookies);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusNonAuthoritativeInformation,
			Message: "Failed get user information",
			Data: nil,
		}
	}

	ownerId, err := vs.UserRepository.GetYoutubeIdById(userId);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get owner information",
			Data: nil,
		}
	}

	video, err := vs.VideoRepository.GetById(videoId)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get video information",
			Data: err,
		}
	}

	comments, err := vs.CommentRepository.GetByVideoId(videoId)

	video.Comments = comments;
	video.Id = videoId;

	if video.Owner != ownerId {
		return dto.Response{
			Status: fiber.StatusForbidden,
			Message: "Action is not permitted",
			Data: nil,
		}
	}

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Success get video information",
		Data: video,
	}

}
