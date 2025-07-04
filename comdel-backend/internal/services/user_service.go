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

/*

	this function is used to store user information
	@params user dto.GoogleProfile
	@params oauthToken *oauth2.Token

	return value = dto.Response

*/
func SaveUser(user dto.GoogleProfile, oauthToken *oauth2.Token) dto.Response {
	conn := config.LoadDatabase();		/*Load Database*/
	oauthConfig := config.OAuthConfig()	/*Load OAuth Config*/

	var googleId string;		/*value to store g_id (it is available or not)*/
	var userId string;

	tx, err := conn.Begin(context.Background());	/* Starting database transaction */

	/*
		Check if transaction is successfully created or not,
		if not success, returning dto.Response
	*/
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "failed to start transaction",
			Data: nil,
		}
	}

	/*
		making token source for created youtube service
		and returning dto.Response if opeation failed
	*/
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

	err = tx.QueryRow(context.Background(), "SELECT g_id FROM user_info WHERE g_id=$1", user.GId).Scan(&googleId);
	
	if err != nil && err != pgx.ErrNoRows {
		tx.Rollback(context.Background())
		return dto.Response{
			Status:  fiber.StatusBadRequest,
			Message: "Database error checking g_id",
			Data:    err.Error(),
		}
	}

	if googleId == "" {
		log.Info("Creating a new account");
		err := tx.QueryRow(
			context.Background(), 
			"INSERT INTO user_info (name, given_name, email, isverified, picture, g_id, youtube_id, title_snippet) VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING user_id",
			user.Name, user.GivenName, user.Email, user.VerifiedEmail, user.Picture, user.GId, channel.Items[0].Id, channel.Items[0].Snippet.Title,
		).Scan(&userId)
	
		if err != nil {
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "failed to insert a profile",
				Data: err.Error(),
			}
		}

		log.Info("Refresh Token", oauthToken.RefreshToken);
		_, err = tx.Exec(
			context.Background(),
			"INSERT INTO oauth_token (access_token, refresh_token, expiry, owner_id) VALUES ($1, $2, $3, $4)",
			oauthToken.AccessToken, oauthToken.RefreshToken, oauthToken.Expiry, userId,
		)

		if err != nil {
			log.Info(err);
			return dto.Response{
				Status: fiber.StatusBadRequest,
				Message: "Failed to save token",
				Data: err,
			}
		}
	}

	defer tx.Rollback(context.Background());

	err = tx.QueryRow(
		context.Background(),
		"SELECT user_id FROM user_info WHERE g_id=$1",
		googleId,
	).Scan(&userId);

	jwt, err := helper.GenerateToken(userId);
	log.Info(userId)

	if err != nil{
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to generate JWT Token for authentication",
			Data: nil,
		}
	}

	user.Token = jwt;

	if err := GrantSubscriptionAccess(tx, userId); err != nil {
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


func GetUser(cookies string) dto.Response {

	conn := config.LoadDatabase();

	var user dto.User;
	var videosId []string;
	var videos []dto.Videos;
	
	if cookies == "" {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get credentials informations",
			Data: nil,
		}
	}


	token, err := helper.VerifyToken(cookies);

	
	if err != nil {
		log.Info("Error : ", err);
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Cannot verify JWT Token",
			Data: err,
		}
	}

	err = conn.QueryRow(
		context.Background(),
		"SELECT user_id, name, email, given_name, picture, videos, g_id, subscription, premium_plan FROM user_info WHERE user_id=$1",
		token["userId"],
	).Scan(&user.UserId,&user.Name, &user.Email, &user.GivenName, &user.Picture, &videosId, &user.GId, &user.Subscription, &user.PremiumPlan)

	log.Info("Premium Plan ", user.PremiumPlan)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get user information",
			Data: err,
		}
	}


	for _, id := range videosId {
		var video dto.Videos;

		err = conn.QueryRow(
			context.Background(),
			"SELECT title, thumbnail, strategy, scheduler, owner FROM videos WHERE videos_id=$1",
			id,
		).Scan(&video.Title, &video.Thumbnail, &video.Strategy, &video.Scheduler, &video.Owner);

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

func UploadVideo(cookie string, upload dto.UploadVideos) dto.Response {

	if upload.Link == "" || upload.Scheduler == "" || upload.Strategy == "" {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "invalid request query",
			Data: nil,
		}
	}

	/* Variable Declaration */
	conn := config.LoadDatabase();
	oauthConfig := config.OAuthConfig();
	var channelId string;
	var video dto.Videos;
	
	/* 
		=== START ===
		Parsing Url End Get the value 
	*/
	youtubeLink, err := url.Parse(upload.Link);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to parse youtube link",
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
	/* === END ===*/



	/*
		=== START ===
		Validate user cookies
	*/
	userId, err := helper.VerifyAndGet(cookie);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "invalid authentication token",
			Data: nil,
		}
	}

	err = conn.QueryRow(
		context.Background(),
		"SELECT youtube_id FROM user_info WHERE user_id=$1",
		userId,
	).Scan(&channelId);

	var token oauth2.Token

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: err.Error(),
			Data: nil,
		}
	}
	/* === END ===*/


	err = conn.QueryRow(
			context.Background(),
			"SELECT access_token, refresh_token, expiry FROM oauth_token WHERE owner_id=$1",
			userId,
		).Scan(&token.AccessToken, &token.RefreshToken, &token.Expiry)

	var tokenSource = oauthConfig.TokenSource(context.Background(), &token);
	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(tokenSource));

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to authenticate to youtube",
			Data: err,
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

	

	_, err = tx.Exec(
		context.Background(),
		"INSERT INTO videos(videos_id, title, owner, thumbnail, strategy, scheduler) VALUES ($1, $2, $3, $4, $5, $6)",
		video.Id, video.Title, video.Owner, video.Thumbnail, upload.Strategy, upload.Scheduler,
	)

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to add a new videos",
			Data:err,
		}
	}

	var cooldown time.Time = time.Now().Add((time.Hour * 24) * 7);

	_, err = tx.Exec(
		context.Background(),
		"UPDATE user_info SET videos = array_append(videos, $1), upload_cooldown=$2 WHERE user_id=$3",
		video.Id, cooldown, userId,
	)

	defer tx.Rollback(context.Background());

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to update a new videos",
			Data:nil,
		}
	}
	
	var initialDetection dto.Response = FetchAndDeleteComment(cookie, videoId)
	
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

func GetComments(cookie string, link string) dto.Response {
	conn := config.LoadDatabase();
	oauthConfig := config.OAuthConfig()

	var channelId string;
	
	cookieToken, err := helper.VerifyToken(cookie);
	
	if err != nil {
		return dto.Response{	
			Status: fiber.StatusBadRequest,
			Message: "Failed to verify jwt token",
			Data: nil,
		}
	}

	userId, ok := cookieToken["userId"].(string)
	if !ok || userId == "" {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Invalid or empty userId in token",
			Data: nil,
		}
	}



	var token oauth2.Token;

	log.Info("User Id : ", userId);
	
	err = conn.QueryRow(
		context.Background(),
		"SELECT access_token, refresh_token, expiry FROM oauth_token WHERE owner_id=$1",
		userId,
		).Scan(&token.AccessToken, &token.RefreshToken, &token.Expiry)
		
	if err != nil {
		log.Info("Error : ", err);
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to fetch oauth token",
			Data: err,
		}
	}
		
	tokenSource := oauthConfig.TokenSource(context.Background(), &token)
	
	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(tokenSource))
	
	parsedUrl, err := url.Parse(link);
	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to parse youtube url",
			Data: nil,
		}
	}
	
	var youtubeId string = parsedUrl.Query()["v"][0];
	log.Info("Youtube ID : ", youtubeId);
	videoResponse, err := youtubeService.Videos.List([]string{"snippet"}).Id(youtubeId).Do();

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get videos data",
			Data: err,
		}
	}

	log.Info("Video Response: ",videoResponse);
	if len(videoResponse.Items) == 0 {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "failed to get videos information, no data accepted",
			Data: nil,
		}
	}
	
	// var videoOwner string = videoResponse.Items[0].Snippet.ChannelId;


	err = conn.QueryRow(
		context.Background(),
		"SELECT youtube_id FROM user_info WHERE user_id=$1",
		userId,
	).Scan(&channelId);

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to retrieve user youtube information",
			Data: err,
		}
	}

	// if channelId != videoOwner {
	// 	return dto.Response{
	// 		Status: fiber.StatusBadRequest,
	// 		Message: "Illegal action, action is not permitted. you are not the owner of the video",
	// 		Data: nil,
	// 	}
	// }

	
	comments, err := youtubeService.CommentThreads.List([]string{"snippet"}).VideoId(youtubeId).MaxResults(10).Do();

	if err != nil {
		return dto.Response{
			Status: fiber.StatusBadRequest,
			Message: "Failed to get comment threads",
			Data: err,
		}
	}

	return dto.Response{
		Status: fiber.StatusOK,
		Message: "Success Fetch comment",
		Data: comments,
	}
}