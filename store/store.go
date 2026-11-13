// ไฟล์ store/store.go: เชื่อม PostgreSQL ผ่าน pool และรวมตัวช่วยอ่านแถวข้อมูล/แปลง error ที่ใช้ร่วมกัน
package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store ใช้ connection pool ร่วมกันสำหรับทุกคำสั่งของ PostgreSQL repository
type Store struct{ Pool *pgxpool.Pool }

// ให้ compiler ตรวจว่า Store มี method ครบตามสัญญา Repository
var _ Repository = (*Store)(nil)

// Ping ตรวจว่าฐานข้อมูลตอบสนองภายในขอบเขตเวลาของ context
func (s *Store) Ping(ctx context.Context) error { return s.Pool.Ping(ctx) }

// dbError แปลงกรณีไม่พบแถว ข้อมูลซ้ำ และข้อจำกัดฐานข้อมูลเป็น error กลางของ repository
func dbError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	var e *pgconn.PgError
	if errors.As(err, &e) {
		switch e.Code {
		case "23505":
			return ErrConflict
		case "23503", "23502", "23514":
			return ErrForeignKey
		}
	}
	return err
}

// scanner ทำให้ตัวช่วยอ่านแถวใช้ได้ทั้ง QueryRow และ Rows ที่วนอ่านหลายแถว
type scanner interface{ Scan(...any) error }

// scanMember อ่านคอลัมน์สมาชิกตามลำดับ memberColumns รวม password hash ที่ใช้ตรวจรหัสผ่านภายใน
func scanMember(row scanner) (Member, error) {
	var m Member
	err := row.Scan(&m.ID, &m.Email, &m.PasswordHash, &m.DisplayName, &m.Role)
	return m, dbError(err)
}

// scanRestaurant อ่านรายละเอียดร้านที่ join จังหวัด ประเภทอาหาร และสถิติรีวิวแล้ว
func scanRestaurant(row scanner) (Restaurant, error) {
	var r Restaurant
	err := row.Scan(&r.ID, &r.Name, &r.Province, &r.Category, &r.Address, &r.Description, &r.Latitude, &r.Longitude, &r.ImageURL, &r.Rating, &r.ReviewCount)
	return r, dbError(err)
}

// scanReview อ่านรีวิวพร้อมผู้เขียน และแปลงเวลาเป็น UTC ก่อนส่งให้ handler
func scanReview(row scanner) (Review, error) {
	var r Review
	err := row.Scan(&r.ID, &r.RestaurantID, &r.Rating, &r.Text, &r.CreatedAt, &r.Member.ID, &r.Member.DisplayName)
	r.CreatedAt = r.CreatedAt.UTC()
	return r, dbError(err)
}

// affected ตรวจผลคำสั่งเขียนข้อมูล และคืน ErrNotFound เมื่อไม่มีแถวถูกเปลี่ยนแปลง
func affected(tag pgconn.CommandTag, err error) error {
	if err != nil {
		return dbError(err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
