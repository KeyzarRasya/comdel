package dto

import "github.com/KeyzarRasya/comdel-server/internal/model"


type GoogleProfile struct {
	GId				string 	`json:"id"`
	Email 			string	`json:"email"`;
	Name			string	`json:"name"`;
	VerifiedEmail	bool	`json:"verified_email"`
	GivenName		string	`json:"given_name"`
	Picture			string 	`json:"picture"`
	Token			string	`json:"token"`
}

func (gp *GoogleProfile) Parse() model.User {
	var user model.User;

	user.GId = gp.GId;
	user.Email = gp.Email;
	user.Name = gp.Name;
	user.GivenName = gp.GivenName;
	user.Picture = gp.Picture;
	user.VerifiedEmail = gp.VerifiedEmail

	return user;
}