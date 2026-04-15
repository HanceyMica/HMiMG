package models

import "time"

// ============================================================
// 用户模型 - 对应数据库表 hmimg_users
// ============================================================
type User struct {
	ID        uint32    `gorm:"primaryKey"`                   // 用户唯一标识，主键自增
	Username  string    `gorm:"size:255;uniqueIndex;not null"` // 用户名，唯一索引，不允许为空
	Password  string    `gorm:"size:255;not null"`             // 密码（bcrypt 加密后的哈希值），不允许为空
	Email     string    `gorm:"size:255;not null"`             // 邮箱，不允许为空
	Phone     string    `gorm:"size:64;not null"`              // 手机号，不允许为空，最大64字符
	Role      string    `gorm:"size:32;not null;default:user"` // 角色：admin（管理员）或 user（普通用户），默认为 user
	CreatedAt time.Time `gorm:"autoCreateTime"`                // 创建时间，GORM 自动填充
	UpdatedAt time.Time `gorm:"autoUpdateTime"`                // 更新时间，GORM 自动填充

	// GORM 关联定义
	Albums      []Album       `gorm:"foreignKey:CreatedBy"`     // 用户创建的所有相册
	Collections []Collection  `gorm:"foreignKey:CreatedBy"`    // 用户创建的所有收藏集
	Images      []Image       `gorm:"foreignKey:UploadedBy"`   // 用户上传的所有图片
}

// TableName 指定该模型对应的数据库表名
func (User) TableName() string { return "hmimg_users" }

// ============================================================
// 相册模型 - 对应数据库表 hmimg_albums
// ============================================================
type Album struct {
	ID          uint32    `gorm:"primaryKey"`                   // 相册唯一标识，主键自增
	Name        string    `gorm:"size:255;not null"`            // 相册名称，不允许为空
	Description string    `gorm:"type:text"`                     // 相册描述，支持长文本
	CreatedBy   *uint32   `gorm:"index"`                        // 创建者用户ID，外键关联 hmimg_users(id)，可为空（用户删除后设为NULL）
	CoverImage  *string   `gorm:"size:255"`                     // 封面图片文件名，可为空
	CreatedAt   time.Time `gorm:"autoCreateTime"`                // 创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`                // 更新时间

	// GORM 关联定义
	Images []Image `gorm:"foreignKey:AlbumID"`                  // 相册中的所有图片
}

// TableName 指定该模型对应的数据库表名
func (Album) TableName() string { return "hmimg_albums" }

// ============================================================
// 收藏集模型 - 对应数据库表 hmimg_collections
// 收藏集可以包含多个相册或其他收藏集，形成树形结构
// ============================================================
type Collection struct {
	ID          uint32          `gorm:"primaryKey"`             // 收藏集唯一标识，主键自增
	Name        string          `gorm:"size:255;not null"`       // 收藏集名称，不允许为空
	Description string          `gorm:"type:text"`              // 收藏集描述，支持长文本
	CreatedBy   *uint32         `gorm:"index"`                  // 创建者用户ID，外键关联 hmimg_users(id)，可为空
	CreatedAt   time.Time       `gorm:"autoCreateTime"`         // 创建时间
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`         // 更新时间

	// GORM 关联定义
	Items []CollectionItem `gorm:"foreignKey:CollectionID"`    // 收藏集中的所有条目
}

// TableName 指定该模型对应的数据库表名
func (Collection) TableName() string { return "hmimg_collections" }

// ============================================================
// 图片模型 - 对应数据库表 hmimg_images
// ============================================================
type Image struct {
	ID           uint32    `gorm:"primaryKey"`                   // 图片唯一标识，主键自增
	Filename     string    `gorm:"size:255;not null"`            // 存储的文件名（唯一），不允许为空
	OriginalName string    `gorm:"size:255;not null"`            // 用户上传时的原始文件名，不允许为空
	Path         string    `gorm:"size:255;not null"`            // 文件相对路径，不允许为空
	Size         int64     `gorm:"not null"`                    // 文件大小（字节），不允许为空
	Mimetype     string    `gorm:"size:255;not null"`            // 文件 MIME 类型（如 image/jpeg），不允许为空
	AlbumID      uint32    `gorm:"index;not null"`              // 所属相册ID，外键关联 hmimg_albums(id)，不允许为空
	UploadedBy   *uint32   `gorm:"index"`                        // 上传者用户ID，外键关联 hmimg_users(id)，可为空
	CreatedAt    time.Time `gorm:"autoCreateTime"`               // 上传时间
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`               // 更新时间

	// GORM 关联定义
	Album Album `gorm:"foreignKey:AlbumID"`                     // 图片所属的相册
}

// TableName 指定该模型对应的数据库表名
func (Image) TableName() string { return "hmimg_images" }

// ============================================================
// 收藏集条目模型 - 对应数据库表 hmimg_collection_items
// 用于建立收藏集与相册/其他收藏集之间的多对多关系
// ============================================================
type CollectionItem struct {
	ID           uint32 `gorm:"primaryKey"`                                      // 条目唯一标识，主键自增
	CollectionID uint32 `gorm:"index;not null;uniqueIndex:uq_collection_item"`   // 所属收藏集ID，外键，不允许为空，与后两列组成唯一索引防止重复
	ItemType     string `gorm:"size:32;not null;uniqueIndex:uq_collection_item"` // 条目类型：album（相册）或 collection（收藏集），不允许为空
	ItemID       uint32 `gorm:"not null;uniqueIndex:uq_collection_item"`         // 条目ID（对应相册或收藏集的ID），不允许为空，与前两列组成唯一索引
}

// TableName 指定该模型对应的数据库表名
func (CollectionItem) TableName() string { return "hmimg_collection_items" }

// ============================================================
// 系统设置模型 - 对应数据库表 hmimg_settings
// 以 key-value 形式存储系统配置
// ============================================================
type Setting struct {
	Key       string    `gorm:"primaryKey;size:255"`           // 设置项的键名，主键
	Value     string    `gorm:"size:255;not null"`             // 设置项的值，不允许为空
	CreatedAt time.Time `gorm:"autoCreateTime"`                // 创建时间
	UpdatedAt time.Time `gorm:"autoUpdateTime"`                // 更新时间
}

// TableName 指定该模型对应的数据库表名
func (Setting) TableName() string { return "hmimg_settings" }
