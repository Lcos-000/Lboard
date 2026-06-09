package repository

import (
	"context"
	"errors"
	"whiteboard/server/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound 未找到错误
var ErrNotFound = errors.New("not found")

// UserRepository 用户仓库
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.Exec(ctx, `
INSERT INTO users (id, username, email, password_hash)
VALUES ($1, $2, $3, $4)
`, user.ID, user.Username, user.Email, user.PasswordHash)
	return err
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	row := r.db.QueryRow(ctx, `
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE email = $1
`, email)
	var user model.User
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	row := r.db.QueryRow(ctx, `
SELECT id, username, email, password_hash, created_at, updated_at
FROM users
WHERE id = $1
`, userID)
	var user model.User
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
