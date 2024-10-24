package services

import (
	"context"
	"crypto/rsa"
	"time"

	"auth/internal/domain"
)

type usersStore interface {
	Create(ctx context.Context, user domain.User) (err error)
	ByUsername(ctx context.Context, username string) (user *domain.User, err error)
	ByUid(ctx context.Context, uid string) (user *domain.User, err error)
}

type authServiceConfig interface {
	PrivateKey() *rsa.PrivateKey
	PublicKey() *rsa.PublicKey
	HashKey() string
	Salt() int
	RefreshTokenExpires() time.Duration
	AccessTokenExpires() time.Duration
}
