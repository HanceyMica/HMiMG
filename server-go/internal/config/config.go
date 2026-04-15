package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Config 应用程序配置结构体
// 所有配置项从环境变量读取，部分字段有默认值
type Config struct {
	Port      int    // HTTP 服务监听端口，默认 9108
	JWTSecret string // JWT 令牌签名密钥，用于登录认证

	// 数据库配置
	DBDriver string // 数据库驱动类型：mysql、postgres 等
	DBHost   string // 数据库主机地址
	DBPort   int    // 数据库端口
	DBUser   string // 数据库用户名
	DBPass   string // 数据库密码
	DBName   string // 数据库名称

	UploadDir string // 用户上传文件的存储目录
}

// Load 从环境变量加载配置，返回配置实例
// 如果环境变量未设置，则使用默认值
func Load() (Config, error) {
	// 尝试加载 .env 文件（如果存在）
	_ = loadDotEnvIfExists(".env")

	// 从环境变量读取各配置项，使用默认值作为 fallback
	dbPort := getInt("DB_PORT", 3306)
	uploadDir, err := filepath.Abs(getString("UPLOAD_DIR", "uploads"))
	if err != nil {
		return Config{}, err
	}

	return Config{
		Port:      getInt("PORT", 9108),
		JWTSecret: getString("JWT_SECRET", "YOUR_JWT_SECRET_HERE"),
		DBDriver:  getString("DB_DRIVER", "mysql"),
		DBHost:    getString("DB_HOST", "127.0.0.1"),
		DBPort:    dbPort,
		DBUser:    getString("DB_USER", "hmimg"),
		DBPass:    getString("DB_PASSWORD", "hmimg_password"),
		DBName:    getString("DB_NAME", "hmimg_db"),
		UploadDir: uploadDir,
	}, nil
}

// getString 获取字符串类型的环境变量值
// 如果环境变量未设置或为空，返回默认值
func getString(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

// getInt 获取整数类型的环境变量值
// 如果环境变量未设置或解析失败，返回默认值
func getInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return parsed
}

// loadDotEnvIfExists 尝试加载 .env 文件并设置环境变量
// 如果文件不存在则静默忽略，不会返回错误
// 支持的处理格式：
//   - KEY=value
//   - export KEY=value
//   - KEY="value" 或 KEY='value'（引号会被去除）
func loadDotEnvIfExists(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在时静默忽略
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释行
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 移除 "export " 前缀（如果有）
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		// 按等号分割 key 和 value
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)

		// 跳过空 key
		if key == "" {
			continue
		}

		// 去除首尾引号
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}

		// 仅当环境变量未设置时才写入（不覆盖已存在的环境变量）
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, val)
		}
	}
	return scanner.Err()
}
