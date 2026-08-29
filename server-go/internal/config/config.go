package config

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const defaultJWTSecret = "YOUR_JWT_SECRET_HERE"

// DatabaseConfig 数据库连接配置，安装向导 Step 2 提交
type DatabaseConfig struct {
	Driver string
	Host   string
	Port   int
	User   string
	Pass   string
	Name   string
}

type Config struct {
	Port      int
	JWTSecret string

	DBDriver string
	DBHost   string
	DBPort   int
	DBUser   string
	DBPass   string
	DBName   string

	// DBConfigured 表示环境变量或 .env 中是否显式配置了数据库连接
	DBConfigured bool

	UploadDir string

	// FrontendDir 前端静态文件目录（前后端不分离部署），为空时后端不托管前端
	FrontendDir string

	TrustProxy bool
}

// HasDatabaseConfig 检查数据库连接是否已显式配置（而非使用默认值）
func HasDatabaseConfig() bool {
	return os.Getenv("DB_USER") != "" || os.Getenv("DB_HOST") != "" || os.Getenv("DB_NAME") != ""
}

// SaveDatabaseConfig 将数据库配置写入 .env 文件
// 保留文件中其他非 DB_* 配置行（如 JWT_SECRET、PORT）
func SaveDatabaseConfig(dbcfg DatabaseConfig) error {
	dbKeys := map[string]string{
		"DB_DRIVER":   dbcfg.Driver,
		"DB_HOST":     dbcfg.Host,
		"DB_PORT":     strconv.Itoa(dbcfg.Port),
		"DB_USER":     dbcfg.User,
		"DB_PASSWORD": dbcfg.Pass,
		"DB_NAME":     dbcfg.Name,
	}

	var kept []string
	f, err := os.Open(".env")
	if err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			key, _, ok := strings.Cut(strings.TrimPrefix(line, "export "), "=")
			if !ok {
				continue
			}
			if _, isDB := dbKeys[strings.TrimSpace(key)]; isDB {
				continue
			}
			kept = append(kept, line)
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	var out strings.Builder
	if len(kept) > 0 {
		out.WriteString(strings.Join(kept, "\n"))
		out.WriteString("\n\n")
	}
	out.WriteString("# Database (written by installer)\n")
	for _, k := range []string{"DB_DRIVER", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"} {
		out.WriteString(k + "=" + dbKeys[k] + "\n")
	}

	for k, v := range dbKeys {
		_ = os.Setenv(k, v)
	}
	return os.WriteFile(".env", []byte(out.String()), 0o600)
}

func Load() (Config, error) {
	_ = loadDotEnvIfExists(".env")

	port := getInt("PORT", 9108)
	dbPort := getInt("DB_PORT", 3306)
	uploadDir, err := filepath.Abs(getString("UPLOAD_DIR", "uploads"))
	if err != nil {
		return Config{}, err
	}
	frontendDir := ""
	if v := strings.TrimSpace(os.Getenv("FRONTEND_DIR")); v != "" {
		frontendDir, err = filepath.Abs(v)
		if err != nil {
			return Config{}, err
		}
	}
	jwtSecret, err := resolveJWTSecret()
	if err != nil {
		return Config{}, err
	}
	trustProxy := getBool("TRUST_PROXY", false)
	return Config{
		Port:         port,
		JWTSecret:    jwtSecret,
		DBDriver:     getString("DB_DRIVER", "mysql"),
		DBHost:       getString("DB_HOST", "127.0.0.1"),
		DBPort:       dbPort,
		DBUser:       getString("DB_USER", "hmimg"),
		DBPass:       getString("DB_PASSWORD", "hmimg_password"),
		DBName:       getString("DB_NAME", "hmimg_db"),
		DBConfigured: HasDatabaseConfig(),
		UploadDir:    uploadDir,
		FrontendDir:  frontendDir,
		TrustProxy:   trustProxy,
	}, nil
}

func resolveJWTSecret() (string, error) {
	envSecret := strings.TrimSpace(os.Getenv("JWT_SECRET"))

	// 环境变量已显式配置且足够强：直接使用
	if envSecret != "" && envSecret != defaultJWTSecret && len(envSecret) >= 16 {
		return envSecret, nil
	}

	// 显式配置了过短的密钥（非占位符）视为错误配置，release 模式下拒绝启动
	if envSecret != "" && envSecret != defaultJWTSecret && len(envSecret) < 16 {
		if os.Getenv("GIN_MODE") == "release" {
			return "", errors.New("JWT_SECRET is too short (>= 16 chars required) in release mode")
		}
		log.Printf("WARNING: JWT_SECRET is too short; generating a strong secret instead")
		envSecret = ""
	}

	// 环境变量缺省或为占位符时，优先使用 .env 中已持久化的密钥
	if envSecret == "" || envSecret == defaultJWTSecret {
		if fileSecret := readEnvFileValue(".env", "JWT_SECRET"); fileSecret != "" && fileSecret != defaultJWTSecret && len(fileSecret) >= 16 {
			return fileSecret, nil
		}
	}

	// 首次运行：生成强随机密钥并写回 .env，重启后保持稳定
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	generated := hex.EncodeToString(buf)
	if err := persistEnvValue(".env", "JWT_SECRET", generated); err != nil {
		// 写入失败不阻断启动：仅本次进程内生效，重启后重新生成
		log.Printf("WARNING: failed to persist generated JWT_SECRET to .env (%v); it will change on restart", err)
	} else {
		log.Println("NOTICE: JWT_SECRET not configured; generated a strong secret and saved it to .env")
	}
	_ = os.Setenv("JWT_SECRET", generated)
	return generated, nil
}

// readEnvFileValue 从 .env 风格文件中读取指定键的值（不修改进程环境）
func readEnvFileValue(filePath, key string) string {
	f, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(k) != key {
			continue
		}
		v = strings.TrimSpace(v)
		if len(v) >= 2 {
			if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
				v = v[1 : len(v)-1]
			}
		}
		return v
	}
	return ""
}

// persistEnvValue 将键值写入 .env 文件：已有则原位替换，没有则追加
func persistEnvValue(filePath, key, value string) error {
	var lines []string
	data, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		lines = strings.Split(strings.TrimRight(string(data), "\r\n"), "\n")
	}

	replaced := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if strings.HasPrefix(trimmed, "export ") {
			trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "export "))
		}
		k, _, ok := strings.Cut(trimmed, "=")
		if ok && strings.TrimSpace(k) == key {
			lines[i] = key + "=" + value
			replaced = true
			break
		}
	}
	if !replaced {
		if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) != "" {
			lines = append(lines, "")
		}
		lines = append(lines, key+"="+value)
	}

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")+"\n"), 0o600)
}

func getBool(key string, def bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if v == "" {
		return def
	}
	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	}
	return def
}

func getString(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

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

func loadDotEnvIfExists(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.TrimSpace(val)
		if key == "" {
			continue
		}
		if len(val) >= 2 {
			if (val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, val)
		}
	}
	return scanner.Err()
}
