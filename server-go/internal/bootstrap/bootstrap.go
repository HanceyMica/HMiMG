package bootstrap

import (
	"errors"
	"os"

	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// EnsureUploadDir 确保上传目录存在，不存在则自动创建
// 使用 0o755 权限（所有者读写执行，其他人读执行）
//
// 参数：
//   - cfg: 应用配置对象，包含 UploadDir（上传统目录路径）
//
// 返回值：
//   - error: 创建目录失败时的错误信息
func EnsureUploadDir(cfg config.Config) error {
	return os.MkdirAll(cfg.UploadDir, 0o755)
}

// AutoMigrate 执行数据库自动迁移
// 根据 Model 定义自动创建或更新数据库表结构
// 支持的表：User、Album、Collection、Image、CollectionItem、Setting
//
// 参数：
//   - db: GORM 数据库连接实例
//
// 返回值：
//   - error: 迁移失败时的错误信息
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Album{},
		&models.Collection{},
		&models.Image{},
		&models.CollectionItem{},
		&models.Setting{},
	)
}

// SeedDefaults 执行数据库默认数据初始化
// 包括：系统默认设置、管理员账号
// 此函数在应用首次启动时调用，确保必要的初始数据存在
//
// 参数：
//   - db: GORM 数据库连接实例
//
// 返回值：
//   - error: 初始化失败时的错误信息
func SeedDefaults(db *gorm.DB) error {
	// 初始化系统默认设置
	if err := seedSettings(db); err != nil {
		return err
	}
	// 初始化管理员账号
	if err := seedAdmin(db); err != nil {
		return err
	}
	return nil
}

// seedSettings 初始化系统默认设置
// 如果设置项已存在则跳过，避免覆盖用户修改过的值
//
// 默认设置项：
//   - max_users: 最大用户数，默认 3
//   - allow_registration: 是否允许自行注册，默认 false（关闭）
//   - website_title: 网站标题，默认 "HMiMG"
//   - default_language: 默认语言，默认 "zh"
func seedSettings(db *gorm.DB) error {
	// 定义所有默认设置项
	defaults := []models.Setting{
		{Key: "max_users", Value: "3"},
		{Key: "allow_registration", Value: "false"},
		{Key: "website_title", Value: "HMiMG"},
		{Key: "default_language", Value: "zh"},
	}

	// 遍历每个默认设置项
	for _, s := range defaults {
		var existing models.Setting
		err := db.Where("`key` = ?", s.Key).First(&existing).Error

		// 如果设置项已存在，跳过
		if err == nil {
			continue
		}
		// 如果是其他错误（非"记录不存在"），返回错误
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		// 创建新的设置项记录
		if err := db.Create(&s).Error; err != nil {
			return err
		}
	}
	return nil
}

// seedAdmin 初始化管理员账号
// 仅当 admin 账号不存在时创建，密码使用 bcrypt 加密
//
// 默认管理员账号：
//   - 用户名: admin
//   - 密码: admin（运行时加密，不会以明文存储）
//   - 邮箱: admin@yourdomaname.com
//   - 手机: +8613200000000
//   - 角色: admin
func seedAdmin(db *gorm.DB) error {
	// 检查 admin 账号是否已存在
	var existing models.User
	err := db.First(&existing, "username = ?", "admin").Error
	// 如果已存在，直接返回（不重复创建）
	if err == nil {
		return nil
	}
	// 如果是其他错误，返回错误
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 使用 bcrypt 生成密码哈希
	// cost 参数为 10，表示计算强度（值越高越安全但越慢）
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), 10)
	if err != nil {
		return err
	}

	// 创建 admin 用户
	admin := models.User{
		Username: "admin",
		Password: string(hashedPassword), // 存储加密后的哈希值
		Email:    "admin@yourdomaname.com",
		Phone:    "+8613200000000",
		Role:     "admin",
	}
	return db.Create(&admin).Error
}
