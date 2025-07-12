package mock

type MockAuthenticator struct {
	GenerateTokenFunc 		func(userId string) 	(string, error);
	VerfiyFunc				func(token string)		(map[string]interface{}, error)
	GetUserIdByCookieFunc	func(cookie string)		(string, error)
}

func (ma *MockAuthenticator) GenerateToken(userId string) (string, error) {
	return ma.GenerateTokenFunc(userId) 
}

func (ma *MockAuthenticator) Verify(token string) (map[string]interface{}, error) {
	return ma.VerfiyFunc(token);
}

func (ma *MockAuthenticator) GetUserIdByCookie(cookie string) (string, error) {
	return ma.GetUserIdByCookieFunc(cookie);
}