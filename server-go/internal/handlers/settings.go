package handlers

import (
	"fmt"
	"net/http"

	"hmimg-server-go/internal/dbstate"
	"hmimg-server-go/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SettingsHandler 系统设置处理器
// 负责获取和更新系统配置项
type SettingsHandler struct{}

// db 获取当前数据库连接（来自 dbstate）
func (h SettingsHandler) db() *gorm.DB {
	return dbstate.DB()
}

// GetPublic 获取公开的设置项（无需认证即可访问）
// 返回给客户端的设置，如网站标题、是否允许注册等
//
// GET /api/settings/public
// 响应示例：{"allow_registration": false, "website_title": "HMiMG", "default_language": "zh"}
func (h SettingsHandler) GetPublic(c *gin.Context) {
	var settings []models.Setting
	// 只查询公开的设置项
	if err := h.db().Where("`key` IN ?", []string{"allow_registration", "website_title", "default_language"}).Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load settings"})
		return
	}

	// 构建响应，添加默认值
	out := gin.H{
		"allow_registration": false, // 默认关闭注册
		"website_title":      "HMiMG",
		"default_language":   "zh",
	}

	// 覆盖默认值为数据库中的实际值
	for _, s := range settings {
		if s.Key == "allow_registration" {
			// 字符串 "true" 转换为布尔值
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

// GetAll 获取所有设置项（需要管理员权限）
// 返回完整的设置列表，包括内部设置
//
// GET /api/settings（需管理员认证）
// 响应示例：{"allow_registration": "false", "max_users": "3", "website_title": "HMiMG", ...}
func (h SettingsHandler) GetAll(c *gin.Context) {
	var settings []models.Setting
	if err := h.db().Find(&settings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load settings"})
		return
	}

	// 将设置列表转换为 key-value Map
	out := gin.H{}
	for _, s := range settings {
		out[s.Key] = s.Value
	}
	c.JSON(http.StatusOK, out)
}

// updateSettingsRequest 更新设置请求体结构
type updateSettingsRequest struct {
	MaxUsers          interface{} `json:"max_users"`          // 最大用户数
	AllowRegistration interface{} `json:"allow_registration"` // 是否允许注册（true/false）
	WebsiteTitle      interface{} `json:"website_title"`      // 网站标题
	DefaultLanguage   interface{} `json:"default_language"`   // 默认语言
}

// Update 处理更新设置请求（需要管理员权限）
// 支持部分更新，只更新提供的字段
//
// PUT /api/settings（需管理员认证）
// 请求体：{"max_users": 10, "allow_registration": true, "website_title": "New Title", ...}
// 成功响应：{"message": "Settings updated"}
func (h SettingsHandler) Update(c *gin.Context) {
	var req updateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 更新最大用户数（如果提供）
	if req.MaxUsers != nil {
		if err := h.upsert("max_users", fmt.Sprintf("%v", req.MaxUsers)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update max_users: " + err.Error()})
			return
		}
	}

	// 更新允许注册设置（如果提供）
	if req.AllowRegistration != nil {
		if err := h.upsert("allow_registration", fmt.Sprintf("%v", req.AllowRegistration)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update allow_registration: " + err.Error()})
			return
		}
	}

	// 更新网站标题（如果提供）
	if req.WebsiteTitle != nil {
		if err := h.upsert("website_title", fmt.Sprintf("%v", req.WebsiteTitle)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update website_title: " + err.Error()})
			return
		}
	}

	// 更新默认语言（如果提供）
	if req.DefaultLanguage != nil {
		if err := h.upsert("default_language", fmt.Sprintf("%v", req.DefaultLanguage)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update default_language: " + err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}

// upsert 插入或更新单个设置项
// 如果设置项已存在则更新，不存在则创建
//
// 参数：
//   - key: 设置项的键名
//   - value: 设置项的新值
//
// 返回值：
//   - error: 操作失败时的错误信息
func (h SettingsHandler) upsert(key, value string) error {
	var s models.Setting
	err := h.db().Where("`key` = ?", key).First(&s).Error
	if err != nil {
		// 记录不存在，创建新记录
		if err == gorm.ErrRecordNotFound {
			return h.db().Create(&models.Setting{Key: key, Value: value}).Error
		}
		return err
	}
	// 记录已存在，更新值
	return h.db().Model(&s).Update("value", value).Error
}
