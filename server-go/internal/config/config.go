package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port      int
	JWTSecret string

	DBDriver string
	DBHost   string
	DBPort   int
	DBUser   string
	DBPass   string
	DBName   string

	UploadDir string
}

func Load() (Config, error) {
	_ = loadDotEnvIfExists(".env")

	port := getInt("PORT", 9108)
	dbPort := getInt("DB_PORT", 3306)
	return Config{
		Port:      port,
		JWTSecret: getString("JWT_SECRET", "YOUR_JWT_SECRET_HERE"),
		DBDriver:  getString("DB_DRIVER", "mysql"),
		DBHost:    getString("DB_HOST", "127.0.0.1"),
		DBPort:    dbPort,
		DBUser:    getString("DB_USER", "hmimg"),
		DBPass:    getString("DB_PASSWORD", "hmimg_password"),
		DBName:    getString("DB_NAME", "hmimg_db"),
		UploadDir: getString("UPLOAD_DIR", "uploads"),
	}, nil
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
