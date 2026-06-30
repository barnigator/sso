package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/barnigator/sso/internal/auth/deps"
	domain "github.com/barnigator/sso/internal/auth/domain"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const fn = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", fn, deps.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", fn, err)

	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (domain.User, error) {
	const fn = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user domain.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", fn, deps.ErrUserNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const fn = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", fn, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool
	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", fn, deps.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (domain.App, error) {
	const fn = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return domain.App{}, fmt.Errorf("%s: %w", fn, err)
	}

	row := stmt.QueryRowContext(ctx, appID)
	var app domain.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.App{}, deps.ErrAppNotFound
		}

		return domain.App{}, fmt.Errorf("%s: %w", fn, err)
	}

	return app, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
