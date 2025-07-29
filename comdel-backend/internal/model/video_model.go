package model

import (
	"fmt"
	"os"
)

type Videos struct {
	Id 			string		`json:"id" redis:"id"`
	Title 		string		`json:"title" redis:"title"`
	Thumbnail 	string		`json:"thumbnail" redis:"thumbnail"`
	Owner 		string		`json:"owner" redis:"owner"`
	Strategy	string		`json:"strategy" redis:"strategy"`
	Scheduler	string		`json:"scheduler" redis:"scheduler"`
	Comments	[]*Comment	`json:"comments"`
	DeletedComment	[]Comment	`json:"deletedComment"`
}

func (v *Videos) RedisHashString() []string {
	hashVid := []string{
		"id", v.Id,
		"title", v.Title,
		"thumbnail", v.Thumbnail,
		"owner", v.Owner,
		"strategy", v.Strategy,
		"scheduler", v.Scheduler,
	}

	return hashVid;
}

func RedisVideoKey(videoId string) string {
	redisKey := os.Getenv("REDIS_VIDEO_KEY")
	return fmt.Sprintf("%s:%s", redisKey, videoId);
}