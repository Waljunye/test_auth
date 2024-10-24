package dto

import "auth/internal/domain"

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	User     User   `json:"user"`
}

type SignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

func BuildSignInResponse(accessToken, refreshToken string, user domain.User) *SignInResponse {
	return &SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: User{
			Uuid: user.Uuid,
			Name: user.Username,
		},
	}
}

type User struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type SignUpResponse struct{}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}
