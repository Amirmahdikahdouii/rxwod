package auth

type RegisterCommand struct {
	Email       string
	Password    string
	DisplayName string
}

type LoginCommand struct {
	Email    string
	Password string
}

type TokenDTO struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type AccessTokenDTO struct {
	AccessToken string
	ExpiresIn   int64
}

type UserDTO struct {
	ID          string
	Email       string
	DisplayName string
}
