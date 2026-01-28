package models

import "time"

type User struct {
	ID        uint32    `gorm:"primaryKey"`
	Username  string    `gorm:"size:255;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;not null"`
	Phone     string    `gorm:"size:64;not null"`
	Role      string    `gorm:"size:32;not null;default:user"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string { return "hmimg_users" }

type Album struct {
	ID          uint32    `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	CreatedBy   *uint32   `gorm:"index"`
	CoverImage  *string   `gorm:"size:255"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Album) TableName() string { return "hmimg_albums" }

type Collection struct {
	ID          uint32    `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	CreatedBy   *uint32   `gorm:"index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Collection) TableName() string { return "hmimg_collections" }

type Image struct {
	ID           uint32    `gorm:"primaryKey"`
	Filename     string    `gorm:"size:255;not null"`
	OriginalName string    `gorm:"size:255;not null"`
	Path         string    `gorm:"size:255;not null"`
	Size         int64     `gorm:"not null"`
	Mimetype     string    `gorm:"size:255;not null"`
	AlbumID      uint32    `gorm:"index;not null"`
	UploadedBy   *uint32   `gorm:"index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (Image) TableName() string { return "hmimg_images" }

type CollectionItem struct {
	ID           uint32 `gorm:"primaryKey"`
	CollectionID uint32 `gorm:"index;not null;uniqueIndex:uq_collection_item"`
	ItemType     string `gorm:"size:32;not null;uniqueIndex:uq_collection_item"`
	ItemID       uint32 `gorm:"not null;uniqueIndex:uq_collection_item"`
}

func (CollectionItem) TableName() string { return "hmimg_collection_items" }

type Setting struct {
	Key       string    `gorm:"primaryKey;size:255"`
	Value     string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Setting) TableName() string { return "hmimg_settings" }
