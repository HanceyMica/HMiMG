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

	registerErr := h.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&models.User{}).Count(&count).Error; err != nil {
			return &registerFailure{http.StatusInternalServerError, "Failed to check user limit"}
		}
		if count >= int64(maxUsers) {
			return &registerFailure{http.StatusForbidden, "User limit reached"}
		}
		var existing models.User
		err := tx.First(&existing, "username = ?", req.Username).Error
		if err == nil {
			return &registerFailure{http.StatusBadRequest, "Username taken"}
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return &registerFailure{http.StatusInternalServerError, "Failed to check username"}
		}
		if err := tx.Create(&user).Error; err != nil {
			var dup models.User
			if q := tx.First(&dup, "username = ?", req.Username).Error; q == nil {
				return &registerFailure{http.StatusBadRequest, "Username taken"}
			}
			return &registerFailure{http.StatusInternalServerError, "Failed to create user"}
		}
		return nil
	})
	if registerErr != nil {
		var rf *registerFailure
		if errors.As(registerErr, &rf) {
			c.JSON(rf.status, gin.H{"error": rf.message})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}

type registerFailure struct {
	status  int
	message string
}

func (e *registerFailure) Error() string {
	return e.message
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid old password"})
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

func (h AuthHandler) ListUsers(c *gin.Context) {
	var users []models.User
	if err := h.DB.Order("id ASC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load users"})
		return
	}
	out := make([]gin.H, 0, len(users))
	for _, u := range users {
		out = append(out, gin.H{
			"id":         u.ID,
			"username":   u.Username,
			"email":      u.Email,
			"phone":      u.Phone,
			"role":       u.Role,
			"created_at": u.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, out)
}

type updateRoleRequest struct {
	Role string `json:"role"`
}

func (h AuthHandler) UpdateUserRole(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil || (req.Role != "admin" && req.Role != "user") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	var target models.User
	if err := h.DB.First(&target, "id = ?", uint32(id64)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	selfIDAny, _ := c.Get(middleware.ContextUserIDKey)
	selfID, _ := selfIDAny.(uint32)
	if target.ID == selfID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot change your own role"})
		return
	}
	if target.Role == "admin" && req.Role == "user" {
		var adminCount int64
		if err := h.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check admins"})
			return
		}
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot demote the last admin"})
			return
		}
	}

	if err := h.DB.Model(&models.User{}).Where("id = ?", target.ID).Update("role", req.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Role updated"})
}

func (h AuthHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	targetID := uint32(id64)

	selfIDAny, _ := c.Get(middleware.ContextUserIDKey)
	selfID, _ := selfIDAny.(uint32)
	if targetID == selfID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete your own account"})
		return
	}

	var target models.User
	if err := h.DB.First(&target, "id = ?", targetID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if target.Role == "admin" {
		var adminCount int64
		if err := h.DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check admins"})
			return
		}
		if adminCount <= 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete the last admin"})
			return
		}
	}

	if err := h.DB.Delete(&models.User{}, "id = ?", targetID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func defaultString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}
