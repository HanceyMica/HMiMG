// Package dbstate 进程内共享数据库连接与安装状态
// 安装向导在 Step 2 建立连接后调用 SetDB，全部 handler 通过 DB() 取实时连接
package dbstate

import (
	"sync"

	"gorm.io/gorm"
)

var (
	mu        sync.RWMutex
	database  *gorm.DB
	installed bool
	dbError   string
)

// SetDB 设置当前数据库连接（可空，表示尚未配置）
func SetDB(db *gorm.DB) {
	mu.Lock()
	defer mu.Unlock()
	database = db
}

// DB 返回当前数据库连接，未配置时返回 nil
func DB() *gorm.DB {
	mu.RLock()
	defer mu.RUnlock()
	return database
}

// SetInstalled 更新安装完成标志
func SetInstalled(v bool) {
	mu.Lock()
	defer mu.Unlock()
	installed = v
}

// Installed 返回是否已完成安装
func Installed() bool {
	mu.RLock()
	defer mu.RUnlock()
	return installed
}

// SetDBError 记录启动时数据库连接错误（安装向导展示用）
func SetDBError(msg string) {
	mu.Lock()
	defer mu.Unlock()
	dbError = msg
}

// DBError 返回最近一次数据库连接错误
func DBError() string {
	mu.RLock()
	defer mu.RUnlock()
	return dbError
}
