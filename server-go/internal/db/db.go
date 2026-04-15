package db

import (
	"fmt"
	"strings"

	"hmimg-server-go/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open 根据配置创建并返回 GORM 数据库连接实例
// 支持 MySQL 和 PostgreSQL 两种数据库驱动
//
// 参数：
//   - cfg: 应用配置对象，包含数据库连接信息
//
// 返回值：
//   - *gorm.DB: 数据库连接实例
//   - error: 连接失败时的错误信息
func Open(cfg config.Config) (*gorm.DB, error) {
	// 将驱动名称转为小写以支持多种输入格式
	driver := strings.ToLower(cfg.DBDriver)

	switch driver {
	case "mysql", "mysql2":
		// 构建 MySQL DSN（Data Source Name）
		// 格式：用户名:密码@tcp(主机:端口)/数据库名?参数
		// 参数说明：
		//   - parseTime=true：将日期类字符串自动解析为 time.Time
		//   - charset=utf8mb4：使用 UTF-8 编码，支持 Emoji 等特殊字符
		//   - loc=Local：使用本地时区
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local",
			cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName,
		)
		// Logger 设置为 Warn 模式，减少日志输出量
		return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})

	case "pg", "postgres", "postgresql":
		// 构建 PostgreSQL DSN
		// 格式：host=主机 user=用户名 password=密码 dbname=数据库名 port=端口 sslmode=禁用
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
			cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort,
		)
		return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})

	default:
		// 不支持的数据库驱动
		return nil, fmt.Errorf("unsupported DB_DRIVER: %s", cfg.DBDriver)
	}
}
