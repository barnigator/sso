package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/barnigator/sso/internal/auth/deps"
	"github.com/barnigator/sso/internal/auth/domain"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "postgres.New"

	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: ping: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const fn = "repository.postgres.SaveUser"

	var id int64
	err := s.db.QueryRowContext(
		ctx,
		"INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id",
		email,
		passHash,
	).Scan(&id)

	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, deps.ErrUserExists
		}

		return 0, fmt.Errorf("%s: %w", fn, err)
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, email string) (domain.User, error) {
	const fn = "repository.postgres.GetUser"

	var user domain.User
	err := s.db.QueryRowContext(
		ctx,
		"SELECT id, email, pass_hash FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", fn, deps.ErrUserNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", fn, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const fn = "repository.postgres.IsAdmin"

	var IsAdmin bool
	err := s.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM admins WHERE user_id = $1)",
		userID,
	).Scan(&IsAdmin)
	if err != nil {

		return false, fmt.Errorf("%s: %w", fn, err)
	}

	return IsAdmin, nil
}

func (s *Storage) GetApp(ctx context.Context, appID int) (domain.App, error) {
	const fn = "repository.postgres.GetApp"

	var app domain.App

	err := s.db.QueryRowContext(
		ctx,
		"SELECT id, name, secret FROM apps WHERE id = $1",
		appID).Scan(&app.ID, &app.Name, &app.Secret)
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
