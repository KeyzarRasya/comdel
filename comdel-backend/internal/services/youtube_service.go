package services

import (
	"comdel-backend/internal/config"
	"context"
	"errors"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeService interface {
	ChannelInfo(*oauth2.Token)		(*youtube.Channel, error)
	Video(*oauth2.Token, string)	(*youtube.Video, error)
	Comments(token *oauth2.Token, videoId string)  ([]*youtube.CommentThread,error)
	DeleteComment(token *oauth2.Token, commentId string) error;
}

type YoutubeServiceImpl struct {
	OAuthProvider config.OAuthProvider
}

func NewYoutubeService(oauthProvider config.OAuthProvider) YoutubeServiceImpl {
	return YoutubeServiceImpl{OAuthProvider: oauthProvider}
}

func (ys *YoutubeServiceImpl) ChannelInfo(oauthToken *oauth2.Token) (*youtube.Channel, error) {
	tokenSource := ys.OAuthProvider.TokenSource(context.Background(), oauthToken);
	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(tokenSource));
	if err != nil {
		return nil, errors.New("Failed to create Youtube Service");
	}

	channel, err := youtubeService.Channels.List([]string{"id", "snippet"}).Mine(true).Do();
	if err != nil {
		return nil, errors.New("Failed to get channel list");
	}

	return channel.Items[0], err
}

func (ys *YoutubeServiceImpl) Video(oauthToken *oauth2.Token, videoId string) (*youtube.Video, error) {
	tokenSource := ys.OAuthProvider.TokenSource(context.Background(), oauthToken)
	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(tokenSource));
	if err != nil {
		return nil, errors.New("Failed to create Youtube Service");
	}

	videoResponse, err := youtubeService.Videos.List([]string{"snippet"}).Id(videoId).Do();
	if err != nil {
		return nil, err;
	}

	if len(videoResponse.Items) == 0 {
		return nil, errors.New("Invalid video id, no videos were found")
	}

	return videoResponse.Items[0], nil;
}

func (ys *YoutubeServiceImpl) Comments(token *oauth2.Token, videoId string) ([]*youtube.CommentThread,error) {
	tokenSource := ys.OAuthProvider.TokenSource(context.Background(), token)
	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(tokenSource));
	if err != nil {
		return nil, errors.New("Failed to create Youtube Service");
	}

	commentThreads, err := youtubeService.CommentThreads.List([]string{"snippet"}).VideoId(videoId).Do();
	if err != nil {
		return nil, err
	}

	return commentThreads.Items, nil
}

func (ys *YoutubeServiceImpl) DeleteComment(token *oauth2.Token, commentId string) error {
	tokenSource := ys.OAuthProvider.TokenSource(context.Background(), token)
	youtubeService, err := youtube.NewService(context.Background(), option.WithTokenSource(tokenSource));
	if err != nil {
		return errors.New("Failed to create Youtube Service");
	}
	return youtubeService.Comments.Delete(commentId).Do();
}