package handlers

import (
	"fmt"
	"net/http"

	"hmimg-server-go/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SettingsHandler struct {
	DB *gorm.DB
}

func (h SettingsHandler) GetPublic(c *gin.Context) {
	var settings []models.Setting
	if err := h.DB.Where("`key` IN ?", []string{"allow_registration", "website_title", "default_language"}).Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load settings"})
		return
	}
	out := gin.H{"allow_registration": false, "website_title": "HMiMG", "default_language": "zh"}
	for _, s := range settings {
		if s.Key == "allow_registration" {
			out["allow_registration"] = s.Value == "true"
		}
		if s.Key == "website_title" && s.Value != "" {
			out["website_title"] = s.Value
		}
		if s.Key == "default_language" && s.Value != "" {
			out["default_language"] = s.Value
		}
	}
	c.JSON(http.StatusOK, out)
}

func (h SettingsHandler) GetAll(c *gin.Context) {
	var settings []models.Setting
	if err := h.DB.Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load settings"})
		return
	}
	out := gin.H{}
	for _, s := range settings {
		out[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, out)
}

type updateSettingsRequest struct {
	MaxUsers          interface{} `json:"max_users"`
	AllowRegistration interface{} `json:"allow_registration"`
	WebsiteTitle      interface{} `json:"website_title"`
	DefaultLanguage   interface{} `json:"default_language"`
}

func (h SettingsHandler) Update(c *gin.Context) {
	var req updateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if req.MaxUsers != nil {
		if err := h.upsert("max_users", fmt.Sprintf("%v", req.MaxUsers)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update max_users: " + err.Error()})
			return
		}
	}
	if req.AllowRegistration != nil {
		if err := h.upsert("allow_registration", fmt.Sprintf("%v", req.AllowRegistration)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update allow_registration: " + err.Error()})
			return
		}
	}
	if req.WebsiteTitle != nil {
		if err := h.upsert("website_title", fmt.Sprintf("%v", req.WebsiteTitle)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update website_title: " + err.Error()})
			return
		}
	}
	if req.DefaultLanguage != nil {
		if err := h.upsert("default_language", fmt.Sprintf("%v", req.DefaultLanguage)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default_language: " + err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}

func (h SettingsHandler) upsert(key, value string) error {
	var s models.Setting
	err := h.DB.Where("`key` = ?", key).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return h.DB.Create(&models.Setting{Key: key, Value: value}).Error
		}
		return err
	}
	return h.DB.Model(&s).Update("value", value).Error
}
