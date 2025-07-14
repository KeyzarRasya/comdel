package mock

import (
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

/*
	ChannelInfo(*oauth2.Token)		(*youtube.Channel, error)
	Video(*oauth2.Token, string)	(*youtube.Video, error)
*/

type MockYoutubeService struct {
	ChannelInfoFunc   func(*oauth2.Token) (*youtube.Channel, error)
	VideoFunc         func(*oauth2.Token, string) (*youtube.Video, error)
	CommentsFunc      func(token *oauth2.Token, videoId string) ([]*youtube.CommentThread, error)
	DeleteCommentFunc func(token *oauth2.Token, commentId string) error
}

// Comments implements services.YoutubeService.
func (mys MockYoutubeService) Comments(token *oauth2.Token, videoId string) ([]*youtube.CommentThread, error) {
	return mys.CommentsFunc(token, videoId)

}

// DeleteComment implements services.YoutubeService.
func (mys MockYoutubeService) DeleteComment(token *oauth2.Token, commentId string) error {
	return mys.DeleteCommentFunc(token, commentId)

}

// ChannelInfo implements services.YoutubeService.
func (mys MockYoutubeService) ChannelInfo(token *oauth2.Token) (*youtube.Channel, error) {
	return mys.ChannelInfoFunc(token)
}

// Video implements services.YoutubeService.
func (mys MockYoutubeService) Video(token *oauth2.Token, videoId string) (*youtube.Video, error) {
	return mys.VideoFunc(token, videoId)
}

