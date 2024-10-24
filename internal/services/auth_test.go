package services

import (
	"auth/internal/domain"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func authConfigConfigureMock(ctrl *gomock.Controller) (result authServiceConfig) {
	cfg := NewMockauthServiceConfig(ctrl)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	// Получение публичного ключа из приватного
	publicKey := &privateKey.PublicKey

	cfg.EXPECT().PrivateKey().AnyTimes().Return(privateKey)
	cfg.EXPECT().PublicKey().AnyTimes().Return(publicKey)
	cfg.EXPECT().Salt().AnyTimes().Return(bcrypt.DefaultCost)
	cfg.EXPECT().HashKey().AnyTimes().Return("test")
	cfg.EXPECT().RefreshTokenExpires().AnyTimes().Return(time.Minute * 2)
	cfg.EXPECT().AccessTokenExpires().AnyTimes().Return(time.Minute)

	result = cfg
	return
}

func TestAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	cfg := authConfigConfigureMock(ctrl)
	mockedStoreUsers := NewMockusersStore(ctrl)

	usrsByUuid := make(map[string]domain.User)
	mockedStoreUsers.EXPECT().Create(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(_ context.Context, user domain.User) (err error) {
		usrsByUuid[user.Uuid] = user
		return
	})
	mockedStoreUsers.EXPECT().ByUsername(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(_ context.Context, username string) (user *domain.User, err error) {
		founded := false
		for _, usr := range usrsByUuid {
			if usr.Username == username {
				user = &usr
				founded = true
			}
		}
		if !founded {
			err = sql.ErrNoRows
		}
		return
	})
	mockedStoreUsers.EXPECT().ByUid(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(func(_ context.Context, uid string) (user domain.User, err error) {
		var ok bool
		if user, ok = usrsByUuid[uid]; ok {
			return
		}
		err = sql.ErrNoRows
		return
	})

	t.Run("works correctly", func(t *testing.T) {
		serviceAuth := NewAuthService(mockedStoreUsers, cfg)

		err := serviceAuth.SignUp(context.Background(), "test", "test")
		assert.NoError(t, err)

		at, rt, err := serviceAuth.SignIn(context.Background(), "test", "test")
		assert.NoError(t, err)

		assert.NotEmpty(t, rt)
		assert.NotEmpty(t, at)

		at, rt, err = serviceAuth.Refresh(context.Background(), rt)
	})
}
