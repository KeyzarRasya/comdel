package services

import (
	"comdel-backend/internal/config"
	"comdel-backend/internal/dto"
	"comdel-backend/internal/helper"
	"comdel-backend/internal/inference"
	"comdel-backend/internal/repository"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type CommentService interface {
	FetchAndDeleteComment(string, string) dto.Response;
}

type CommentServiceImpl struct {
	UserRepository repository.UserRepository;
	TokenRepository repository.TokenRepository;
	YoutubeService YoutubeService
	CommentRepository repository.CommentRepository
	VideoRepository repository.VideoRepository
	OauthProvider config.OAuthProvider
	DBLoader config.DBLoader
}

func NewCommentService(
	userRepository repository.UserRepository,
	tokenRepository repository.TokenRepository,
	youtubeService YoutubeService,
	commentRepository repository.CommentRepository,
	videoRepository repository.VideoRepository,
	oauthProvider config.OAuthProvider,
	dbLoader config.DBLoader,
) CommentServiceImpl {
	return CommentServiceImpl{
		UserRepository: userRepository,
		TokenRepository: tokenRepository,
		YoutubeService: youtubeService,
		CommentRepository: commentRepository,
		VideoRepository: videoRepository,
		OauthProvider: oauthProvider,
		DBLoader: dbLoader,
	}
}

func (cs *CommentServiceImpl) FetchAndDeleteComment(cookie string, videoId string) dto.Response {
	conn, err := cs.DBLoader.Load()
	if err != nil {
		return dto.Response{
			Status: fiber.StatusInternalServerError,
			Message: "Failed to load database",
			Data: err.Error(),
		}
	}

	userId, err := helper.VerifyAndGet(cookie);

	var modelAPI inference.ModelAPI;
	var deletedComment []dto.Comment;
	var notDetectedComment []dto.Comment;

	var deletedCommentsId []string;
	var notDeletedCommentsId []string;

	channelId, err := cs.UserRepository.GetYoutubeIdById(userId);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get youtube id",
			Data: err.Error(),
		}
	}

	token, err := cs.TokenRepository.GetByOwnerId(userId)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get token",
			Data: err.Error(),
		}
	}

	vidItem, err := cs.YoutubeService.Video(token, videoId)
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get video informattion",
			Data: err,
		}
	}

	if channelId != vidItem.Snippet.ChannelId {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "action not permitted, you are not the owner of the video",
			Data: nil,
		}
	}
	
	comments, err := cs.YoutubeService.Comments(token, videoId);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get comment threads",
			Data: err,
		}
	}

	for _, comment := range comments {
		var commentObject dto.Comment;
		c := comment.Snippet.TopLevelComment.Snippet;

		commentObject.ChannelId = c.ChannelId;
		commentObject.ChannelUrl = c.AuthorChannelUrl;
		commentObject.DisplayName = c.AuthorDisplayName;
		commentObject.PublishedAt = c.PublishedAt;
		commentObject.ProfileUrl = c.AuthorProfileImageUrl;
		commentObject.TextDisplay = c.TextDisplay;
		commentObject.Yid = comment.Id;
		commentObject.VideoId = c.VideoId;

		result, msg, err := modelAPI.Detect(commentObject.TextDisplay).Get();

		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: msg,
				Data: err,
			}
		}

		if result == inference.NOT_DETECTED {
			notDetectedComment = append(notDetectedComment, commentObject);
		}

		if result == inference.DETECTED {
			err := cs.YoutubeService.DeleteComment(token, commentObject.Yid)
			if err != nil {
				return dto.Response{
					Status: fiber.StatusBadRequest,
					Message: "failed to delete comments",
					Data: err,
				}
			}

			log.Info("Comment deleted ", commentObject.TextDisplay);
			deletedComment = append(deletedComment, commentObject);
			log.Info(deletedComment);
		}

		if result == inference.ERROR {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: msg,
				Data: err,
			}
		}
	}

	tx, err := conn.Begin(context.Background());
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to start database interaction",
			Data: err,
		}
	}
	defer tx.Rollback(context.Background());

	for _, comment := range deletedComment {
		comment.Id, err = cs.CommentRepository.Save(tx, comment, true)
		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to update deleted comment",
				Data: err,
			}
		}
		deletedCommentsId = append(deletedCommentsId, comment.Yid);
	}

	for _, comment := range notDetectedComment {
		comment.Id, err = cs.CommentRepository.Save(tx, comment, false)
		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to update comment",
				Data: err,
			}
		}
		notDeletedCommentsId = append(notDeletedCommentsId, comment.Yid)
	}

	err = cs.VideoRepository.UpdateComment(tx, deletedCommentsId, videoId, true)
	if err != nil {
		log.Info(err);
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to update deleted comment information",
			Data: err,
		}
	}

	err = cs.VideoRepository.UpdateComment(tx, notDeletedCommentsId, videoId, false)
	if err != nil {
		log.Info(err);
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to update video information",
			Data: err,
		}
	}

	tx.Commit(context.Background());

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Success fetch and upload new comment",
		Data: nil,
	}
}