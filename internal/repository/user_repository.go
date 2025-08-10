package repository

import (
	"context"
	"time"

	"github.com/LuizFernando991/golang-auth-microservice/internal/model"
	"github.com/jmoiron/sqlx"
)

type RefreshTokenRow struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id int64) (*model.User, error)
	SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, token string) (*RefreshTokenRow, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllRefreshTokensForUser(ctx context.Context, userID int64) error
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *model.User) error {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRowxContext(ctx, query, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	err := r.db.GetContext(ctx, &u, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) SaveRefreshToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1,$2,$3)", userID, token, expiresAt)
	return err
}

func (r *userRepo) FindRefreshToken(ctx context.Context, token string) (*RefreshTokenRow, error) {
	var rt RefreshTokenRow
	err := r.db.GetContext(ctx, &rt, "SELECT id, user_id, token, expires_at FROM refresh_tokens WHERE token=$1", token)
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *userRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE token=$1", token)
	return err
}

func (r *userRepo) DeleteAllRefreshTokensForUser(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE user_id=$1", userID)
	return err
}
