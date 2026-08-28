package auth

import (
	"errors"
	"time"

	"hmimg-server-go/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 令牌中存储的声明（Payload）
// 包含用户的身份信息，用于在请求中识别用户身份
type Claims struct {
	ID                   uint32 `json:"id"`       // 用户ID
	Username             string `json:"username"` // 用户名
	Role                 string `json:"role"`     // 用户角色：admin 或 user
	jwt.RegisteredClaims        // JWT 标准声明（过期时间、签发时间等）
}

// SignToken 使用用户信息生成 JWT 令牌
// 采用 HS256 算法进行签名，令牌有效期为 24 小时
//
// 参数：
//   - user: 用户模型实例，包含用户ID、用户名、角色
//   - secret: JWT 签名密钥，应保密且足够复杂
//
// 返回值：
//   - string: 生成的 JWT 令牌字符串
//   - error: 生成失败时的错误信息
func SignToken(user models.User, secret string) (string, error) {
	claims := Claims{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置令牌过期时间为当前时间 + 24 小时
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			// 设置令牌签发时间为当前时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}
	// 使用 HS256 算法签名，生成令牌字符串
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// ParseToken 解析并验证 JWT 令牌
//
// 参数：
//   - tokenString: JWT 令牌字符串
//   - secret: JWT 签名密钥（必须与签发时使用的密钥一致）
//
// 返回值：
//   - Claims: 解析成功时返回令牌中的声明
//   - error: 解析或验证失败时的错误信息
func ParseToken(tokenString, secret string) (Claims, error) {
	// 解析令牌，验证签名和过期时间
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法是否为 HS256
		return []byte(secret), nil
	})
	if err != nil {
		return Claims{}, err
	}

	// 提取声明并验证令牌有效性
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, errors.New("invalid token")
	}
	return *claims, nil
}
