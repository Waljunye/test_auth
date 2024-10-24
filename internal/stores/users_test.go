package stores

import (
	"context"
	"testing"

	"auth/testdata/generators"

	"github.com/stretchr/testify/assert"
)

func TestUsersStore(t *testing.T) {
	expectedUser := generators.RandomUser()
	dbConn.Exec("TRUNCATE TABLE users CASCADE")

	t.Run("Create(ok)", func(t *testing.T) {
		err := usersStore.Create(context.Background(), expectedUser)
		assert.NoError(t, err)
	})

	t.Run("ByUUID(ok)", func(t *testing.T) {
		actualUser, err := usersStore.ByUid(context.Background(), expectedUser.Uuid)
		assert.NoError(t, err)

		assert.Equal(t, expectedUser, actualUser)
	})

	t.Run("ByUsername(ok)", func(t *testing.T) {
		actualUser, err := usersStore.ByUsername(context.Background(), expectedUser.Username)
		assert.NoError(t, err)

		assert.Equal(t, expectedUser, actualUser)
	})

}
