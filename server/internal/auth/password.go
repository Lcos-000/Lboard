package auth

import "golang.org/x/crypto/bcrypt"

// HashPassword 对密码进行哈希处理
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword 检查密码是否匹配哈希值
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// bcrypt库自动提供hash次数对齐,无需手动指定
