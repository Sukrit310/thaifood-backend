// ไฟล์ auth/auth.go: จัดการรหัสผ่านด้วย bcrypt และสร้าง/ตรวจสอบ JWT สำหรับระบุตัวสมาชิก
// Package auth hashes member passwords and issues the API's access tokens.
package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	issuer   = "thaifood-guide"
	audience = "thaifood-api"
)

// ErrInvalidToken ใช้รายงาน token หรือการตั้งค่าผู้ออก token ที่ไม่ผ่านการตรวจสอบ
var ErrInvalidToken = errors.New("invalid access token")

// Issuer เก็บ secret สำหรับลงลายเซ็น/ตรวจ JWT และอายุ token; TTL ไม่เป็นบวกจะใช้หนึ่งชั่วโมง
type Issuer struct {
	Secret []byte
	TTL    time.Duration
}

// Issue สร้าง JWT แบบ HS256 โดยเก็บรหัสสมาชิกใน subject และกำหนดอายุของ token
func (i *Issuer) Issue(memberID int64) (string, error) {
	if i == nil || len(i.Secret) < 32 || memberID < 1 {
		return "", ErrInvalidToken
	}
	ttl := i.TTL
	if ttl <= 0 {
		ttl = time.Hour
	}
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Audience:  jwt.ClaimStrings{audience},
		Subject:   strconv.FormatInt(memberID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(i.Secret)
}

// Parse ตรวจลายเซ็น อัลกอริทึม ผู้สร้าง ผู้รับ และเวลาของ JWT ก่อนคืนรหัสสมาชิก
func (i *Issuer) Parse(value string) (int64, error) {
	if i == nil || len(i.Secret) < 32 {
		return 0, ErrInvalidToken
	}
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(value, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return i.Secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer(issuer),
		jwt.WithAudience(audience), jwt.WithExpirationRequired(), jwt.WithIssuedAt(), jwt.WithStrictDecoding())
	if err != nil || token == nil || !token.Valid || claims.IssuedAt == nil {
		return 0, ErrInvalidToken
	}
	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrInvalidToken
	}
	return id, nil
}

// HashPassword สร้าง bcrypt hash จากรหัสผ่าน 8–72 ไบต์ โดยคงช่องว่างในรหัสผ่านไว้
func HashPassword(password string) (string, error) {
	if len(password) < 8 || len(password) > 72 {
		return "", errors.New("password must contain 8 to 72 bytes")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword เปรียบเทียบรหัสผ่านกับ bcrypt hash และคืน false เมื่อไม่ตรงหรือ hash ผิดรูปแบบ
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
