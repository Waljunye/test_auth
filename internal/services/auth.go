package services

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"

	"auth/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	cfg        authServiceConfig
	usersStore usersStore
}

func NewAuthService(us usersStore, cfg authServiceConfig) *AuthService {
	return &AuthService{cfg: cfg, usersStore: us}
}

func (as *AuthService) SignUp(ctx context.Context, username, password string) (err error) {
	alreadyExist, err := as.usersStore.ByUsername(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return
	}
	if alreadyExist != nil {
		err = errors.New("user already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), as.cfg.Salt())
	if err != nil {
		return err
	}
	err = as.usersStore.Create(ctx, domain.User{
		Uuid:     uuid.New().String(),
		Username: username,
		Password: string(hashedPassword),
	})
	log.Info().Str("username", username).Msg("user signed up")

	if err != nil {
		return
	}

	return
}

func (as *AuthService) SignIn(ctx context.Context, username, password string) (refreshToken, accessToken string, user *domain.User, err error) {
	user, err = as.usersStore.ByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrSignIn{}
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			err = ErrSignIn{}
		}
		return
	}
	log.Info().Str("username", user.Username).Msg("user signed in")

	refreshToken, err = as.generateToken(as.cfg.PrivateKey(), user.Username, as.cfg.RefreshTokenExpires())
	if err != nil {
		return
	}
	accessToken, err = as.generateToken(as.cfg.PrivateKey(), user.Username, as.cfg.AccessTokenExpires())
	if err != nil {
		return
	}
	return
}

func (as *AuthService) Refresh(ctx context.Context, oldRefreshToken string) (accessToken, refreshToken string, err error) {
	tokenParsed, err := jwt.Parse(oldRefreshToken, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodRS512.Alg() {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return as.cfg.PublicKey(), nil
	})
	if err != nil {
		return
	}
	var (
		claims jwt.MapClaims
		ok     bool
	)
	if claims, ok = tokenParsed.Claims.(jwt.MapClaims); ok {
		err = errors.New("no claims in token or its corrupted")
		return
	}

	username := claims["sub"].(string)
	log.Info().Str("username", username).Msg("user refreshed")

	user, err := as.usersStore.ByUsername(ctx, username)
	if err != nil || user == nil {
		if errors.Is(err, sql.ErrNoRows) || user == nil {
			err = errors.New("username that claims in subject does not exist")
			log.Warn().Err(err).Str("username", username).Msg("PRIVATE KEY LEAK. SHUTTING DOWN SERVER")
			panic(errors.Wrap(err, "PRIVATE KEY LEAK. SHUTTING DOWN SERVER"))
		}
		return
	}

	refreshToken, err = as.generateToken(as.cfg.PrivateKey(), user.Username, as.cfg.RefreshTokenExpires())
	if err != nil {
		return
	}

	accessToken, err = as.generateToken(as.cfg.PrivateKey(), user.Username, as.cfg.AccessTokenExpires())
	if err != nil {
		return
	}

	return
}

func (as *AuthService) generateToken(key *rsa.PrivateKey, subject string, expires time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
	})

	return token.SignedString(key)
}
