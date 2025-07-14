package model


type Comment struct {
	Id 				string		`json:"id"`
	Yid 			string		`json:"youtubeCommentId"`
	PublishedAt		string		`json:"publishedAt"`
	ChannelId		string		`json:"channelId"`
	ChannelUrl		string		`json:"channelUrl"`
	DisplayName		string		`json:"displayName"`
	ProfileUrl		string		`json:"profileUrl"`
	TextDisplay		string		`json:"textDisplay"`
	VideoId 		string		`json:"videoId"`
	Isdetected		bool		`json:"isDetected"`
}