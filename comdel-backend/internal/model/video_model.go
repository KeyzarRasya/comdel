package model

type Videos struct {
	Id 			string		`json:"id"`
	Title 		string		`json:"title"`
	Thumbnail 	string		`json:"thumbnail"`
	Owner 		string		`json:"owner"`
	Strategy	string		`json:"strategy"`
	Scheduler	string		`json:"scheduler"`
	Comments	[]*Comment	`json:"comments"`
	DeletedComment	[]Comment	`json:"deletedComment"`
}