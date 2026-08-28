package main

import (
	"fmt"
	"log"

	"hmimg-server-go/internal/bootstrap"
	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/db"
	"hmimg-server-go/internal/server"
)

// 程序入口函数
// 负责初始化配置、连接数据库、创建表、启动 HTTP 服务器
func main() {
	// 步骤 1：加载配置
	// 从环境变量或 .env 文件读取配置（PORT、DB_*、JWT_SECRET、UPLOAD_DIR 等）
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err) // 配置加载失败，程序终止
	}

	// 步骤 2：确保上传目录存在
	// 如果目录不存在会自动创建（权限 0o755）
	if err := bootstrap.EnsureUploadDir(cfg); err != nil {
		log.Fatal(err)
	}

	// 步骤 3：连接数据库
	// 支持 MySQL 和 PostgreSQL，根据配置中的 DB_DRIVER 决定
	database, err := db.Open(cfg)
	if err != nil {
		log.Fatal(err) // 数据库连接失败，程序终止
	}

	// 步骤 4：执行数据库自动迁移
	// 根据 Model 定义创建或更新数据库表结构
	// 包括：User、Album、Collection、Image、CollectionItem、Setting
	if err := bootstrap.AutoMigrate(database); err != nil {
		log.Fatal(err)
	}

	// 步骤 5：初始化默认数据
	// 包括系统默认设置（allow_registration、max_users 等）和管理员账号
	// 首次启动时自动创建 admin 账号（用户名/密码均为 admin）
	if err := bootstrap.SeedDefaults(database); err != nil {
		log.Fatal(err)
	}

	// 步骤 6：创建并启动 HTTP 服务器
	// 注册所有路由并监听指定端口
	r := server.NewRouter(database, cfg)
	addr := fmt.Sprintf(":%d", cfg.Port) // 格式化为 ":9108" 这样的字符串
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
