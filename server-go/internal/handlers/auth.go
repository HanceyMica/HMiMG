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

type AuthHandler struct {
	DB  *gorm.DB
	Cfg config.Config
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	var user models.User
	if err := h.DB.First(&user, "username = ?", req.Username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token, err := auth.SignToken(user, h.Cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  gin.H{"id": user.ID, "username": user.Username, "role": user.Role},
	})
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

func (h AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	settings, err := loadSettingsMap(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load settings"})
		return
	}
	if settings["allow_registration"] != "true" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Registration is closed"})
		return
	}
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

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
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

type updateProfileRequest struct {
	Username       *string `json:"username"`
	Email          *string `json:"email"`
	Phone          *string `json:"phone"`
	OldPassword    *string `json:"oldPassword"`
	NewPassword    *string `json:"password"`
	ConfirmNewPass *string `json:"confirmPassword"`
}

func (h AuthHandler) UpdateProfile(c *gin.Context) {
	userIDAny, _ := c.Get(middleware.ContextUserIDKey)
	userID, _ := userIDAny.(uint32)

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	if err := h.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	updates := map[string]interface{}{}
	if req.Username != nil && *req.Username != "" {
		updates["username"] = *req.Username
	}
	if req.Email != nil && *req.Email != "" {
		updates["email"] = *req.Email
	}
	if req.Phone != nil && *req.Phone != "" {
		updates["phone"] = *req.Phone
	}

	passwordChanged := false
	if req.NewPassword != nil && *req.NewPassword != "" {
		if req.OldPassword == nil || *req.OldPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Old password is required to set a new password"})
			return
		}
		if req.ConfirmNewPass != nil && *req.ConfirmNewPass != "" && *req.ConfirmNewPass != *req.NewPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The new passwords do not match!"})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*req.OldPassword)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
			return
		}
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.NewPassword), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		updates["password"] = string(hashed)
		passwordChanged = true
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "passwordChanged": false})
		return
	}
	if err := h.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully", "passwordChanged": passwordChanged})
}

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

func defaultString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
