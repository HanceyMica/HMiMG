package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveJWTSecretPersistsAndStable(t *testing.T) {
	dir := t.TempDir()
	prevWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prevWd)

	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("GIN_MODE")

	s1, err := resolveJWTSecret()
	if err != nil {
		t.Fatalf("first resolve: %v", err)
	}
	if len(s1) < 16 {
		t.Fatalf("secret too short: %d", len(s1))
	}
	raw, err := os.ReadFile(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatalf(".env not written: %v", err)
	}
	if !strings.Contains(string(raw), "JWT_SECRET="+s1) {
		t.Fatalf(".env does not contain secret; content=%q", string(raw))
	}

	// 环境变量留空（模拟无注入），应复用文件中的密钥
	s2, err := resolveJWTSecret()
	if err != nil {
		t.Fatalf("second resolve: %v", err)
	}
	if s1 != s2 {
		t.Fatalf("secret regenerated: s1=%s s2=%s", s1, s2)
	}
}

func TestResolveJWTSecretPlaceholderEnvUsesFile(t *testing.T) {
	dir := t.TempDir()
	prevWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prevWd)

	os.WriteFile(filepath.Join(dir, ".env"), []byte("JWT_SECRET=0123456789abcdef0123456789abcdef\n"), 0o600)

	t.Setenv("JWT_SECRET", defaultJWTSecret) // compose 注入占位符场景
	s, err := resolveJWTSecret()
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if s != "0123456789abcdef0123456789abcdef" {
		t.Fatalf("expected file secret, got %q", s)
	}
}

func TestResolveJWTSecretShortCustomFailsInRelease(t *testing.T) {
	dir := t.TempDir()
	prevWd, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(prevWd)

	t.Setenv("JWT_SECRET", "short")
	t.Setenv("GIN_MODE", "release")
	if _, err := resolveJWTSecret(); err == nil {
		t.Fatal("expected error for short custom secret in release mode")
	}
}
