package v1

import (
	"auth/internal/domain"
	"context"
)

type authService interface {
	SignUp(ctx context.Context, username, password string) (err error)
	SignIn(ctx context.Context, username, password string) (refreshToken, accessToken string, user *domain.User, err error)
	Refresh(ctx context.Context, oldRefreshToken string) (accessToken, refreshToken string, err error)
}
