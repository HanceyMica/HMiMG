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
		TrustProxy:   trustProxy,
	}, nil
}

func resolveJWTSecret() (string, error) {
	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret != "" && secret != defaultJWTSecret && len(secret) >= 16 {
		return secret, nil
	}
	if os.Getenv("GIN_MODE") == "release" {
		return "", errors.New("JWT_SECRET must be set to a strong random value (>= 16 chars) in release mode")
	}
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	generated := hex.EncodeToString(buf)
	log.Println("WARNING: JWT_SECRET not configured; generated an ephemeral secret. Tokens invalidate on restart. Set JWT_SECRET in .env or environment.")
	return generated, nil
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
