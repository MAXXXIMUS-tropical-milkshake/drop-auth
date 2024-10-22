package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/core"
	"github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/lib/postgres"
	constraints "github.com/MAXXXIMUS-tropical-milkshake/beatflow-auth/internal/store/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type store struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) core.UserStore {
	return &store{pg}
}

func (s *store) GetUserByEmail(ctx context.Context, email string) (user *core.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user = new(core.User)

	stmt := `SELECT id, username, email, password_hash, is_deleted FROM users WHERE email = $1`
	err = s.DB.QueryRowContext(ctx, stmt, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *store) GetUserByUsername(ctx context.Context, username string) (user *core.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user = new(core.User)

	stmt := `SELECT id, username, email, password_hash, is_deleted FROM users WHERE username = $1`
	err = s.DB.QueryRowContext(ctx, stmt, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrInvalidCredentials
		}
		return nil, err
	}

	return user, nil
}

func (s *store) GetUserByID(ctx context.Context, userID int) (user *core.User, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user = new(core.User)

	stmt := `SELECT id, username, email, password_hash, is_deleted FROM users WHERE id = $1`
	err = s.DB.QueryRowContext(ctx, stmt, userID).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.IsDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *store) AddUser(ctx context.Context, user core.User) (userID int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `INSERT INTO users (username, email, password_hash)
	VALUES ($1, $2, $3) RETURNING id`

	err = s.DB.QueryRowContext(ctx, stmt, user.Username, user.Email, user.PasswordHash).Scan(&userID)
	if err != nil {
		var pg *pgconn.PgError
		if ok := errors.As(err, &pg); ok && pg.Code == pgerrcode.UniqueViolation {
			switch {
			case pg.ConstraintName == constraints.UniqueUsernameConstraint:
				return 0, core.ErrUsernameAlreadyExists
			case pg.ConstraintName == constraints.UniqueEmailConstraint:
				return 0, core.ErrEmailAlreadyExists
			default:
				return 0, core.ErrAlreadyExists
			}
		}
		return 0, err
	}

	return userID, err
}

func (s *store) UpdateUser(ctx context.Context, user core.UpdateUser) (userID int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var password *string
	if user.Password != nil {
		password = &user.Password.NewPassword
	}

	stmt := `UPDATE users SET
	password_hash = COALESCE($1, password_hash),
	username = COALESCE($2, username),
	email = COALESCE($3, email),
	updated_at = DEFAULT
	WHERE id = $4
	RETURNING id`
	err = s.DB.QueryRowContext(ctx, stmt, password, user.Username, user.Email, user.ID).Scan(&userID)
	if err != nil {
		var pg *pgconn.PgError
		if ok := errors.As(err, &pg); ok && pg.Code == pgerrcode.UniqueViolation {
			switch {
			case pg.ConstraintName == constraints.UniqueUsernameConstraint:
				return 0, core.ErrUsernameAlreadyExists
			case pg.ConstraintName == constraints.UniqueEmailConstraint:
				return 0, core.ErrEmailAlreadyExists
			default:
				return 0, core.ErrAlreadyExists
			}
		}
		return 0, err
	}

	return userID, nil
}

func (s *store) DeleteUser(ctx context.Context, userID int) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `UPDATE users SET
	is_deleted = true
	WHERE id = $1`

	err := s.DB.QueryRowContext(ctx, stmt, userID).Err()
	if err != nil {
		return err
	}

	return nil
}
