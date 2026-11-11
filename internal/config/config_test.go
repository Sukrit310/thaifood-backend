package config

import (
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://example:secret@localhost/db")
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
	t.Setenv("PORT", "")
	t.Setenv("ALLOWED_ORIGINS", "http://localhost:3000, https://thaifood.example")
	c, err := Load()
	if err != nil || c.Port != "8080" || len(c.AllowedOrigins) != 2 {
		t.Fatalf("load config: %#v %v", c, err)
	}
	for _, tc := range []struct{ key, value string }{
		{"JWT_SECRET", "short"}, {"DATABASE_URL", ""}, {"PORT", "0"}, {"PORT", "65536"},
		{"PORT", "bad"}, {"ALLOWED_ORIGINS", "*"}, {"ALLOWED_ORIGINS", "https://example.com/path"},
	} {
		t.Run(tc.key+tc.value, func(t *testing.T) {
			t.Setenv(tc.key, tc.value)
			if _, err := Load(); err == nil {
				t.Fatal("invalid configuration accepted")
			}
		})
	}
}
