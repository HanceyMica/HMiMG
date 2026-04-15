package middleware

import (
	"net/http"
	"strings"

	"hmimg-server-go/internal/auth"
	"hmimg-server-go/internal/config"

	"github.com/gin-gonic/gin"
)

// Gin Context 中存储用户信息的键名常量
// 用于在请求生命周期内传递已认证用户的信息
const (
	ContextUserIDKey   = "user_id"   // 用户ID（uint32 类型）
	ContextUsernameKey = "username" // 用户名（string 类型）
	ContextRoleKey     = "role"     // 用户角色（string 类型，admin 或 user）
)

// RequireAuth 认证中间件
// 验证请求头中的 JWT 令牌，确保用户已登录
// 验证通过后，将用户信息（ID、用户名、角色）存入 Gin Context，供后续处理器使用
//
// 使用方式：router.Use(middleware.RequireAuth(cfg))
//
// 参数：
//   - cfg: 应用配置对象，包含 JWT 签名密钥
//
// 返回值：
//   - gin.HandlerFunc: Gin 中间件处理函数
func RequireAuth(cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取 Authorization 字段
		header := c.GetHeader("Authorization")

		// 验证 Authorization 头存在且格式为 "Bearer <token>"
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// 提取令牌字符串（去除 "Bearer " 前缀）
		tokenString := strings.TrimPrefix(header, "Bearer ")

		// 解析并验证 JWT 令牌
		claims, err := auth.ParseToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// 将解析出的用户信息存入 Gin Context，供后续处理器使用
		c.Set(ContextUserIDKey, claims.ID)
		c.Set(ContextUsernameKey, claims.Username)
		c.Set(ContextRoleKey, claims.Role)

		// 调用后续处理器
		c.Next()
	}
}

// RequireAdmin 管理员权限中间件
// 确保已认证用户的角色为 admin
// 通常与 RequireAuth 配合使用，先验证登录再验证权限
//
// 使用方式：router.GET("/admin", middleware.RequireAuth(cfg), middleware.RequireAdmin(), adminHandler)
//
// 返回值：
//   - gin.HandlerFunc: Gin 中间件处理函数
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Context 中获取用户角色
		role, _ := c.Get(ContextRoleKey)

		// 验证角色是否为 admin
		if roleStr, ok := role.(string); !ok || roleStr != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}

		// 调用后续处理器
		c.Next()
	}
}
