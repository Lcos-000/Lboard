package service

import (
	"context"
	"errors"
	"strings"
	"whiteboard/server/internal/auth"
	"whiteboard/server/internal/model"
	"whiteboard/server/internal/repository"
	"whiteboard/server/pkg/uuid"
)

var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrInvalidLogin    = errors.New("invalid email or password")
	ErrForbidden       = errors.New("forbidden")
)

type AuthService struct {
	users *repository.UserRepository
	jwt   *auth.JWTManager
}

// NewAuthService 创建认证服务
func NewAuthService(users *repository.UserRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{
		users: users,
		jwt:   jwt,
	}
}

// RegisterInput 注册输入参数
type RegisterInput struct {
	Username string
	Email    string
	Password string
}

// LoginInput 登录输入参数
type LoginInput struct {
	Email    string
	Password string
}

// AuthResult 认证结果
type AuthResult struct {
	User  *model.User `json:"user"`
	Token string      `json:"token"`
}

// Register 注册用户
func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*AuthResult, error) {
	// 数据清洗
	username := strings.TrimSpace(in.Username)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	password := strings.TrimSpace(in.Password)
	// 数据校验
	if username == "" || email == "" || len(password) < 6 {
		return nil, ErrInvalidArgument
	}
	// 密码哈希
	hash, err := auth.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:           uuid.NewString(),
		Username:     username,
		Email:        email,
		PasswordHash: hash,
	}
	// 用户数据存储
	if err := s.users.Create(ctx, user); err != nil {
		return nil, err
	}
	// JWT生成
	token, err := s.jwt.Generate(user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	return &AuthResult{
		User:  user,
		Token: token,
	}, nil
}

// Login 登录用户
func (s *AuthService) Login(ctx context.Context, in LoginInput) (*AuthResult, error) {
	// 数据清洗
	email := strings.ToLower(strings.TrimSpace(in.Email))
	password := strings.TrimSpace(in.Password)
	if email == "" || password == "" {
		return nil, ErrInvalidLogin
	}
	// 用户查询
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidLogin
	}

	// 检查密码是否匹配
	if !auth.CheckPassword(password, user.PasswordHash) {
		return nil, ErrInvalidLogin
	}

	// 通过检验,返回JWT
	token, err := s.jwt.Generate(user.ID, user.Email)
	if err != nil {
		return nil, err
	}
	return &AuthResult{
		User:  user,
		Token: token,
	}, nil
}

// Me 获取当前用户信息
func (s *AuthService) Me(ctx context.Context, userID string) (*model.User, error) {
	return s.users.GetByID(ctx, userID)
}
