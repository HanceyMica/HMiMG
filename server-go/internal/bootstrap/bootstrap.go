package bootstrap

import (
	"errors"
	"os"

	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func EnsureUploadDir(cfg config.Config) error {
	return os.MkdirAll(cfg.UploadDir, 0o755)
}

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

func SeedDefaults(db *gorm.DB) error {
	if err := seedSettings(db); err != nil {
		return err
	}
	if err := seedAdmin(db); err != nil {
		return err
	}
	return nil
}

func seedSettings(db *gorm.DB) error {
	defaults := []models.Setting{
		{Key: "max_users", Value: "3"},
		{Key: "allow_registration", Value: "false"},
		{Key: "website_title", Value: "HMiMG"},
		{Key: "default_language", Value: "zh"},
	}
	for _, s := range defaults {
		var existing models.Setting
		err := db.Where("`key` = ?", s.Key).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := db.Create(&s).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedAdmin(db *gorm.DB) error {
	var existing models.User
	err := db.First(&existing, "username = ?", "admin").Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), 10)
	if err != nil {
		return err
	}
	admin := models.User{
		Username: "admin",
		Password: string(hashedPassword),
		Email:    "admin@yourdomaname.com",
		Phone:    "+8613200000000",
		Role:     "admin",
	}
	return db.Create(&admin).Error
}
