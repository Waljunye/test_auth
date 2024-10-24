package stores

import (
	"auth/internal/domain"
	"context"
	"database/sql"
	"github.com/stretchr/testify/suite"
)

type UsersStore_ByUid_TestSuite struct {
	suite.Suite
	InitialData []domain.User
	InputData   struct {
		Uuid string
	}
	Expected struct {
		User domain.User
		Err  error
	}
	Db  *sql.DB
	Ctx context.Context
}

func (s *UsersStore_ByUid_TestSuite) SetupTest() {
	stmt, _ := s.Db.PrepareContext(s.Ctx, `INSERT INTO users (uuid, username, password) VALUES ($1, $2, $3);`)
	defer stmt.Close()

	for _, initial := range s.InitialData {
		stmt.ExecContext(s.Ctx, initial.Uuid, initial.Username, initial.Password)
	}

	return
}

func (s *UsersStore_ByUid_TestSuite) TearDownTest() {
	stmt, _ := s.Db.PrepareContext(s.Ctx, `DELETE FROM users WHERE uuid IN($1);`)
	defer stmt.Close()

	uids := make([]string, 0, len(s.InitialData))
	for _, initial := range s.InitialData {
		uids = append(uids, initial.Uuid)
	}

	stmt.Exec(s.Ctx, stmt)

	return
}
