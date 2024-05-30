package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/belimang/internal/entity"
	"github.com/sprint-id/belimang/internal/ierr"
)

type userRepo struct {
	conn *pgxpool.Pool
}

func newUserRepo(conn *pgxpool.Pool) *userRepo {
	return &userRepo{conn}
}

func (u *userRepo) Insert(ctx context.Context, user entity.User) (string, error) {
	const query = `INSERT INTO users (id, username, email, password, is_admin, created_at)
                   VALUES (gen_random_uuid(), $1, $2, $3, $4, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	var userID string
	err := u.conn.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.IsAdmin).Scan(&userID)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return "", ierr.ErrDuplicate
		}
		return "", err
	}

	return userID, nil
}

func (u *userRepo) GetByUsername(ctx context.Context, cred string) (entity.User, error) {
	user := entity.User{}
	q := `SELECT id, username, email, password, is_admin FROM users
	WHERE username = $1`

	err := u.conn.QueryRow(ctx,
		q, cred).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsAdmin)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return user, ierr.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (u *userRepo) IsAdmin(ctx context.Context, sub string) (bool, error) {
	var isAdmin bool
	q := `SELECT is_admin FROM users WHERE id = $1`

	err := u.conn.QueryRow(ctx, q, sub).Scan(&isAdmin)
	if err != nil {
		return false, err
	}

	return isAdmin, nil
}
