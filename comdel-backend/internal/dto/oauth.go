package dto

type GoogleProfile struct {
	GId				string 	`json:"id"`
	Email 			string	`json:"email"`;
	Name			string	`json:"name"`;
	VerifiedEmail	bool	`json:"verified_email"`
	GivenName		string	`json:"given_name"`
	Picture			string 	`json:"picture"`
	Token			string	`json:"token"`
}