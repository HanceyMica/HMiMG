package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"hmimg-server-go/internal/auth"
	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/middleware"
	"hmimg-server-go/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler 认证处理器，负责用户登录、注册、个人信息更新
type AuthHandler struct {
	DB  *gorm.DB         // 数据库连接实例
	Cfg config.Config    // 应用配置，包含 JWT 密钥等
}

// loginRequest 登录请求体结构
type loginRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码（明文）
}

// Login 处理用户登录请求
// 验证用户名和密码，成功则返回 JWT 令牌
//
// POST /api/login
// 请求体：{"username": "...", "password": "..."}
// 成功响应：{"token": "...", "user": {"id": ..., "username": "...", "role": "..."}}
// 失败响应：401 {"error": "Invalid credentials"}
func (h AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	// 绑定并验证 JSON 请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 根据用户名查询用户
	var user models.User
	if err := h.DB.First(&user, "username = ?", req.Username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 使用 bcrypt 验证密码
	// CompareHashAndPassword 会自动处理哈希比对，无需关心密码加密细节
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 密码验证成功，生成 JWT 令牌
	token, err := auth.SignToken(user, h.Cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}

	// 返回令牌和用户信息
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username, "role": user.Role},
	})
}

// registerRequest 注册请求体结构
type registerRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码（明文）
	Email    string `json:"email"`    // 邮箱
	Phone    string `json:"phone"`    // 手机号
}

// Register 处理用户注册请求
// 根据系统设置判断是否允许注册，并检查用户数量限制
//
// POST /api/register
// 请求体：{"username": "...", "password": "...", "email": "...", "phone": "..."}
// 成功响应：{"message": "Registered successfully"}
// 失败响应：403 {"error": "Registration is closed"} 或 403 {"error": "User limit reached"}
func (h AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 加载系统设置，检查是否允许注册
	settings, err := loadSettingsMap(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load settings"})
		return
	}
	if settings["allow_registration"] != "true" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Registration is closed"})
		return
	}

	// 检查用户数量是否达到上限
	maxUsers, _ := strconv.Atoi(defaultString(settings["max_users"], "100"))
	var count int64
	if err := h.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user limit"})
		return
	}
	if count >= int64(maxUsers) {
		c.JSON(http.StatusForbidden, gin.H{"error": "User limit reached"})
		return
	}

	// 检查用户名是否已被使用
	var existing models.User
	err = h.DB.First(&existing, "username = ?", req.Username).Error
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username taken"})
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check username"})
		return
	}

	// 使用 bcrypt 加密密码
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 创建新用户，角色默认为 "user"
	user := models.User{
		Username: req.Username,
		Password: string(hashed),
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     "user",
	}
	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}

// updateProfileRequest 更新个人信息请求体结构
type updateProfileRequest struct {
	Username       *string `json:"username"`        // 新用户名（可选）
	Email          *string `json:"email"`           // 新邮箱（可选）
	Phone          *string `json:"phone"`           // 新手机号（可选）
	OldPassword    *string `json:"oldPassword"`    // 旧密码（修改密码时必填）
	NewPassword    *string `json:"password"`       // 新密码（修改密码时填写）
	ConfirmNewPass *string `json:"confirmPassword"` // 确认新密码（修改密码时填写）
}

// UpdateProfile 处理更新个人信息请求
// 可更新用户名、邮箱、手机号，如需修改密码需提供旧密码
//
// PUT /api/admin/update（需认证）
// 请求体：{"username": "...", "email": "...", "phone": "...", "oldPassword": "...", "password": "...", "confirmPassword": "..."}
// 成功响应：{"message": "Profile updated successfully", "passwordChanged": true/false}
func (h AuthHandler) UpdateProfile(c *gin.Context) {
	// 从 Context 获取当前登录用户的 ID（由 RequireAuth 中间件设置）
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 查询当前用户信息
	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 准备更新的字段映射
	updates := map[string]interface{}{}

	// 处理用户名更新
	if req.Username != nil && *req.Username != "" {
		updates["username"] = *req.Username
	}
	// 处理邮箱更新
	if req.Email != nil && *req.Email != "" {
		updates["email"] = *req.Email
	}
	// 处理手机号更新
	if req.Phone != nil && *req.Phone != "" {
		updates["phone"] = *req.Phone
	}

	// 处理密码修改
	passwordChanged := false
	if req.NewPassword != nil && *req.NewPassword != "" {
		// 修改密码必须提供旧密码
		if req.OldPassword == nil || *req.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Old password is required to set a new password"})
			return
		}
		// 确认新密码必须与新密码一致
		if req.ConfirmNewPass != nil && *req.ConfirmNewPass != "" && *req.ConfirmNewPass != *req.NewPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The new passwords do not match!"})
			return
		}
		// 验证旧密码是否正确
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*req.OldPassword)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
			return
		}
		// 用 bcrypt 加密新密码
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.NewPassword), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updates["password"] = string(hashed)
		passwordChanged = true
	}

	// 如果没有需要更新的字段，直接返回
	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "passwordChanged": false})
		return
	}

	// 执行数据库更新
	if err := h.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "passwordChanged": passwordChanged})
}

// loadSettingsMap 从数据库加载所有设置项为 Map 结构
// 用于快速查询各设置项的值
func loadSettingsMap(db *gorm.DB) (map[string]string, error) {
	var settings []models.Setting
	if err := db.Find(&settings).Error; err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, s := range settings {
		out[s.Key] = s.Value
	}
	return out, nil
}

// defaultString 如果值为空字符串则返回默认值
func defaultString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
