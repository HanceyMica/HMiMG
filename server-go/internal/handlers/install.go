package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"hmimg-server-go/internal/auth"
	"hmimg-server-go/internal/bootstrap"
	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/db"
	"hmimg-server-go/internal/dbstate"
	"hmimg-server-go/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// serverVersion 当前服务端版本（安装向导展示用）
const serverVersion = "1.0.0"

// InstallHandler 安装向导处理器
// 全部接口仅在未安装状态下可用，安装完成后统一返回 404
type InstallHandler struct {
	Cfg config.Config
}

// installStep 进度状态机：
//
//	""         未建表、未配置数据库
//	"database" 数据库已连接且表已创建
//	"admin"    管理员已创建
//	"done"     安装完成（installed=true）
const (
	installStepNone     = ""
	installStepDatabase = "database"
	installStepAdmin    = "admin"
	installStepDone     = "done"
)

// Status GET /api/install/status
// 返回安装进度与环境检查信息，向导据此决定起始步骤
func (h InstallHandler) Status(c *gin.Context) {
	installed := dbstate.Installed()
	gdb := dbstate.DB()

	step := installStepNone
	if installed {
		step = installStepDone
	} else if gdb != nil {
		if s, err := bootstrap.GetSetting(gdb, "install_step"); err == nil && s != "" {
			step = s
		} else {
			// 已连接但未记录进度：表已就绪（启动时迁移或上次中断）
			step = installStepDatabase
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"installed":       installed,
		"has_db":          gdb != nil,
		"db_error":        dbstate.DBError(),
		"step":            step,
		"version":         serverVersion,
		"upload_writable": bootstrap.UploadDirWritable(h.Cfg),
	})
}

type installDatabaseRequest struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	TestOnly bool   `json:"test_only"`
}

// Database POST /api/install/database
// 测试并保存数据库连接，随后执行建表迁移
func (h InstallHandler) Database(c *gin.Context) {
	if dbstate.Installed() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}

	var req installDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	req.Driver = strings.ToLower(strings.TrimSpace(req.Driver))
	req.Host = strings.TrimSpace(req.Host)
	req.User = strings.TrimSpace(req.User)
	req.Name = strings.TrimSpace(req.Name)
	if req.Driver == "" || req.Host == "" || req.User == "" || req.Name == "" || req.Port <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All database fields are required"})
		return
	}

	testCfg := h.Cfg
	testCfg.DBDriver = req.Driver
	testCfg.DBHost = req.Host
	testCfg.DBPort = req.Port
	testCfg.DBUser = req.User
	testCfg.DBPass = req.Password
	testCfg.DBName = req.Name

	gdb, err := db.Open(testCfg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Connection failed: " + err.Error()})
		return
	}
	if err := db.Ping(gdb, 5e9); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Connection failed: " + err.Error()})
		return
	}

	if req.TestOnly {
		sqlDB, _ := gdb.DB()
		_ = sqlDB.Close()
		c.JSON(http.StatusOK, gin.H{"message": "Connection OK"})
		return
	}

	// 建表
	if err := bootstrap.AutoMigrate(gdb); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Migration failed: " + err.Error()})
		return
	}

	// 持久化连接配置到 .env（重启后直接按配置连接）
	if err := config.SaveDatabaseConfig(config.DatabaseConfig{
		Driver: req.Driver,
		Host:   req.Host,
		Port:   req.Port,
		User:   req.User,
		Pass:   req.Password,
		Name:   req.Name,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save config to .env: " + err.Error()})
		return
	}

	dbstate.SetDBError("")
	dbstate.SetDB(gdb)
	if err := bootstrap.SetSetting(gdb, "install_step", installStepDatabase); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record install step"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database configured and tables created"})
}

type installAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

// Admin POST /api/install/admin
// 创建管理员账号，要求 install_step=database
func (h InstallHandler) Admin(c *gin.Context) {
	if requireInstallStep(c, installStepDatabase) {
		return
	}

	var req installAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	if len(req.Username) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 2 characters"})
		return
	}
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters"})
		return
	}

	gdb := dbstate.DB()

	// 已存在管理员（如向导重放）则拒绝重复创建
	var count int64
	if err := gdb.Model(&models.User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing users"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Admin account already exists"})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	admin := models.User{
		Username: req.Username,
		Password: string(hashed),
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     "admin",
	}
	if err := gdb.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}
	if err := bootstrap.SetSetting(gdb, "install_step", installStepAdmin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record install step"})
		return
	}

	token, err := auth.SignToken(admin, h.Cfg.JWTSecret)
	if err != nil {
		token = ""
	}
	c.JSON(http.StatusOK, gin.H{"message": "Admin created", "token": token})
}

type installSiteRequest struct {
	WebsiteTitle      string `json:"website_title"`
	DefaultLanguage   string `json:"default_language"`
	MaxUsers          int    `json:"max_users"`
	AllowRegistration bool   `json:"allow_registration"`
}

// Site POST /api/install/site
// 写入站点设置并锁定安装，要求 install_step=admin
func (h InstallHandler) Site(c *gin.Context) {
	if requireInstallStep(c, installStepAdmin) {
		return
	}

	var req installSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	req.WebsiteTitle = strings.TrimSpace(req.WebsiteTitle)
	if req.WebsiteTitle == "" {
		req.WebsiteTitle = "HMiMG"
	}
	switch req.DefaultLanguage {
	case "en", "zh", "ja":
	default:
		req.DefaultLanguage = "zh"
	}
	if req.MaxUsers < 1 {
		req.MaxUsers = 100
	}

	gdb := dbstate.DB()

	settings := []models.Setting{
		{Key: "website_title", Value: req.WebsiteTitle},
		{Key: "default_language", Value: req.DefaultLanguage},
		{Key: "max_users", Value: intToString(req.MaxUsers)},
		{Key: "allow_registration", Value: boolToString(req.AllowRegistration)},
	}
	for _, s := range settings {
		if err := bootstrap.SetSetting(gdb, s.Key, s.Value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save settings"})
			return
		}
	}

	// 原子锁定：installed 标志写入后全部安装接口停用
	if err := bootstrap.SetSetting(gdb, "installed", "true"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to lock installer"})
		return
	}
	if err := bootstrap.SetSetting(gdb, "install_step", installStepDone); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record install step"})
		return
	}
	dbstate.SetInstalled(true)

	c.JSON(http.StatusOK, gin.H{"message": "Installation completed"})
}

// requireInstallStep 校验安装进度状态机，未满足条件时已写响应并返回 true
func requireInstallStep(c *gin.Context, expected string) bool {
	if dbstate.Installed() {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return true
	}
	gdb := dbstate.DB()
	if gdb == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Database not configured"})
		return true
	}
	step, err := bootstrap.GetSetting(gdb, "install_step")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read install step"})
		return true
	}
	if step == "" {
		step = installStepDatabase
	}
	if step != expected {
		c.JSON(http.StatusConflict, gin.H{"error": "Invalid install step: " + step})
		return true
	}
	return false
}

func intToString(v int) string {
	return strconv.Itoa(v)
}

func boolToString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}
