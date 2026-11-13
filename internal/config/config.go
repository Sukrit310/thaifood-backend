// ไฟล์ internal/config/config.go: อ่านและตรวจสอบ environment variables ที่จำเป็นสำหรับเปิด API
package config

import (
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Config รวมค่าที่ใช้เชื่อมฐานข้อมูล ลงลายเซ็น JWT อนุญาต CORS และเลือกพอร์ต HTTP
type Config struct {
	DatabaseURL    string
	JWTSecret      string
	AllowedOrigins []string
	Port           string
}

// Load ตรวจ DATABASE_URL, JWT_SECRET, PORT และรายการ origin พร้อมตั้งค่าเริ่มต้นของพอร์ตและ CORS
func Load() (Config, error) {
	c := Config{DatabaseURL: os.Getenv("DATABASE_URL"), JWTSecret: os.Getenv("JWT_SECRET"), Port: os.Getenv("PORT")}
	if c.DatabaseURL == "" {
		return c, errors.New("DATABASE_URL is required")
	}
	if len(c.JWTSecret) < 32 {
		return c, errors.New("JWT_SECRET must contain at least 32 bytes")
	}
	if c.Port == "" {
		c.Port = "8080"
	}
	p, err := strconv.Atoi(c.Port)
	if err != nil || p < 1 || p > 65535 {
		return c, errors.New("PORT must be an integer from 1 to 65535")
	}
	origins := os.Getenv("ALLOWED_ORIGINS")
	if strings.TrimSpace(origins) == "" {
		origins = "http://localhost:3000"
	}
	for _, origin := range strings.Split(origins, ",") {
		origin = strings.TrimSpace(origin)
		u, err := url.Parse(origin)
		if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
			return c, errors.New("ALLOWED_ORIGINS must contain comma-separated http(s) origins without paths")
		}
		c.AllowedOrigins = append(c.AllowedOrigins, origin)
	}
	return c, nil
}
