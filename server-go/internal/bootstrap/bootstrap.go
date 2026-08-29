// Package bootstrap 启动引导：目录准备、表结构迁移、安装状态判定
package bootstrap

import (
	"errors"
	"os"

	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/models"

	"gorm.io/gorm"
)

// EnsureUploadDir 确保上传目录存在，不存在则自动创建
// 使用 0o755 权限（所有者读写执行，其他人读执行）
func EnsureUploadDir(cfg config.Config) error {
	return os.MkdirAll(cfg.UploadDir, 0o755)
}

// UploadDirWritable 检查上传目录是否可写（安装向导环境检查用）
func UploadDirWritable(cfg config.Config) bool {
	probe := cfg.UploadDir + string(os.PathSeparator) + ".write_test"
	if err := os.WriteFile(probe, []byte("ok"), 0o600); err != nil {
		return false
	}
	_ = os.Remove(probe)
	return true
}

// AutoMigrate 执行数据库自动迁移
// 根据 Model 定义自动创建或更新数据库表结构
// 支持的表：User、Album、Collection、Image、CollectionItem、Setting
func AutoMigrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(
		&models.User{},
		&models.Album{},
		&models.Collection{},
		&models.Image{},
		&models.CollectionItem{},
		&models.Setting{},
	)
}

// IsInstalled 判断是否已完成安装（settings 表存在且 installed=true）
// settings 表不存在视为未安装（全新库或未执行建表）
func IsInstalled(gdb *gorm.DB) bool {
	if !gdb.Migrator().HasTable(&models.Setting{}) {
		return false
	}
	v, err := GetSetting(gdb, "installed")
	return err == nil && v == "true"
}

// MarkInstalled 写入安装完成标志
func MarkInstalled(gdb *gorm.DB) error {
	return SetSetting(gdb, "installed", "true")
}

// LegacyHasUsers 判断是否为旧版本部署（users 表已存在且有数据但缺 installed 标志）
// 用于升级兼容：旧部署启动时自动补装标志，无需重新走安装向导
func LegacyHasUsers(gdb *gorm.DB) bool {
	if !gdb.Migrator().HasTable(&models.User{}) {
		return false
	}
	var count int64
	if err := gdb.Model(&models.User{}).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

// GetSetting 读取单个设置项，不存在返回空字符串
func GetSetting(gdb *gorm.DB, key string) (string, error) {
	var s models.Setting
	err := gdb.Where("`key` = ?", key).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return s.Value, nil
}

// SetSetting 写入或更新单个设置项
func SetSetting(gdb *gorm.DB, key, value string) error {
	var s models.Setting
	err := gdb.Where("`key` = ?", key).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gdb.Create(&models.Setting{Key: key, Value: value}).Error
	}
	if err != nil {
		return err
	}
	return gdb.Model(&s).Update("value", value).Error
}
