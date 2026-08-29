package main

import (
	"fmt"
	"log"
	"time"

	"hmimg-server-go/internal/bootstrap"
	"hmimg-server-go/internal/config"
	"hmimg-server-go/internal/db"
	"hmimg-server-go/internal/dbstate"
	"hmimg-server-go/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := bootstrap.EnsureUploadDir(cfg); err != nil {
		log.Fatal(err)
	}

	initDatabase(cfg)

	r := server.NewRouter(cfg)
	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

// initDatabase 按配置初始化数据库并判定安装状态
//
// 三种情形：
//  1. 未配置 DB（.env 无 DB_*）→ 安装模式，向导从数据库步骤开始
//  2. 已配置但连接失败 → 安装模式（向导展示连接错误，可重新提交配置）
//  3. 连接成功 → AutoMigrate（幂等，兼顾建表与升级）后判定安装状态：
//     已安装 / 旧部署回填 / 未安装（向导从管理员步骤续起）
func initDatabase(cfg config.Config) {
	if !cfg.DBConfigured {
		log.Println("No database configured; installer mode active, open /install to setup")
		return
	}

	database, err := db.Open(cfg)
	if err != nil {
		dbstate.SetDBError(err.Error())
		log.Printf("Database connect failed (%v); installer mode active, open /install to reconfigure", err)
		return
	}
	if err := db.Ping(database, 5*time.Second); err != nil {
		dbstate.SetDBError(err.Error())
		log.Printf("Database ping failed (%v); installer mode active, open /install to reconfigure", err)
		return
	}

	dbstate.SetDB(database)
	dbstate.SetDBError("")

	if err := bootstrap.AutoMigrate(database); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	if bootstrap.IsInstalled(database) {
		dbstate.SetInstalled(true)
		return
	}

	// 旧版本部署兼容：有用户数据但无安装标志 → 自动补锁
	if step, err := bootstrap.GetSetting(database, "install_step"); err == nil && step == "" {
		if bootstrap.LegacyHasUsers(database) {
			if err := bootstrap.MarkInstalled(database); err != nil {
				log.Fatalf("Failed to mark existing deployment as installed: %v", err)
			}
			dbstate.SetInstalled(true)
			log.Println("Existing deployment detected; installer locked")
			return
		}
	}

	log.Println("Not installed yet; open /install to complete setup")
}
