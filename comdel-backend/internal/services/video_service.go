package services

import (
	"context"
	"net/url"
	"time"

	"github.com/KeyzarRasya/comdel-server/internal/config"
	"github.com/KeyzarRasya/comdel-server/internal/dto"
	"github.com/KeyzarRasya/comdel-server/internal/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func CheckVideoOwnership(link string, cookie string) dto.Response {
	conn := config.LoadDatabase();
	oauthConfig := config.OAuthConfig();

	var token oauth2.Token;
	var channelId string;

	userId, err := helper.VerifyAndGet(cookie);

	isCooldown, err := CheckVideoCoolDown(conn, userId)

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

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get user information",
			Data: nil,
		}
	}

	err = conn.QueryRow(
		context.Background(),
		"SELECT youtube_id FROM user_info WHERE user_id=$1",
		userId,
	).Scan(&channelId);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get Channel Id",
			Data: nil,
		}
	}

	youtubeLink, err := url.Parse(link);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get link information",
			Data: nil,
		}
	}

	queryLen := len(youtubeLink.Query()["v"])
	
	if queryLen == 0 {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Cannot get video information. invalid videos id parameter",
			Data: nil,
		}
	}
	
	var videoId string = youtubeLink.Query()["v"][0];
	log.Info("Video ID: ", videoId)

	err = conn.QueryRow(
		context.Background(),
		"SELECT access_token, refresh_token, expiry FROM oauth_token WHERE owner_id=$1",
		userId,
	).Scan(&token.AccessToken, &token.RefreshToken, &token.Expiry)

	log.Info(token.AccessToken);

	var tokenSource = oauthConfig.TokenSource(context.Background(), &token);
	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(tokenSource));

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to create youtube service",
			Data: nil,
		}
	}

	videoResponse, err := youtubeService.Videos.List([]string{"snippet"}).Id(videoId).Do();

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get video informattion",
			Data: err,
		}
	}

	if len(videoResponse.Items) == 0 {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "invalid video id. no videos found",
			Data: nil,
		}
	}

	vidItem := videoResponse.Items[0];

	if channelId != vidItem.Snippet.ChannelId {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "action not permitted, you are not the owner of the video",
			Data: nil,
		}
	}

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Video terverifikasi milik anda",
		Data: nil,
	}

}

func CheckVideoCoolDown(conn *pgx.Conn, userID string) (bool, error) {
	var subscription string
	var cooldown *time.Time

	err := conn.QueryRow(
		context.Background(),
		`SELECT subscription, upload_cooldown FROM user_info WHERE user_id = $1`,
		userID,
	).Scan(&subscription, &cooldown)
	if err != nil {
		return false, err
	}

	if subscription == "NONE" || cooldown == nil || time.Now().After(*cooldown) {
		return false, nil
	}
	return true, nil
}



func GetVideoInformation(videoId string, cookies string) dto.Response {
	conn := config.LoadDatabase();

	var video dto.Videos;
	var ownerId string;
	var comments []dto.Comment;
	var videoComments []string;

	userId, err := helper.VerifyAndGet(cookies);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusNonAuthoritativeInformation,
			Message: "Failed get user information",
			Data: nil,
		}
	}
	
	err = conn.QueryRow(
		context.Background(),
		"SELECT youtube_id FROM user_info WHERE user_id=$1",
		userId,
	).Scan(&ownerId);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get owner information",
			Data: nil,
		}
	}

	err = conn.QueryRow(
		context.Background(),
		"SELECT title, thumbnail, strategy, scheduler, comments, owner FROM videos WHERE videos_id=$1",
		videoId,
	).Scan(&video.Title, &video.Thumbnail, &video.Strategy, &video.Scheduler, &videoComments, &video.Owner)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get video information",
			Data: err,
		}
	}

	var yCommentId, publishedAt, channelId, channelUrl, displayName, profileUrl, textDisplay string;
	var isDetected bool;

	rows, err := conn.Query(
		context.Background(),
		"SELECT DISTINCT y_comment_id, published_at, channel_id, author_channel_url, display_name, profile_url, text_display, isdetected FROM comments WHERE video_id=$1",
		videoId,
	)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "failed while fetch comments",
			Data: nil,
		}
	}

	_, err = pgx.ForEachRow(rows, []any{&yCommentId, &publishedAt, &channelId, &channelUrl, &displayName, &profileUrl, &textDisplay, &isDetected}, func () error  {
		comments = append(comments, dto.Comment{
			Yid: yCommentId,
			PublishedAt: publishedAt,
			ChannelId: channelId,
			ChannelUrl: channelUrl,
			DisplayName: displayName,
			ProfileUrl: profileUrl,
			TextDisplay: textDisplay,
			Isdetected: isDetected,
		})
		return nil
	})

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to appending comments",
			Data: err,
		}
	}

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
