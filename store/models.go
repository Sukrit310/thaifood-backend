// ไฟล์ store/models.go: กำหนดโครงสร้างข้อมูล JSON, input/filter และสัญญา Repository ที่แยก handler ออกจากฐานข้อมูล
package store

import (
	"context"
	"errors"
	"time"
)

// Error กลางที่ handler ใช้แปลงผลลัพธ์จากฐานข้อมูลเป็น HTTP โดยไม่พึ่งรายละเอียด driver
var (
	ErrNotFound   = errors.New("not found")
	ErrConflict   = errors.New("conflict")
	ErrForeignKey = errors.New("invalid reference")
	ErrForbidden  = errors.New("forbidden")
	ErrLastAdmin  = errors.New("last admin")
)

// Member เก็บข้อมูลสมาชิก; json:"-" ทำให้ PasswordHash ไม่ถูกส่งกลับใน API response
type Member struct {
	ID           int64  `json:"id"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	Role         string `json:"role"`
	PasswordHash string `json:"-"`
}

// Province เป็นข้อมูลจังหวัดที่ใช้แสดงชื่อไทย ภูมิภาค และเป็นรายการอ้างอิงของร้าน
type Province struct {
	ID     int64  `json:"id"`
	NameTh string `json:"nameTh"`
	Region string `json:"region"`
}

// Category เป็นประเภทอาหารที่ใช้จัดหมวดหมู่และกรองร้าน
type Category struct {
	ID     int64  `json:"id"`
	NameTh string `json:"nameTh"`
}

// Restaurant คือข้อมูลร้านสำหรับ API; pointer รองรับค่า null สำหรับพิกัด รูป และคะแนนที่ยังไม่มี
type Restaurant struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Province    string   `json:"province"`
	Category    string   `json:"category"`
	Address     string   `json:"address"`
	Description string   `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	ImageURL    *string  `json:"imageUrl"`
	Rating      *float64 `json:"rating"`
	ReviewCount int64    `json:"reviewCount"`
}

// ReviewMember เปิดเผยเฉพาะ ID และชื่อแสดงผลของผู้เขียนรีวิว
type ReviewMember struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"displayName"`
}

// Review รวมคะแนน ข้อความ เวลา และข้อมูลผู้เขียนที่เปิดเผยได้
type Review struct {
	ID           int64        `json:"id"`
	RestaurantID int64        `json:"restaurantId"`
	Rating       int          `json:"rating"`
	Text         string       `json:"text"`
	CreatedAt    time.Time    `json:"createdAt"`
	Member       ReviewMember `json:"member"`
}

// RestaurantFilter เก็บชื่อจังหวัด/ประเภท คำค้น และการแบ่งหน้าที่ผ่าน validation จาก handler แล้ว
type RestaurantFilter struct {
	Province, Category, Query string
	Page, Limit               int
}

// RestaurantInput ใช้ pointer แยกฟิลด์ที่ละไว้เพื่อคงค่าเดิมตอนแก้ไข
// ImageURLSet แยกการไม่ส่ง imageUrl ออกจากการส่ง null เพื่อล้างรูป
type RestaurantInput struct {
	Name        *string  `json:"name"`
	ProvinceID  *int64   `json:"provinceId"`
	CategoryID  *int64   `json:"categoryId"`
	Address     *string  `json:"address"`
	Description *string  `json:"description"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	ImageURL    *string  `json:"imageUrl"`
	ImageURLSet bool     `json:"-"`
}

// Repository กำหนดงานข้อมูลที่ handler เรียกใช้ โดยไม่ผูกกับ PostgreSQL implementation โดยตรง
type Repository interface {
	Ping(context.Context) error
	CreateMember(context.Context, string, string, string) (Member, error)
	MemberByEmail(context.Context, string) (Member, error)
	MemberByID(context.Context, int64) (Member, error)
	UpdateMember(context.Context, int64, string) (Member, error)
	UpdateMemberRole(context.Context, int64, int64, string) (Member, error)
	UpdatePassword(context.Context, int64, string, string) error
	ListMembers(context.Context) ([]Member, error)
	ListProvinces(context.Context) ([]Province, error)
	ListCategories(context.Context) ([]Category, error)
	ListRestaurants(context.Context, RestaurantFilter) ([]Restaurant, int64, error)
	RestaurantByID(context.Context, int64) (Restaurant, error)
	CreateRestaurant(context.Context, RestaurantInput) (Restaurant, error)
	UpdateRestaurant(context.Context, int64, RestaurantInput) (Restaurant, error)
	DeleteRestaurant(context.Context, int64) error
	ListReviews(context.Context, int64) ([]Review, error)
	CreateReview(context.Context, int64, int64, int, string) (Review, error)
	DeleteReview(context.Context, int64, int64, bool) error
	ListFavorites(context.Context, int64) ([]Restaurant, error)
	CreateFavorite(context.Context, int64, int64) error
	DeleteFavorite(context.Context, int64, int64) error
}
