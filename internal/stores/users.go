package stores

import (
	"auth/internal/domain"
	"auth/libs/pg"
	"context"
	"database/sql"
)

type UsersStore struct {
	db *sql.DB
}

func NewUsersStore(db *sql.DB) *UsersStore {
	return &UsersStore{
		db: db,
	}
}

func (us *UsersStore) ByUsername(ctx context.Context, username string) (user *domain.User, err error) {
	stmt, err := us.db.PrepareContext(ctx, `SELECT uuid, username, password FROM users WHERE username = $1`)
	if err != nil {
		pg.CloseStmt(stmt)
		return
	}
	defer pg.CloseStmt(stmt)

	dto := DTOUser{}
	row := stmt.QueryRowContext(ctx, username)
	err = row.Scan(&dto.Uuid, &dto.Username, &dto.Password)
	if err != nil {
		return
	}

	user = &domain.User{
		Uuid:     dto.Uuid,
		Username: dto.Username,
		Password: dto.Password,
	}
	return
}

func (us *UsersStore) Create(ctx context.Context, user domain.User) (err error) {
	stmt, err := us.db.PrepareContext(ctx, `INSERT INTO users (uuid, username, password) VALUES ($1, $2, $3)`)
	if err != nil {
		pg.CloseStmt(stmt)
		return
	}
	defer pg.CloseStmt(stmt)

	var ra int64
	res, err := stmt.ExecContext(ctx, user.Uuid, user.Username, user.Password)
	if err != nil {
		return
	}
	ra, err = res.RowsAffected()
	if err != nil {
		return
	}
	if ra == 0 {
		err = NewErrNoRowsWasAffected("\"users.Create\"")
		return
	}

	return
}
func (us *UsersStore) ByUid(ctx context.Context, uid string) (user *domain.User, err error) {
	stmt, err := us.db.PrepareContext(ctx, `SELECT uuid, username, password FROM users WHERE uuid = $1`)
	if err != nil {
		pg.CloseStmt(stmt)
		return
	}
	defer pg.CloseStmt(stmt)

	dto := DTOUser{}
	row := stmt.QueryRowContext(ctx, uid)
	err = row.Scan(&dto.Uuid, &dto.Username, &dto.Password)
	if err != nil {
		return
	}

	user = &domain.User{
		Uuid:     dto.Uuid,
		Username: dto.Username,
		Password: dto.Password,
	}
	return
}

type DTOUser struct {
	Uuid     string
	Username string
	Password string
}
