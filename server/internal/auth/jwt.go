package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义JWT claims结构体,这是jwt库生成token还有解析token时的需要的结构体
type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager JWT管理器,负责生成和解析JWT，它包含密钥和过期时间,用于生成和验证JWT
// 收束行为实现属性和行为的封装
type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

// NewJWTManager 创建JWT管理器
func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

// Generate 生成JWT
func (m *JWTManager) Generate(userID, email string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  userID,
			IssuedAt: jwt.NewNumericDate(now),
			// 过期时间为当前时间加上过期时间
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成JWT token字符串，使用HS256签名算法
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Parse 解析JWT
func (m *JWTManager) Parse(tokenString string) (*Claims, error) {
	// 解析JWT token字符串，这里的token是一个*jwt.Token结构体
	// 自动检查token是否过期
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	// 类型断言防止载荷不是Claims类型
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
