package model

type User struct {
	UserId			string			`json:"userId"`
	Name 			string			`json:"name"`
	GivenName		string			`json:"givenName"`
	Email			string			`json:"email"`
	VerifiedEmail	bool			`json:"verifiedEmail"`
	Subscription	string			`json:"subscription"`
	PremiumPlan		string			`json:"premiumPlan"`
	Isverified		bool			`json:"isVerified"`
	Picture			string			`json:"picture"`
	Videos			[]*Videos		`json:"videos"`
	GId				string			`json:"g_id"`
	YoutubeId		string			`json:"youtubeId"`
	TitleSnippet	string			`json:"title_snippet"`
}