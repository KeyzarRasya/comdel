package services

import (
	"comdel-backend/internal/dto"
	"comdel-backend/internal/inference"
	"context"
	"errors"
	"fmt"

	"comdel-backend/internal/config"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jackc/pgx/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)



func CronFetchDelete() error {
	conn := config.LoadDatabase()
	oauthConfig := config.OAuthConfig()
	modelAPI := inference.ModelAPI{}

	type TokenData struct {
		Token  oauth2.Token
		UserID string
	}
	var tokenList []TokenData

	rows, err := conn.Query(
		context.Background(),
		"SELECT access_token, refresh_token, expiry, owner_id FROM oauth_token",
	)
	if err != nil {
		log.Info(fmt.Sprintf("DB query error in CronFetchDelete: %s", err))
		return errors.New("Failed to get Oauth Token")
	}
	defer rows.Close()

	// Drain rows first
	for rows.Next() {
		var token oauth2.Token
		var userId string
		if err := rows.Scan(&token.AccessToken, &token.RefreshToken, &token.Expiry, &userId); err != nil {
			log.Info("Failed to scan oauth token row: " + err.Error())
			continue
		}
		tokenList = append(tokenList, TokenData{Token: token, UserID: userId})
	}

	// Now process each token-user pair after rows are closed
	for _, t := range tokenList {
		if err := processUser(conn, &oauthConfig, modelAPI, t.Token, t.UserID); err != nil {
			log.Info("Error processing user: " + err.Error())
			// Optionally continue or return based on needs
		}
	}

	return nil
}

func processUser(conn *pgx.Conn, oauthConfig *oauth2.Config, modelAPI inference.ModelAPI, token oauth2.Token, userId string) error {
	videoIDs, err := fetchUserVideos(conn, userId)
	if err != nil {
		log.Info(err.Error())
		return errors.New("Failed to fetch user videos")
	}

	youtubeService, err := youtube.NewService(
		context.Background(),
		option.WithTokenSource(oauthConfig.TokenSource(context.Background(), &token)),
	)
	if err != nil {
		log.Info("Failed to create youtube service")
		return errors.New("Failed to get youtube service")
	}

	for _, videoID := range videoIDs {
		if err := processVideo(conn, youtubeService, modelAPI, videoID); err != nil {
			log.Info("Failed to process user videos")
			return errors.New("Failed to process user videos")
		}
	}
	return nil
}

func processVideo(conn *pgx.Conn, yt *youtube.Service, modelAPI inference.ModelAPI, videoID string) error {
	comments, err := fetchComments(yt, videoID)
	if err != nil {
		return errors.New("Failed to fetch commen")
	}

	var deletedComments, notDetectedComments []dto.Comment

	for _, comment := range comments {
		result, _, err := modelAPI.Detect(comment.TextDisplay).Get()
		if err != nil {
			log.Info("Failed to detect model API")
			return errors.New("Failed to detech model api")
		}

		log.Info(result)
		log.Info(comment.TextDisplay)
		log.Info("================")

		switch result {
		case inference.DETECTED:
			if err := deleteComment(yt, comment); err != nil {
				log.Info("Failed to delete comment")
				return errors.New("Failed to delete comment")
			}
			deletedComments = append(deletedComments, comment)
		case inference.NOT_DETECTED:
			notDetectedComments = append(notDetectedComments, comment)
		case inference.ERROR:
			return errors.New("Model API Error")
		}
	}

	return saveResultsToDB(conn, videoID, deletedComments, notDetectedComments)
}

func saveResultsToDB(conn *pgx.Conn, videoID string, detected, undetected []dto.Comment) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Info("Failed to create transaction")
		return errors.New("Failed to create transaction db")
	}
	defer tx.Rollback(context.Background())

	var detectedIDs, undetectedIDs []string

	for _, c := range detected {
		if err := insertComment(tx, c, true); err != nil {
			log.Info("Failed to insert comment")
			return errors.New("Failed to inser comment")
		}
		detectedIDs = append(detectedIDs, c.Yid)
	}
	for _, c := range undetected {
		if err := insertComment(tx, c, false); err != nil {
			log.Info("Failed to insert undetected comment")
			return errors.New("Failed to insert undetected comment")
		}
		undetectedIDs = append(undetectedIDs, c.Yid)
	}

	if _, err := tx.Exec(
		context.Background(),
		"UPDATE videos SET deleted_comments = $1, comments = $2 WHERE videos_id = $3",
		detectedIDs, undetectedIDs, videoID,
	); err != nil {
		log.Info("update video information")
		return errors.New("Failed to update video info")
	}

	return tx.Commit(context.Background())
}

func fetchUserVideos(conn *pgx.Conn, userId string) ([]string, error) {
	var videoIDs []string
	err := conn.QueryRow(
		context.Background(),
		"SELECT videos FROM user_info WHERE user_id = $1",
		userId,
	).Scan(&videoIDs)
	return videoIDs, err
}

func fetchComments(yt *youtube.Service, videoID string) ([]dto.Comment, error) {
	res, err := yt.CommentThreads.List([]string{"snippet"}).VideoId(videoID).Do()
	if err != nil {
		log.Info("Failed to get youtube comment threads")
		return nil, errors.New("Failed to get comment threads")
	}
	var comments []dto.Comment
	for _, item := range res.Items {
		snippet := item.Snippet.TopLevelComment.Snippet
		comments = append(comments, dto.Comment{
			Yid:         item.Id,
			ChannelId:   snippet.ChannelId,
			ChannelUrl:  snippet.AuthorChannelUrl,
			DisplayName: snippet.AuthorDisplayName,
			PublishedAt: snippet.PublishedAt,
			ProfileUrl:  snippet.AuthorProfileImageUrl,
			TextDisplay: snippet.TextDisplay,
			VideoId:     snippet.VideoId,
		})
	}
	return comments, nil
}

func deleteComment(yt *youtube.Service, comment dto.Comment) error {
	return yt.Comments.Delete(comment.Yid).Do()
}

func insertComment(tx pgx.Tx, comment dto.Comment, detected bool) error {
	return tx.QueryRow(
		context.Background(),
		"INSERT INTO comments(y_comment_id, published_at, channel_id, author_channel_url, display_name, profile_url, text_display, video_id, isdetected) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING comment_id",
		comment.Yid, comment.PublishedAt, comment.ChannelId, comment.ChannelUrl, comment.DisplayName, comment.ProfileUrl, comment.TextDisplay, comment.VideoId, detected,
	).Scan(&comment.Id)
}


