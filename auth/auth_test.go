package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestIssueAndParse(t *testing.T) {
	i := &Issuer{Secret: []byte(strings.Repeat("s", 32))}
	token, err := i.Issue(42)
	if err != nil {
		t.Fatal(err)
	}
	id, err := i.Parse(token)
	if err != nil || id != 42 {
		t.Fatalf("id=%d, err=%v", id, err)
	}
	if _, err := i.Issue(0); err == nil {
		t.Fatal("accepted nonpositive member ID")
	}
	if _, err := (&Issuer{Secret: []byte("short")}).Issue(1); err == nil {
		t.Fatal("accepted short secret")
	}
}

func TestRejectInvalidTokens(t *testing.T) {
	secret := []byte(strings.Repeat("s", 32))
	i := &Issuer{Secret: secret, TTL: time.Hour}
	valid := func() jwt.RegisteredClaims {
		return jwt.RegisteredClaims{Issuer: "thaifood-guide", Audience: jwt.ClaimStrings{"thaifood-api"}, Subject: "12", IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour))}
	}
	cases := []struct {
		name   string
		mutate func(*jwt.RegisteredClaims)
		method jwt.SigningMethod
		key    []byte
	}{
		{"expired", func(c *jwt.RegisteredClaims) { c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute)) }, jwt.SigningMethodHS256, secret},
		{"missing expiry", func(c *jwt.RegisteredClaims) { c.ExpiresAt = nil }, jwt.SigningMethodHS256, secret},
		{"future issuance", func(c *jwt.RegisteredClaims) { c.IssuedAt = jwt.NewNumericDate(time.Now().Add(time.Hour)) }, jwt.SigningMethodHS256, secret},
		{"missing issuance", func(c *jwt.RegisteredClaims) { c.IssuedAt = nil }, jwt.SigningMethodHS256, secret},
		{"wrong issuer", func(c *jwt.RegisteredClaims) { c.Issuer = "attacker" }, jwt.SigningMethodHS256, secret},
		{"wrong audience", func(c *jwt.RegisteredClaims) { c.Audience = jwt.ClaimStrings{"another-api"} }, jwt.SigningMethodHS256, secret},
		{"invalid subject", func(c *jwt.RegisteredClaims) { c.Subject = "0" }, jwt.SigningMethodHS256, secret},
		{"wrong algorithm", func(c *jwt.RegisteredClaims) {}, jwt.SigningMethodHS384, secret},
		{"forged signature", func(c *jwt.RegisteredClaims) {}, jwt.SigningMethodHS256, []byte(strings.Repeat("x", 32))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := valid()
			tc.mutate(&c)
			token, err := jwt.NewWithClaims(tc.method, c).SignedString(tc.key)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := i.Parse(token); err == nil {
				t.Fatal("accepted invalid token")
			}
		})
	}
	for _, token := range []string{"", "garbage", "a.b.c"} {
		if _, err := i.Parse(token); err == nil {
			t.Fatal("accepted malformed token")
		}
	}
}

func TestPasswordHashing(t *testing.T) {
	hash, err := HashPassword("Password123")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "Password123" || !CheckPassword(hash, "Password123") || CheckPassword(hash, "wrong123") {
		t.Fatal("password hash verification failed")
	}
	second, err := HashPassword("Password123")
	if err != nil {
		t.Fatal(err)
	}
	if second == hash {
		t.Fatal("password salts must differ")
	}
	for _, password := range []string{"", "short", strings.Repeat("a", 73), strings.Repeat("ก", 25)} {
		if _, err := HashPassword(password); err == nil {
			t.Fatal("accepted invalid password length")
		}
	}
	if CheckPassword("not-a-hash", "Password123") {
		t.Fatal("accepted malformed hash")
	}
}
