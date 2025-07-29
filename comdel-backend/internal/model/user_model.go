package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type User struct {
	UserId			string			`json:"userId" redis:"userId"`
	Name 			string			`json:"name" redis:"name"`
	GivenName		string			`json:"givenName" redis:"givenName"`
	Email			string			`json:"email" redis:"email"`
	VerifiedEmail	bool			`json:"verifiedEmail"`
	Subscription	string			`json:"subscription" redis:"subscription"`
	PremiumPlan		string			`json:"premiumPlan" redis:"premiumPlan"`
	Isverified		bool			`json:"isVerified"`
	Picture			string			`json:"picture" redis:"picture"`
	Videos			[]*Videos		`json:"videos"`
	GId				string			`json:"g_id" redis:"gid"`
	YoutubeId		string			`json:"youtubeId"`
	TitleSnippet	string			`json:"title_snippet"`
}

func (u *User) RedisHashString() []string {
	hashField := []string{
		"userId", u.UserId,
		"name", u.Name,
		"givenName", u.GivenName,
		"email", u.Email,
		"subscription", u.Subscription,
		"premiumPlan", u.PremiumPlan,
		"picture", u.Picture,
		"gid", u.GId,
	}

	return hashField
}

func (u *User) JSON() string {
	jsonBytes, err := json.Marshal(u)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}


func RedisUserKey(userId string) string {
	redisKey := os.Getenv("REDIS_USER_KEY")
	return fmt.Sprintf("%s:%s", redisKey, userId);
}

func RedisUserVideoKey(userId string) string {
	redisKey := os.Getenv("REDIS_VIDEO_USER_KEY")
	return fmt.Sprintf("%s:%s", redisKey, userId);
}